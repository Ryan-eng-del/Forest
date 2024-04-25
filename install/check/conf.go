package check

import (
	"fmt"
	"go-gateway/install/template"
	"go-gateway/install/tool"
	"os"
	"strings"
)


var (
	ConfPath = tool.ForestGatewayPath + "/conf/dev/"
)

func initBase() error {
	tool.LogInfo.Println("init base conf")
	fileName := ConfPath + "base.toml"
	redisClient := RedisClient.Host + ":" + RedisClient.Port
	baseConf := strings.Replace(template.BaseConf, "#REDIS_CLIENT", redisClient, 1)
	baseConf = strings.Replace(baseConf, "#REDIS_PWD", RedisClient.Pwd, 1)
	err := os.WriteFile(fileName, []byte(baseConf), 0666); if err != nil{
		return err
	}
	return nil
}

func initRedis() error {
	tool.LogInfo.Println("init redis conf")
	fileName := ConfPath + "redis.toml"
	redisClient := RedisClient.Host + ":" + RedisClient.Port
	redisConf := strings.Replace(template.RedisConf, "#REDIS_CLIENT", redisClient, 1)
	redisConf = strings.Replace(redisConf, "#REDIS_PWD", RedisClient.Pwd, 1)
	err := os.WriteFile(fileName, []byte(redisConf), 0666); if err != nil{
		return err
	}
	return nil
}

func initMysql() error {
	tool.LogInfo.Println("init mysql conf")

	fileName := ConfPath + "mysql.toml"
	mysqlClient := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		MysqlClient.User,
		MysqlClient.Pwd,
		MysqlClient.Host,
		MysqlClient.Port,
		MysqlClient.Database)
	mysqlConf := strings.Replace(template.MysqlConf, "#MYSQL_CLIENT", mysqlClient, 1)
	err := os.WriteFile(fileName, []byte(mysqlConf), 0666); if err != nil{
		return err
	}
	return nil
}



func InitConf() error {
	tool.LogInfo.Println("init conf start")
	err := initBase()
	if err != nil {
		return err
	}

	err = initRedis()
	if err != nil {
		return err
	}

	err = initMysql()
	if err != nil {
		return err
	}
	tool.LogInfo.Println("init conf end")
	return nil

}