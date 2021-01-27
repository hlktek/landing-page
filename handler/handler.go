package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"oauth2-go-service/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

type AuthContent struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var authContent AuthContent

func init() {
	fmt.Println(os.Getenv("GOOGLE_CLIENT_ID"))
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     config.GetConfig("GOOGLE_CLIENT_ID"),
		ClientSecret: config.GetConfig("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func HandleMain(c *gin.Context) {
	c.HTML(http.StatusOK, "main.tmpl", gin.H{})
}

func HandleGoogleCallback(c *gin.Context) {

	content, token, err := getUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")

		// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// 	return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"id":            content.ID,
		"email":         content.Email,
		"verifiedEmail": content.VerifiedEmail,
		"picture":       content.Picture,
		"token":         token,
	})
}

// HandleGoogleLogin handle login
func HandleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
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
