package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"oauth2-go-service/config"
	"oauth2-go-service/model"
	"time"

	"golang.org/x/oauth2"
)

func getTopWinner(startDate time.Time, endDate time.Time, category string) (model.TopWinner, error) {
	var topWinnerData model.TopWinner
	var requestTW = model.TopWinnerRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Size:      5,
		From:      0,
		OrderBy:   "totalWinLoss",
		OrderType: "desc",
		UserType:  "user",
		Category:  category,
	}
	byteRequestBody, err := json.Marshal(requestTW)
	if err != nil {
		return topWinnerData, err
	}
	requestBody := bytes.NewBuffer(byteRequestBody)
	response, err := http.Post(config.GetConfig("TOP_WINNER_URL"), "application/json", requestBody)
	if err != nil {
		return model.TopWinner{}, err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &topWinnerData)
	return topWinnerData, nil
}

func getJackpotHistory(startDate time.Time, endDate time.Time) (model.JackpotHistory, error) {
	var jackpotData model.JackpotHistory
	var requestTW model.JackpotRequest
	requestTW.Query.StartDate = startDate
	requestTW.Query.EndDate = endDate
	requestTW.Query.UserType = "user"
	requestTW.Paging.From = 0
	requestTW.Paging.Size = 10
	byteRequestBody, err := json.Marshal(requestTW)
	if err != nil {
		return jackpotData, err
	}
	requestBody := bytes.NewBuffer(byteRequestBody)
	response, err := http.Post(config.GetConfig("JACKPOT_URL"), "application/json", requestBody)
	if err != nil {
		return jackpotData, err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &jackpotData)
	return jackpotData, nil
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
