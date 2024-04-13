package main

import (
	"go-gateway/install/tool"
	"go-gateway/install/check"
	"os"
)

func main() {
	var (
		err error
	)

	tool.InitSystem()

	err = check.InitRedis(); if err != nil{
		tool.LogError.Println("connect redis error")
		tool.LogError.Println(err)
		os.Exit(-1)
	}

	
	// var (
	// 	err error
	// )


}