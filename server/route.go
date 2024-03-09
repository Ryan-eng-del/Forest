package server

import (
	demoController "go-gateway/controller/demo"
	"github.com/gin-gonic/gin"
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

	v1 := router.Group("v1")
	
	// group 中间件
	v1.Use()
	{
		demoController.Register(v1)	
	}

	return router
}