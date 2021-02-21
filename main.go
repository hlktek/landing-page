package main

import (
	"net/http"
	"oauth2-go-service/config"
	"oauth2-go-service/router"

	"oauth2-go-service/data"

	"github.com/jasonlvhit/gocron"
)

func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	go func() {
		gocron.Every(1).Minute().Do(FetchListGameBO)
		<-gocron.Start()
	}()
	routerInit := router.InitRouter()
	routerInit.Run(config.GetConfig("PORT"))
	server := &http.Server{
		Handler: routerInit,
	}
	server.ListenAndServe()
}

// FetchListGameBO fetch listgamefrom BO
func FetchListGameBO() {
	data.GetListGame()
}
