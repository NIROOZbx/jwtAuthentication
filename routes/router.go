package routes

import (
	"jwt-authentication/controller"
	"jwt-authentication/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{

	router:=gin.Default()

	router.POST("/login",controller.LoginHandler)
	router.POST("/register",controller.RegisterHandler)
	router.GET("/logout", controller.Logout)
	router.GET("/home",middleware.AuthMiddleware(),controller.HomeHandler)




	return router

}