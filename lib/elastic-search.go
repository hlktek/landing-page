package lib

import (
	"oauth2-go-service/config"
	"oauth2-go-service/logger"
	"sync"

	"github.com/elastic/go-elasticsearch"
	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	client *elasticsearch.Client
)

//GetEsClient Singleton get es client
func GetEsClient() *elasticsearch.Client {
	if client == nil {
		once.Do(
			func() {
				cfg := elasticsearch.Config{
					Addresses: []string{
						config.GetConfig("ES_URL"),
					},
				}
				var err error
				client, err = elasticsearch.NewClient(cfg)
				if err != nil {
					logger.Error(logrus.Fields{
						"action": "Get es client instance",
					}, "Fail to get es clientt instance: %s", err.Error())
				}
			})
	}
	return client
}
