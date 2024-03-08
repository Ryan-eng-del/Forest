package server

import (
	"fmt"
	baseLib "go-gateway/lib/base"
	confLib "go-gateway/lib/conf"
	lib "go-gateway/lib/func"
	viperLib "go-gateway/lib/vipper"
	"log"
	"net"
	"os"
	"time"
)

var (
	LocalIP net.IP
)


const (
	BaseConfName = "base"
	TimeFormat = "2006-01-02 15:04:05"
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
	if viperLib.ViperInstance == nil {
		viperLib.ViperInstance = &viperLib.ViperLib{}
		viperLib.ViperInstance.ParseConfPath(configPath)
	}

	if viperLib.ViperInstance.ConfEnvPath == "" || viperLib.ViperInstance.ConfEnv == "" {
		return fmt.Errorf("[ERROR] ParseConfPath failed:%s", "confEnvPath and confEnv are required")
	}

	err := viperLib.ViperInstance.InitConfig()

	if err != nil {
		return err
	}

	// 读取配置，初始化 base 模块 (log)
	if lib.IsInArrayString(BaseConfName, modules) {
		baseLib.BaseLibInstance = &baseLib.BaseLib{}
		baseLib.BaseLibInstance.SetPath(BaseConfName, viperLib.ViperInstance.ConfEnvPath)
		if err := baseLib.BaseLibInstance.InitConf(); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	// 读取配置初始化数据库模块 (mysql + gorm)


	// 读取配置初始化缓存模块 (redis)

	if location, err := time.LoadLocation(confLib.BaseConfInstance.TimeLocation); err != nil {
		confLib.TimeLocation = location
	}

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("-----------------------------")
	return nil
}