package server

import (
	adminController "go-gateway/controller/admin"
	appController "go-gateway/controller/app"
	dashboardController "go-gateway/controller/dashboard"
	serviceController "go-gateway/controller/service"
	"log"

	_ "go-gateway/docs"

	confLib "go-gateway/lib/conf"
	mids "go-gateway/middlewares"

	"github.com/gin-gonic/contrib/sessions"
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

	redisConf,ok:=confLib.RedisConfInstance.List["default"]
	if !ok {
		log.Fatalf("redis.default.config err")
	}
	store, err := sessions.NewRedisStore(10, "tcp", redisConf.ProxyList[0], redisConf.Password, []byte("secret"))


	if err != nil {
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}

	authApi := router.Group("/api")
	noAuth := router.Group("/api")
	
	noAuth.Use(
		mids.RecoveryMiddleware(),
		mids.RequestLogMiddleware(),
		mids.TranslationMiddleware(),
	)

	authApi.Use(
		sessions.Sessions("mysession", store),
		mids.RecoveryMiddleware(),
		mids.RequestLogMiddleware(),
		mids.SessionAuthMiddleware(),
		mids.TranslationMiddleware(),
	)

	{
		adminLogin := noAuth.Group("/admin_login")
		adminInfo := noAuth.Group("/admin").Use()
		adminController.Register(adminLogin)
		adminController.RegisterAuth(adminInfo)
	}

	{
		serviceRouter := authApi.Group("/service").Use()
		serviceController.Register(serviceRouter)
	}

	{
		serviceRouter := authApi.Group("/app").Use()
		appController.Register(serviceRouter)
	}

	{
		dashRouter := authApi.Group("/dashboard")
		dashboardController.DashboardRegister(dashRouter)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}