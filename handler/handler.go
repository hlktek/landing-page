package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//internal service
	pb "oauth2-go-service/auth"
	"oauth2-go-service/config"
	"oauth2-go-service/data"
	"oauth2-go-service/logger"
	"oauth2-go-service/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

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
		logger.Debug(logrus.Fields{
			"action": "Get grpc connection",
		}, "Fail to get grpc connecttion : %s", err.Error())
		logger.Error(logrus.Fields{
			"action": "Get grpc connection",
		}, "Fail to get grpc connecttion : %s", err.Error())
		return
	}
	client = pb.NewAuthServiceClient(connection)
	gob.Register(&model.SessionInfo{}) //must register to save struct data in session
}

// HandleMain handle main page
func HandleMain(c *gin.Context) {
	listGameInfo := data.ListGameInfo
	sessionData := model.SessionInfo{}
	session := sessions.Default(c)
	key := session.Get("UserID")
	if key == nil {
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"listGameInfo": listGameInfo,
		})
		return
	}
	byetData, err := json.Marshal(key)
	if err != nil {
		logger.Error(logrus.Fields{
			"action": "Handle Main",
		}, "Fail to unmarshal session key : %s", err.Error())
		return
	}
	json.Unmarshal(byetData, &sessionData)
	listGameInfoToken := []model.GameInfo{}
	for _, gameInfo := range listGameInfo {
		gameInfo.Token = sessionData.Token
		listGameInfoToken = append(listGameInfoToken, gameInfo)
	}
	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"token":        sessionData.Token,
		"displayName":  sessionData.DisplayName,
		"listGameInfo": listGameInfoToken,
	})
}

// HandleGoogleCallback handle call back oauth2
func HandleGoogleCallback(c *gin.Context) {
	content, token, err := getUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		logger.Error(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Fail to get user info : %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	responseGrpc, err := client.GetUserByUserId(context.Background(), &pb.UserIdRequest{UserId: config.GetConfig("USER_ID_PREFIX") + content.Email})
	if err != nil {
		logger.Error(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Get user by user id fail: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
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
			logger.Error(logrus.Fields{
				"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
			}, "Register fail with error : %s", err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		if responseRegister.Errors != nil {
			logger.Error(logrus.Fields{
				"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
			}, "Register fail with error : %s", responseRegister.Errors)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		logger.Debug(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Register success")
	} else {
		logger.Debug(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Account is exists")
	}
	responseSetToken, err := client.SetToken(context.Background(), &pb.SetTokenRequest{
		UserId: config.GetConfig("USER_ID_PREFIX") + content.Email,
		Token:  token,
		Secret: config.GetConfig("SECRET"),
	})

	if err != nil {
		logger.Error(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Set token fail with error : %s", responseSetToken.Errors)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	if responseSetToken.Errors != nil {
		logger.Error(logrus.Fields{
			"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
		}, "Set token fail with error : %s", responseSetToken.Errors)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	logger.Debug(logrus.Fields{
		"userId": config.GetConfig("USER_ID_PREFIX") + content.Email,
	}, "Set token success")

	dataSession := model.SessionInfo{ID: content.ID,
		Email:         content.Email,
		DisplayName:   content.DisplayName,
		VerifiedEmail: content.VerifiedEmail,
		Picture:       content.Picture,
		Token:         token,
	}
	session := sessions.Default(c)
	session.Set("UserID", dataSession)
	session.Save()
	c.Redirect(http.StatusPermanentRedirect, "/")
	return
}

// HandleGoogleLogin handle login
func HandleGoogleLogin(c *gin.Context) {
	session := sessions.Default(c)
	key := session.Get("UserID")

	if key == nil {
		url := googleOauthConfig.AuthCodeURL(oauthStateString)
		c.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
	return
}

func getUserInfo(state string, code string) (model.GoogleUserInfo, string, error) {
	var userInfo model.GoogleUserInfo
	if state != oauthStateString {
		return userInfo, "", fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return userInfo, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return userInfo, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return userInfo, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return userInfo, "", fmt.Errorf("failed to unmarshal response body: %s", err.Error())
	}
	return userInfo, token.AccessToken, nil
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
