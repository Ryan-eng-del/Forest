package handler

import (
	"context"
	redisLib "go-gateway/lib/redis"
	"log"

	"github.com/redis/go-redis/v9"
)


const (
	DataTypeSecond = iota
	DataTypeMinute
	DataTypeHour
)


type DistributedLimiter struct {
	Name     string
	Dtype    int //0=qps 1=qpm 2=qph
	Rate     int64
	Capacity int64
}


func NewDistributedLimiter(name string, dtype int, rate, capacity int64) *DistributedLimiter {
	if dtype == DataTypeMinute {
		rate = capacity / 60
	}
	if dtype == DataTypeHour {
		rate = capacity / 3600
	}
	if rate < 1 {
		rate = 1
	}
	return &DistributedLimiter{
		Name:     name,
		Dtype:    dtype,
		Rate:     rate,
		Capacity: capacity,
	}
}

func (d *DistributedLimiter) Allow() bool {

	luaScript := redis.NewScript(`--$lua = <<<SCRIPT
	local key = KEYS[1]               --每秒一个，如：sv1_1625898937
	local limit = tonumber(ARGV[1])   --限流大小，如：20
	local current = tonumber(redis.call('get', key) or "0")
	local expire = tostring(ARGV[2])
	if ( current == 0 )
	then
		redis.call("INCRBY", key,"1") --自增
			redis.call("expire", key,expire) --2s,1m,1h超时
		return 1
	end
	if ( current + 1 > limit )
	then
			return 0
	else
			redis.call("INCRBY", key,"1") --自增
			return 1
	end
	--SCRIPT;
	`)
	
	expire := 2
	if d.Dtype == DataTypeMinute {
		expire = 60
	}
	
	if d.Dtype == DataTypeHour {
		expire = 3600
	}

	keys := []string{d.Name}
	values := []interface{}{d.Capacity, expire}
	redisClient, err := redisLib.GetRedisPoll("default")

	if err != nil {
		log.Println("DistributedLimiter Allow Error", err)
		return false
	}

	allow, err := luaScript.Run(context.Background(), redisClient, keys, values...).Int()

	if err != nil {
		log.Println("DistributedLimiter Allow Error", err)
		return false
	}

	if allow == 0 {
		return false
	} else {
		return true
	}
}