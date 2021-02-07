package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oauth2-go-service/config"
	"oauth2-go-service/logger"
	"oauth2-go-service/model"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
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

func getTopWinnerChess(startDate time.Time, endDate time.Time) (model.TopWinnerChess, error) {
	var topWinnerChessData model.TopWinnerChess
	var request = model.TopWinnerChessRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}
	byteRequestBody, err := json.Marshal(request)
	if err != nil {
		return topWinnerChessData, err
	}
	requestBody := bytes.NewBuffer(byteRequestBody)
	response, err := http.Post(config.GetConfig("TOP_WINNER_CHESS_URL"), "application/json", requestBody)
	if err != nil {
		return model.TopWinnerChess{}, err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &topWinnerChessData)
	return topWinnerChessData, nil
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

type ElasticDocs struct {
	UserID    string
	Feedback  string
	ServiceId string
	Time      int64
}

// A function for marshaling structs to JSON string
func jsonStruct(doc ElasticDocs) string {

	// Create struct instance of the Elasticsearch fields struct object
	docStruct := &ElasticDocs{
		UserID:    doc.UserID,
		Feedback:  doc.Feedback,
		ServiceId: doc.ServiceId,
		Time:      doc.Time,
	}

	fmt.Println("\ndocStruct:", docStruct)
	fmt.Println("docStruct TYPE:", reflect.TypeOf(docStruct))

	// Marshal the struct to JSON and check for errors
	b, err := json.Marshal(docStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(b)
}

func insertEs(userId string, feedback string, serviceId string, time int64) error {
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	// Create a mapping for the Elasticsearch documents
	var (
		docMap map[string]interface{}
	)
	fmt.Println("docMap:", docMap)
	fmt.Println("docMap TYPE:", reflect.TypeOf(docMap))

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.GetConfig("ES_URL"),
		},
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return fmt.Errorf("Elasticsearch connection error: %s", err)
	}
	doc1 := ElasticDocs{}
	doc1.UserID = userId
	doc1.Feedback = feedback
	doc1.Time = time
	doc1.ServiceId = serviceId
	docStr1 := jsonStruct(doc1)
	timeStr := strconv.Itoa(int(time))
	documentID := timeStr + "-" + userId + "-" + serviceId
	req := esapi.IndexRequest{
		Index:      "feed-back",
		DocumentID: documentID,
		Body:       strings.NewReader(docStr1),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, client)
	if err != nil {
		logger.Error(logrus.Fields{
			"action": "Insert Feed Back",
		}, "Insert feedback fail : %s", err.Error())
		return fmt.Errorf("Insert feedback fail : %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Debug(logrus.Fields{
			"action": "insert es",
		}, "Fail to insert es : %s", err.Error())
		return fmt.Errorf("Fail to insert es")

	} else {
		// Deserialize the response into a map.
		var resMap map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			logger.Debug(logrus.Fields{
				"action": "insert es - parse respone",
			}, "Fail to parse respone : %s", err.Error())
			return fmt.Errorf("Error parsing the response body: %s", err)
		}
		return nil
	}
}
