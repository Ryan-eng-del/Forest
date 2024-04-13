package check

import (
	"context"
	"fmt"
	"go-gateway/install/tool"

	"github.com/redis/go-redis/v9"
)


type Redis struct{
	Host 	 string
	Port 	 string
	Pwd	 	 string
}


var (
	RedisClient Redis
)


func InitRedis() error{
	host, err := tool.Input("please enter redis host (default:127.0.0.1)", "127.0.0.1")
	if err != nil{
		return err
	}

	port, err := tool.Input("please enter redis port (default:6379)", "6379")
	if err != nil{
		return err
	}

	pwd, err := tool.Input("please enter redis pwd (default:null)", "")
	if err != nil{
		return err
	}

	redisClient := Redis{
		Host: host,
		Port: port,
		Pwd: pwd,
	}
	RedisClient = redisClient
	tool.LogInfo.Printf(fmt.Sprintf("redis connect info host:[%s] port:[%s] pwd:[%s]", host, port, pwd))
	err = redisClient.Init();if err !=nil{
		tool.LogError.Println(err)
		return err
	}
	return nil
}

func (r *Redis) Init() error {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://default:%s@%s:%s/0", r.Pwd, r.Host, r.Port))

	if err != nil {
		return err
	}
	rdb := redis.NewClient(opt)
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		tool.LogWarning.Println(err)
		return InitRedis()
	}
	tool.LogInfo.Println("connect redis success")
	return nil
}