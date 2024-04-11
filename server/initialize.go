package server

import (
	"database/sql"
	"fmt"
	baseLib "go-gateway/lib/base"
	confLib "go-gateway/lib/conf"
	lib "go-gateway/lib/func"
	logLib "go-gateway/lib/log"
	mysqlLib "go-gateway/lib/mysql"
	redisLib "go-gateway/lib/redis"
	viperLib "go-gateway/lib/viper"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


const (
	BaseConfName = "base"
	MysqlConfName = "mysql"
	RedisConfName = "redis"
	ZookeeperConfName= "zookeeper"
)

func Migrate() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/go_gateway?charset=utf8&parseTime=true&loc=Asia%2FChongqing&multiStatements=true")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})	
	m,err := migrate.NewWithDatabaseInstance(
			"file:///Users/max/Documents/coding/Backend/Golang/Personal/Go-Gateway/migrations",
			"go_gateway", 
			driver,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}

func InitModule(configPath string) error{
	return initModule(configPath, []string{"base", "mysql", "redis", "zookeeper"})
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
		confLib.LocalIP = ips[0]
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
			return fmt.Errorf("[ERROR] %s%s", time.Now().Format(confLib.TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	// 读取配置初始化数据库模块 (mysql + gorm)
	if lib.IsInArrayString(MysqlConfName, modules) {
		mysqlLib.MysqlLibInstance = &mysqlLib.MysqlLib{}
		mysqlLib.MysqlLibInstance.SetPath(MysqlConfName, viperLib.ViperInstance.ConfEnvPath)
		if err := mysqlLib.MysqlLibInstance.InitConf(); err != nil {
			return fmt.Errorf("[ERROR] %s%s", time.Now().Format(confLib.TimeFormat), " InitMysqlConf:"+err.Error())
		}
	}

	// 读取配置初始化缓存模块 (redis)
	if lib.IsInArrayString(RedisConfName, modules) {
		redisLib.RedisLibInstance = &redisLib.RedisLib{}
		redisLib.RedisLibInstance.SetPath(RedisConfName, viperLib.ViperInstance.ConfEnvPath)
		if err := redisLib.RedisLibInstance.InitConf(); err != nil {
			return fmt.Errorf("[ERROR] %s%s", time.Now().Format(confLib.TimeFormat), " InitRedisConf:"+err.Error())
		}
	}

		// 读取配置初始化缓存模块 (zookeeper)
		// if lib.IsInArrayString(ZookeeperConfName, modules) {
		// 	zooKeeperLib.ZkManageInstance = &zooKeeperLib.ZkManager{}
		// 	zooKeeperLib.ZkManageInstance.SetPath(ZookeeperConfName, viperLib.ViperInstance.ConfEnvPath)
		// 	if err := zooKeeperLib.ZkManageInstance.InitConf(); err != nil {
		// 		return fmt.Errorf("[ERROR] %s%s", time.Now().Format(confLib.TimeFormat), " InitZookeeperConf:"+err.Error())
		// 	}
		// }
	
	if location, err := time.LoadLocation(confLib.BaseConfInstance.TimeLocation); err != nil {
		return fmt.Errorf("[ERROR] %s%s", time.Now().Format(confLib.TimeFormat), " InitTimeLocation:"+err.Error())
	} else {
		confLib.TimeLocation = location
	}

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("-----------------------------")
	return nil
}

func DestroyModule() {
	log.Println("-----------------------------")
	log.Printf(" [INFO] %s\n", " start destroy resources.")
	mysqlLib.CloseDB()
	redisLib.CloseRedis()
	logLib.Close()
	log.Printf(" [INFO] %s\n", " success destroy resources.")
	log.Println("-----------------------------")
}