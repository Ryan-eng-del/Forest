package load_balance

import (
	"fmt"
	"go-gateway/model"
	"log"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)


type LoadBalanceCheckConf struct {
	observers    []Observer
	confIpWeight map[string]string
	activeList   []string
	format       string //单条数据格式 http://%s，%s方便替换成ip地址
	name         string
	closeChan    chan bool
}


func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}


func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceCheckConf) WatchConf() {
	go func(){
		confIpErrNum := map[string]int{}
		log.Printf("checking [%s] config_list:%v active_list:%v\n", s.name, s.confIpWeight, s.activeList)
		OUTFOR:
			for {
				select {
				case <- s.closeChan:
					break OUTFOR
				default:
				}
				changedList := []string{}
				for rs, _ := range s.confIpWeight {
					conn, err := net.DialTimeout("tcp", rs, time.Duration(DefaultCheckTimeout)*time.Second)
					if err == nil {
						conn.Close()
						if _, ok := confIpErrNum[rs]; ok {
							confIpErrNum[rs] = 0
						}
					}

					if err != nil {
						if _, ok := confIpErrNum[rs]; ok {
							confIpErrNum[rs] += 1
						} else {
							confIpErrNum[rs] = 1
						}
					}

					if confIpErrNum[rs] < DefaultCheckMaxErrNum {
						changedList = append(changedList, rs)
					}
				}

				sort.Strings(changedList)
				sort.Strings(s.activeList)

				if !reflect.DeepEqual(changedList, s.activeList) {
					log.Printf("checking [%s] changed config_list:%v changed_list:%v\n", s.name, s.confIpWeight, changedList)
					s.UpdateConf(changedList)
				}
				time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
			}
	}()
}
func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

func (s *LoadBalanceCheckConf) CloseWatch() {
	s.closeChan <- true
	close(s.closeChan)
}

func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceCheckConf(service *model.ServiceDetail) (LoadBalanceConf, error) {
	ipList := service.LoadBalance.GetIPListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	confIpWeight := map[string]string{}
	for i, ip := range ipList {
		weight := weightList[i]
		confIpWeight[ip] = weight
	}
	schema := "http://"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}

	mConf := &LoadBalanceCheckConf{
		name:         service.Info.ServiceName,
		format:       fmt.Sprintf("%s%s", schema, "%s"),
		activeList:   ipList,
		confIpWeight: confIpWeight,
		closeChan:    make(chan bool, 1)}

	mConf.WatchConf()
	return mConf, nil
}

func init() {
	RegisterCheckConfigHandler("upstream_config", NewLoadBalanceCheckConf)
}
