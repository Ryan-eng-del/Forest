package main

import (
	grpcProxyServer "go-gateway/grpcProxy/router"
	"go-gateway/handler"
	httpProxyServer "go-gateway/httpProxy/server"
	lib "go-gateway/lib/conbra"
	"go-gateway/server"
	tcpProxyServer "go-gateway/tcpProxy/router"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Go-Gateway  API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8880
// @BasePath /api
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func main() {
	if err := lib.CmdExecute(); err != nil || lib.GetCmdPanelType() == "" {
		os.Exit(1)
	}

	if lib.GetCmdPanelType() == "proxy" {
		StartProxy()
	}
	if lib.GetCmdPanelType() == "control" {
		StartControl()
	}
	if lib.GetCmdPanelType() == "both" {
		StartBoth()
	}
}

func StartProxy () {
	err := server.InitModule(lib.GetCmdConfPath())
	if err != nil {
		panic(err)
	}
	handler.AppManagerHandler.LoadAndWatch()
	handler.ServiceManagerHandler.LoadAndWatch()

	go func() {
		httpProxyServer.ServerRun()
	}()

	go func ()  {
		tcpProxyServer.TcpManagerHandler.TcpServerRun()
	}()

	go func ()  {
		grpcProxyServer.GrpcManagerHandler.GrpcServerRun()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	httpProxyServer.HttpServerStop()
}

func StartControl() {
	err := server.InitModule(lib.GetCmdConfPath())
	if err != nil {
		panic(err)
	}
	handler.AppManagerHandler.LoadAndWatch()
	handler.ServiceManagerHandler.LoadAndWatch()

	go func () {
		server.HttpServerRun()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	server.HTTPServerStop()
}

func StartBoth() {
	err := server.InitModule(lib.GetCmdConfPath())
	if err != nil {
		panic(err)
	}
	handler.AppManagerHandler.LoadAndWatch()
	handler.ServiceManagerHandler.LoadAndWatch()

	go func() {
		httpProxyServer.ServerRun()
	}()

	go func () {
		server.HttpServerRun()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	httpProxyServer.HttpServerStop()
	server.HTTPServerStop()
}