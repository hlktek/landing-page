package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "oauth2-go-service/auth"
	"oauth2-go-service/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"oauth2-go-service/logger"

	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

//AuthContent auth content model
type AuthContent struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	DisplayName   string `json:"name"`
}

type AuthSession struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	DisplayName   string `json:"name"`
	Token         string `json:"token"`
}

type ResponseWallet struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"resul"`
}

var authContent AuthContent
var responseWallet ResponseWallet
var client pb.AuthServiceClient

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     config.GetConfig("GOOGLE_CLIENT_ID"),
		ClientSecret: config.GetConfig("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
	connection, err := GetGrpcConnection()
	if err != nil {
		fmt.Println(err.Error())
	}
	client = pb.NewAuthServiceClient(connection)
	gob.Register(&AuthSession{})
}

// HandleMain handle main page
func HandleMain(c *gin.Context) {
	sessionData := AuthSession{}
	session := sessions.Default(c)
	key := session.Get("UserID")
	if key == nil {
		c.HTML(http.StatusOK, "main.tmpl", gin.H{})
	}
	byetData, err := json.Marshal(key)
	// sessionData = session.Get("UserID").(AuthContent)
	if err != nil {
		fmt.Errorf("fail to marshal session")
	}
	json.Unmarshal(byetData, &sessionData)
	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"token":       sessionData.Token,
		"displayName": sessionData.DisplayName,
	})
}

// HandleGoogleCallback handle call back oauth2
func HandleGoogleCallback(c *gin.Context) {
	content, token, err := getUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	responseGrpc, err := client.GetUserByUserId(context.Background(), &pb.UserIdRequest{UserId: config.GetConfig("USER_ID_PREFIX") + content.Email})
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	if responseGrpc.DisplayName == "" {
		requestRegister := &pb.RegisterRequest{
			UserName:    content.Email,
			Password:    "123456",
			DisplayName: content.DisplayName,
			FingerPrint: "123456",
			Captcha:     "123456",
			UserId:      config.GetConfig("USER_ID_PREFIX") + content.Email,
			Type:        "USER",
			Active:      true,
			UserDetail: &pb.UserDetail{
				Money:       200000,
				ExtraMoney:  200000,
				DisplayName: content.DisplayName,
				Fullname:    content.DisplayName,
				Email:       content.Email,
				Phone:       "0123456789",
				Address:     "",
				Birthday:    "",
				Country:     "",
				Gender:      "",
				Avatar:      "Avatar0",
				RegisterIp:  "",
				Ip:          "",
				Os:          "",
				Device:      "",
				Browser:     "",
				DenyGameIds: []string{},
				Group:       "",
				Type:        "USER",
			},
		}
		responseRegister, err := client.Register(context.Background(), requestRegister)
		if err != nil {
			fmt.Println(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		if responseRegister.Errors != nil {
			fmt.Println(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		logger.Debug(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "register success")
	} else {
		logger.Debug(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "account is exists")
	}
	responseSetToken, err := client.SetToken(context.Background(), &pb.SetTokenRequest{
		UserId: config.GetConfig("USER_ID_PREFIX") + content.Email,
		Token:  token,
		Secret: config.GetConfig("SECRET"),
	})
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	if responseSetToken.Errors != nil {
		fmt.Println(responseSetToken.Errors)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	logger.Debug(logrus.Fields{
		"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
	}, "set token success")

	// session.Set("UserID", map[string]interface{}{
	// 	"id":            content.ID,
	// 	"email":         content.Email,
	// 	"displayName":   content.DisplayName,
	// 	"verifiedEmail": content.VerifiedEmail,
	// 	"picture":       content.Picture,
	// 	"token":         token,
	// })
	dataSession := AuthSession{ID: content.ID,
		Email:         content.Email,
		DisplayName:   content.DisplayName,
		VerifiedEmail: content.VerifiedEmail,
		Picture:       content.Picture,
		Token:         token,
	}
	session := sessions.Default(c)
	session.Set("UserID", dataSession)
	session.Save()
	fmt.Println(err)

	// return dataSession
	c.Redirect(http.StatusPermanentRedirect, "/")
	// c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 	"id":            content.ID,
	// 	"email":         content.Email,
	// 	"displayName":   content.DisplayName,
	// 	"verifiedEmail": content.VerifiedEmail,
	// 	"picture":       content.Picture,
	// 	"token":         token,
	// })
}

// HandleGoogleLogin handle login
func HandleGoogleLogin(c *gin.Context) {
	// sessionData := AuthSession{}
	session := sessions.Default(c)
	key := session.Get("UserID")
	fmt.Println(key)

	if key == nil {
		url := googleOauthConfig.AuthCodeURL(oauthStateString)
		c.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
	// byetData, err := json.Marshal(key)
	// // sessionData = session.Get("UserID").(AuthContent)
	// if err != nil {
	// 	fmt.Errorf("fail to marshal session")
	// }
	// json.Unmarshal(byetData, &sessionData)
	// c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 	"id":            sessionData.ID,
	// 	"email":         sessionData.Email,
	// 	"displayName":   sessionData.DisplayName,
	// 	"verifiedEmail": sessionData.VerifiedEmail,
	// 	"picture":       sessionData.Picture,
	// 	"token":         sessionData.Token,
	// })
	c.Redirect(http.StatusTemporaryRedirect, "/")
	return

	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func getUserInfo(state string, code string) (AuthContent, string, error) {
	if state != oauthStateString {
		return authContent, "", fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return authContent, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return authContent, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return authContent, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}
	err = json.Unmarshal(contents, &authContent)
	if err != nil {
		return authContent, "", fmt.Errorf("failed to unmarshal response body: %s", err.Error())
	}
	return authContent, token.AccessToken, nil
}

// GetGrpcConnection get
func GetGrpcConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(config.GetConfig("GRPC_HOST"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, err
}
