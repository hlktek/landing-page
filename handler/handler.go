package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jinzhu/now"

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
	client           pb.AuthServiceClient
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.GetConfig("CALLBACK_URL"),
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
	startDate := now.BeginningOfDay()
	endDate := now.EndOfDay()
	jackpotData, err := getJackpotHistory(startDate, endDate)
	if err != nil {
		logger.Error(logrus.Fields{
			"action": "Get Jackpot Data",
		}, "Fail to get jackpot data: %s", err.Error())
	}
	topWinnerData, err := getTopWinner(startDate, endDate, "")
	if err != nil {
		logger.Error(logrus.Fields{
			"action": "Get Top Winner Data",
		}, "Fail to get top winner data: %s", err.Error())
	}
	listGameInfo := data.DataListGameBO.Data
	sessionData := model.SessionInfo{}
	session := sessions.Default(c)
	key := session.Get("UserID")
	if key == nil {
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"listGameInfo":  listGameInfo,
			"topWinnerData": topWinnerData.Data.Data,
			"jackpotData":   jackpotData.Data.Data,
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
	listGameInfoToken := []model.GameInfoBO{}
	for _, gameInfo := range listGameInfo {
		gameInfo.Token = sessionData.Token
		listGameInfoToken = append(listGameInfoToken, gameInfo)
	}
	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"token":         sessionData.Token,
		"displayName":   sessionData.DisplayName,
		"listGameInfo":  listGameInfoToken,
		"topWinnerData": topWinnerData.Data.Data,
		"jackpotData":   jackpotData.Data.Data,
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

//GetTopWinner by category
func GetTopWinner(c *gin.Context) {
	catetgory := c.Query("category")
	startDate := now.BeginningOfDay()
	endDate := now.EndOfDay()

	topWinnerData, err := getTopWinner(startDate, endDate, catetgory)
	if err != nil {
		logger.Debug(logrus.Fields{
			"action": "Get top winner by category",
		}, "Fail to get top winner by category : %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": topWinnerData.Data.Data})
}

//GetTopWinnerChess by category
func GetTopWinnerChess(c *gin.Context) {
	startDate := now.BeginningOfDay()
	endDate := now.EndOfDay()

	topWinnerData, err := getTopWinnerChess(startDate, endDate)
	if err != nil {
		logger.Debug(logrus.Fields{
			"action": "Get top winner chess",
		}, "Fail to get top winner chess : %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": topWinnerData.Data.Data})
}

//GetTopJackpot
func GetJackpotHistory(c *gin.Context) {
	startDate := now.BeginningOfDay()
	endDate := now.EndOfDay()

	jackpotData, err := getJackpotHistory(startDate, endDate)

	if err != nil {
		logger.Debug(logrus.Fields{
			"action": "Get top jackpot",
		}, "Fail to get top jackpot : %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jackpotData.Data.Data})
}

// GetImage get image from BO
func GetImage(c *gin.Context) {
	imageName := c.Query("imageName")
	target := config.GetConfig("GET_GAME_INFO_BO_URL")

	remote, err := url.Parse(target)
	if err != nil {
		// checkErr("parse", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = config.GetConfig("GET_IMAGE_BO_PATH") + imageName //Request API
	proxy.ServeHTTP(c.Writer, c.Request)
}

// AddWallet addwallet
func AddWallet(c *gin.Context) {
	var walletResponse model.Wallet
	userID := c.Query("userId")
	response, err := http.Get(config.GetConfig("WALLET_URL") + "?money=2000000&serverType=Staging&userId=" + userID)
	if err != nil {
		logger.Debug(logrus.Fields{
			"action": "Add wallet",
		}, "Fail to add wallet : %s", err.Error())
		logger.Error(logrus.Fields{
			"action": "Add wallet",
		}, "Fail to add wallet : %s", err.Error())

	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &walletResponse)
	c.JSON(http.StatusOK, gin.H{"data": walletResponse})
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
