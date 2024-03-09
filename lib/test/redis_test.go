package test

import (
	"context"
	"fmt"
	libLog "go-gateway/lib/log"
	libRedis "go-gateway/lib/redis"
	"testing"
	"time"
)

func Test_Redis(t *testing.T) {
	SetUp()
	c, err := libRedis.GetRedisPoll("default")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// 调用SET
	type TraceKey string
	trace := libLog.NewTrace()
	ctx := context.Background()
	ctx = context.WithValue(ctx, TraceKey("trace"), trace)
	redisKey := "test_key1"
	if _, err := c.Set(ctx, redisKey, "test_dpool", 10 * time.Second).Result(); err != nil {
		t.Fatal(err)
	}

	// 调用GET
	v, err := c.Get(ctx, redisKey).Result()
	fmt.Println(v)
	if v != "test_dpool" || err != nil {
		t.Fatal("test redis get fatal!")
	}

	TearDown()
}
