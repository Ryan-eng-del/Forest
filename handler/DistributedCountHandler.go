package handler

import (
	"context"
	"fmt"
	lib "go-gateway/lib/conf"
	libConst "go-gateway/lib/const"
	libLog "go-gateway/lib/log"
	libRedis "go-gateway/lib/redis"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)



type DistributedCountService struct {
	Name        string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
	closeChan   chan bool
}

func NewDistributedCountService(name string, interval time.Duration) *DistributedCountService {
	reqCounter := &DistributedCountService{
		Name: name,
		Interval: interval,
		QPS: 0,
		Unix: 0,
		closeChan: make(chan bool, 1),
	}
	log := libLog.GetLogger()

	go func(){
		defer func() {
			if err := recover(); err != nil {
				log.Error("DistributedCountService error Exiting: %s", err)
			}
		}()

		ticker := time.NewTicker(interval)

		OUTTER:
			for {
				select {
				case <- reqCounter.closeChan:
					continue OUTTER
				case <- ticker.C:
					tickerCount := atomic.LoadInt64(&reqCounter.TickerCount)
					atomic.StoreInt64(&reqCounter.TickerCount, 0)
					currentTime := time.Now()

					dayKey := reqCounter.GetDayKey(currentTime)
					hourKey := reqCounter.GetHourKey(currentTime)

					if tickerCount > 0 {
						if err := libRedis.RedisConfPipline(func(c *redis.Client) error {
							ctx := context.Background()
							if err := c.IncrBy(ctx, dayKey, tickerCount).Err(); err != nil {
								return err
							}
							if err := c.Expire(ctx, dayKey, time.Hour * 48).Err(); err != nil {
								return err
							}

							if err := c.IncrBy(ctx, hourKey, tickerCount).Err(); err != nil {
								return err
							}

							if err := c.Expire(ctx, hourKey, time.Hour * 48).Err(); err != nil {
								return err
							}
							
							return nil
						}); err != nil {

							log.Error("Failed to set redis key %s", err)
							continue
						}
					}

					totalCount, err := reqCounter.GetDayData(currentTime);

					if err != nil {
						log.Error("Failed to get day data %s", err)
						continue
					}

					nowUnix := time.Now().Unix()

					if reqCounter.Unix == 0 {
						reqCounter.Unix = time.Now().Unix()
						continue						
					}

					
				tickerCount = totalCount - reqCounter.TotalCount

					if nowUnix > reqCounter.Unix {
						reqCounter.TotalCount = totalCount
						reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
						reqCounter.Unix = time.Now().Unix()
					}
				}
			}
	}()
	return reqCounter
}


func (o *DistributedCountService) GetDayData(t time.Time) (int64, error) {
	c, err := libRedis.GetRedisPoll("default")
	if err != nil {
		return 0, err
	}

	result, err := c.Get(context.Background(), o.GetDayKey(t)).Result()

	if err != nil {
		return 0, err
	}

	if dayCount, err := strconv.ParseInt(result, 10, 64); err != nil {
		return 0, err
	} else {
		return dayCount, nil
	}
}

func (o *DistributedCountService) GetHourData(t time.Time) (int64, error) {
	c, err := libRedis.GetRedisPoll("default")
	if err != nil {
		return 0, err
	}

	result, err := c.Get(context.Background(), o.GetHourKey(t)).Result()

	if err != nil {
		return 0, err
	}

	if dayCount, err := strconv.ParseInt(result, 10, 64); err != nil {
		return 0, err
	} else {
		return dayCount, nil
	}
}

func (o *DistributedCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(lib.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", libConst.RedisFlowDayKey, dayStr, o.Name)
}


func (o *DistributedCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", libConst.RedisFlowHourKey, hourStr, o.Name)
}

func (o *DistributedCountService) Close() {
	go func(){
		defer func () {
			if err := recover(); err != nil {
				log := libLog.GetLogger()
				log.Error("Failed to close %s", err)
			}
		}()
	}()

	if o.closeChan != nil {
		close(o.closeChan)
		o.closeChan = nil
	}
}

func (c *DistributedCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&c.TickerCount, 1)
	}()
}