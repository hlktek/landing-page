package main

import (
	"net/http"
	"oauth2-go-service/router"
)

func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	routerInit := router.InitRouter()
	routerInit.Run(":3000")
	server := &http.Server{
		Handler: routerInit,
	}
	server.ListenAndServe()
}
