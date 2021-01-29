package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oauth2-go-service/logger"
	"oauth2-go-service/model"
	"strings"

	"github.com/sirupsen/logrus"
)

// DataListGameBO list data bo
var DataListGameBO = model.GameBOConfig{}

func init() {
	postBody, _ := json.Marshal(map[string]interface{}{
		"paging": map[string]int{
			"from": 0,
			"size": 10,
		},
		"query": map[string][]string{
			"ignoreField": {"icon"},
		},
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://128.199.93.192:3025/kgame/v2/game-config/games/list-games", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &DataListGameBO)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, gameBO := range DataListGameBO.Data {
		gameBoPoiter := &DataListGameBO.Data[i]
		trimEx := strings.Replace(gameBO.ExServiceID, "_", "", -1)
		gameBoPoiter.ExServiceID = trimEx
		gameBoPoiter.Link = "https://iframe.staging.ktek.io/" + gameBoPoiter.ExServiceID + "/?token="
	}

	fmt.Println(DataListGameBO.Data)

	logger.Debug(logrus.Fields{
		"action": "get-data-bo",
	}, "Get data list game from BO success")
}
