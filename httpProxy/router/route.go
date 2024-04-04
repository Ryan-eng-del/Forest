package httpRouter

import (
	oauthController "go-gateway/controller/oauth"
	httpMiddlewares "go-gateway/httpProxy/middlewares"
	"go-gateway/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter(mids ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(mids...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	{
		oauth := router.Group("/oauth")
		oauth.Use(middlewares.TranslationMiddleware())
		oauthController.Register(oauth)
	}


	router.Use(
		httpMiddlewares.HTTPAccessModeMiddleware(),
		httpMiddlewares.JWTTokenAuth(),
		httpMiddlewares.HTTPFlowCountMiddleware(),
		httpMiddlewares.HttpFlowLimitMiddleware(),
		httpMiddlewares.HTTPBlackListMiddleware(),
		
	)
	return router
}