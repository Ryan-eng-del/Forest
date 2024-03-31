package server

import (
	adminController "go-gateway/controller/admin"
	appController "go-gateway/controller/app"
	serviceController "go-gateway/controller/service"
	_ "go-gateway/docs"

	mids "go-gateway/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	// 全局中间件
	router.Use(middlewares...)
	router.GET("/ping", func (c *gin.Context)  {
		c.JSON(200, gin.H{
			"message": "Pong",
		})
	})

	api := router.Group("/api")
	api.Use(mids.RequestLogMiddleware(),mids.RecoveryMiddleware(),mids.TranslationMiddleware())

	
	{
		adminLogin := api.Group("/admin_login")
		adminInfo := api.Group("/admin").Use(mids.JWTTokenAuth())
		adminController.Register(adminLogin)
		adminController.RegisterAuth(adminInfo)
	}

	{
		serviceRouter := api.Group("/service").Use(mids.JWTTokenAuth())
		serviceController.Register(serviceRouter)
	}

	{
		serviceRouter := api.Group("/app").Use(mids.JWTTokenAuth())
		appController.Register(serviceRouter)
	}


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}