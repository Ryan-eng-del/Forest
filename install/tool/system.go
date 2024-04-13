package tool

import (
	"os"
	"runtime"
	"strings"
)

var (
	SystemType  string
	CommandType string
	CommandArgs string
	CmdRun 		string
	ForestGatewayPath	string = forestGatewayPath()
)


func InitSystem()  {
	SystemType = runtime.GOOS
	if SystemType == "windows"{
		CommandType = "cmd"
		CommandArgs = "/C"
		CmdRun = "SET GO111MODULE=on&& SET GOPROXY=https://goproxy.cn"
	} else {
		CommandType = "sh"
		CommandArgs = "-c"
		CmdRun = "export GO111MODULE=on && export GOPROXY=https://goproxy.cn"
	}
}


func GetCurrentPath() string{
	path, _ := os.Getwd()
	return strings.Replace(path, "\\", "/", -1)
}


func forestGatewayPath() string{
	path := GetCurrentPath()
	pathArr := strings.Split(path, "/")
	index := len(pathArr)
	pathArr = pathArr[0:index-1]
	path = strings.Join(pathArr, "/")
	return path
}