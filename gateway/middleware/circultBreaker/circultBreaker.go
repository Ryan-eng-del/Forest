package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)
func main() {

	// 统计熔断，将其发送到 8070 服务器上，dashboard 通过 pull 数据，进行可视化
	statisticsHandler := hystrix.NewStreamHandler()
	statisticsHandler.Start()

	go http.ListenAndServe(":8070", statisticsHandler)
	
	hystrix.ConfigureCommand("hystrix", hystrix.CommandConfig{
		// 单次请求完成的超市时间
		Timeout: 1000,
		// 最大并发量
		MaxConcurrentRequests: 1,
		// 熔断后超时重试的时间，也就是熔断打开多久之后，去尝试恢复服务，转换半打开状态
		SleepWindow: 5000,
		// 验证熔断的请求数量，熔断发生之前 10s 的最小故障请求数，10s 有一个故障则熔断
		RequestVolumeThreshold: 1,
		// 根据上一个字段（10s内发生了多少次熔断），去计算熔断的错误百分比，是否达到 1%
		ErrorPercentThreshold: 1,
	})


	for i := 0; i < 1000; i++ {
		err := hystrix.Do("hystrix", func() error {
			if i == 0 {
				return errors.New("service error" + strconv.Itoa(i))
			}
			log.Println("do service ", i)
			return nil
		},func(e error) error {
			log.Println("do 兜底 service ", i)
			return e
		})

		if err != nil {
			log.Println("hystrix error: "+ err.Error())
			time.Sleep(2 * time.Second)
		}
	}
}