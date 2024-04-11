package handler

import (
	lib "go-gateway/lib/const"
	"sync"
	"time"
)

var ServerCountHandler *FlowCounter


type FlowCounter struct {
	RedisFlowCountMap map[string]*DistributedCountService
	Locker sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap: map[string]*DistributedCountService{},
		Locker: sync.RWMutex{},
	}
}

func init() {
	ServerCountHandler = NewFlowCounter()
	ServiceManagerHandler.Regist(ServerCountHandler)
}


func (c *FlowCounter) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		c.GetCounter(lib.FlowServicePrefix + service.Info.ServiceName)
	}

	for _, item := range c.RedisFlowCountMap {
		for _, service := range e.DeleteService {
			if item.Name == lib.FlowServicePrefix + service.Info.ServiceName {
				item.Close()
				delete(c.RedisFlowCountMap, item.Name)
			}
		}
	}
}

func (c *FlowCounter) GetCounter(name string) (*DistributedCountService, error) {
	c.Locker.Lock()
	defer c.Locker.Unlock()


	if item, ok := c.RedisFlowCountMap[name]; ok {
		return item, nil
	}

	newCounter := NewDistributedCountService(name, 1 * time.Second)
	c.RedisFlowCountMap[name] = newCounter
	return newCounter, nil
}