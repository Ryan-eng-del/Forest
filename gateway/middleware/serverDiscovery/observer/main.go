package observer

import (
	zookeeper "go-gateway/gateway/middleware/serverDiscovery/zooKeeper"
	"log"
)


type LoadBalanceObserver struct {
	Conf *ConcreteSubject
}


func (observer *LoadBalanceObserver) Update() {
	log.Printf("Update... LoadBalanceObserver update %+v", observer.Conf.conf)
}


func main() {
	// 被观察者 具体主体
	obConf := NewConCreateSubject(zookeeper.NodeName)
	// 观察者 抽象主体
	loadBalanceObserver := &LoadBalanceObserver{Conf: obConf}
	log.Printf("loadBalanceObserver conf is %+v", loadBalanceObserver.Conf.conf)

	// 观察者 -- 订阅
	obConf.Attach(loadBalanceObserver)
	// 被观察者 -- 发布
	obConf.UpdateConf([]string{"localhost:8001"})
}
