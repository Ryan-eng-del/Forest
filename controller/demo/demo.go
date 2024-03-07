package demoController

import "github.com/gin-gonic/gin"


type demoController struct {}

func (d *demoController) Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "you gays",
	})
}

func Register(g *gin.RouterGroup) {
	demoController := &demoController{}
	g.GET("/hello", demoController.Hello)
}