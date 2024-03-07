package server

import (
	"go-gateway/lib"
	viperLib "go-gateway/lib/vipper"
	"log"
	"net"
	"os"
)

var (
	LocalIP net.IP
)

var (
	ViperLib *viperLib.ViperLib
)


func InitModule(configPath string) error{
	return initModule(configPath, []string{"base", "mysql", "redis"})
}

func initModule(configPath string, modules []string) error {
	if configPath == "" {
		log.Printf("[ERROR] initModule failed: Please specify a config path like %s", "./conf/env")
		os.Exit(1)
	}

	log.Println("-----------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " Start Loading Configs.")

	// 设置本地主机 ip 信息
	ips := lib.GetLocalIPs()

	if len(ips) > 0 {
		LocalIP = ips[0]
	}

	// viper 读取配置
	if ViperLib != nil {
		ViperLib = &viperLib.ViperLib{}
		ViperLib.ParseConfPath(configPath)
	}

	if ViperLib.ConfEnvPath == "" || ViperLib.ConfEnv == "" {
		log.Printf("[ERROR] ParseConfPath failed:%s\n", "confEnvPath and confEnv are required")
	}

	err := ViperLib.InitConfig()

	if err != nil {
		log.Printf("[ERROR] ParseConfPath failed:%s\n", err)
	}


	if lib.IsInArrayString("base", modules) {
		
	}
	return nil


}