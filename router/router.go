package router

import (
	"oauth2-go-service/handler"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//InitRouter init router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	r.Static("/assets", "./assets")

	r.LoadHTMLGlob("templates/*")
	r.GET("/auth/google/callback", handler.HandleGoogleCallback)
	r.GET("/auth/google", handler.HandleGoogleLogin)
	r.GET("/", handler.HandleMain)
	return r
}
