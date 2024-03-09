package lib

import (
	"context"
	libConf "go-gateway/lib/conf"
	libFunc "go-gateway/lib/func"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLib struct {
	ConfPath string 
}

var RedisPoll *redis.Client
var RedisLibInstance *RedisLib
var RedisMapPoll map[string]*redis.Client

func (bL *RedisLib) ParseConfig() error {
	return libFunc.ParseConfigFromFile(bL.ConfPath, libConf.RedisConfInstance)
}


func (rL *RedisLib) InitConf () error {
	if err := rL.ParseConfig(); err != nil {
		return err
	}

	for confName, conf := range libConf.RedisConfInstance.List {
		opt, err := redis.ParseURL(conf.DataSourceName)
		if err != nil {
			return err
		}

		opt.ReadTimeout = time.Millisecond  * time.Duration(conf.ReadTimeout)
		opt.WriteTimeout = time.Millisecond  * time.Duration(conf.WriteTimeout)
		opt.MaxActiveConns = conf.MaxActive
		opt.DialTimeout = time.Millisecond  * time.Duration(conf.ConnTimeout)
		opt.MaxIdleConns = conf.MaxIdle

		rdb := redis.NewClient(opt)
		_, err = rdb.Ping(context.Background()).Result()

		if err != nil {
			return err
		}
		
		if confName ==  "default" {
			RedisPoll = rdb
		}
		
		if RedisMapPoll == nil {
			RedisMapPoll = make(map[string]*redis.Client)
		}

		RedisMapPoll[confName] = rdb
	}
	return nil
}

func (rL *RedisLib) SetPath(fileName string, ConfEnvPath string)  {
	rL.ConfPath = ConfEnvPath + "/" + fileName + ".toml"
}