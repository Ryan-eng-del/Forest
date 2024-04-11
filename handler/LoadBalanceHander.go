package handler

import (
	load_balance "go-gateway/loadBalance"
	"go-gateway/model"
	"sync"
	"time"
)
var LoadBalancerHandler *LoadBalancer
type LoadBalancer struct {
	LoadBanlanceMap   map[string]*LoadBalancerItem
	LoadBanlanceSlice []*LoadBalancerItem
	Locker            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBanlance *load_balance.LoadBalance
	ServiceName  string
	UpdatedAt    time.Time
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBanlanceMap:   map[string]*LoadBalancerItem{},
		LoadBanlanceSlice: []*LoadBalancerItem{},
		Locker:            sync.RWMutex{},
	}
}

func init() {
	LoadBalancerHandler = NewLoadBalancer()
	ServiceManagerHandler.Regist(LoadBalancerHandler)
}


func (lbr *LoadBalancer) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		lbr.GetLoadBalancer(service)
	}
	for _, service := range e.UpdateService {
		lbr.GetLoadBalancer(service)
	}
	newLBSlice := []*LoadBalancerItem{}
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		matched := false
		for _, service := range e.DeleteService {
			if lbrItem.ServiceName == service.Info.ServiceName {
				lbrItem.LoadBanlance.Close()
				matched = true
			}
		}
		if matched {
			delete(lbr.LoadBanlanceMap, lbrItem.ServiceName)
		} else {
			newLBSlice = append(newLBSlice, lbrItem)
		}
	}
	lbr.LoadBanlanceSlice = newLBSlice
}

var balancerMap = map[int]string{
	0: "random",
	1: "round",
	2 :"weight_round",
	3: "consistent_hash",
}

var	loadTypeMap = map[int]string{
	0: "http",
	1: "tcp",
	2: "grpc",
}

func (lbr *LoadBalancer) GetLoadBalancer(service *model.ServiceDetail) (*load_balance.LoadBalance, error) {
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName && lbrItem.UpdatedAt == time.Time(service.Info.UpdateAt) {
			return lbrItem.LoadBanlance, nil
		}
	}
	//fmt.Println("service.Info.LoadBalanceType", service.Info.LoadBalanceType)
	confHandler := load_balance.GetCheckConfigHandler(loadTypeMap[int(service.Info.LoadType)])
	checkConf, err := confHandler(service)
	if err != nil {
		return nil, err
	}
	//fmt.Println("service.Info.LoadBalanceStrategy", service.Info.LoadBalanceStrategy)
	lb := load_balance.LoadBanlanceFactorWithStrategy(load_balance.GetLoadBalanceStrategy(balancerMap[service.LoadBalance.RoundType]), checkConf)
	matched := false
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			matched = true
			lbrItem.LoadBanlance.Close()
			lbrItem.LoadBanlance = lb
			lbrItem.UpdatedAt = time.Time(service.Info.UpdateAt)
		}
	}
	if !matched {
		lbItem := &LoadBalancerItem{
			LoadBanlance: lb,
			ServiceName:  service.Info.ServiceName,
			UpdatedAt:    time.Time(service.Info.UpdateAt),
		}
		lbr.LoadBanlanceSlice = append(lbr.LoadBanlanceSlice, lbItem)
		lbr.Locker.Lock()
		defer lbr.Locker.Unlock()
		lbr.LoadBanlanceMap[service.Info.ServiceName] = lbItem
	}
	return lb, nil
}
