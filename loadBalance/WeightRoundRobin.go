// 加权轮训算法
package loadbalance

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	maxFails = 3;
	failTimeout = time.Second * 4
)

type WeightRoundRobinBalance struct {
	ServeNodes []*ServerNode
	CurIndex int
	conf LoadBalanceConf
}

type WeightParameters struct {
	addr string
	weight int
}



func (r *WeightRoundRobinBalance) Callback(addr string, flag bool) {
	for _, node := range r.ServeNodes {
			if node.Addr == addr {
				if flag {
					if node.EffectiveWeight < node.Weight {
						node.EffectiveWeight++
					}
				} else {
					node.EffectiveWeight--
					node.FailTimes = append(node.FailTimes, time.Now())
					node.RefreshTime()
					node.MaxFails = maxFails - len(node.FailTimes)
				}
			}
	}
}

// func (r *WeightRoundRobinBalance) Add(addrs ...WeightParameters) {
// 	if (len(addrs) > 0) {
// 		for _, addr := range addrs {
// 			node := &ServerNode {
// 				Addr: addr.addr,
// 				Weight: addr.weight,
// 				EffectiveWeight: addr.weight,
// 				MaxFails: maxFails,
// 				FailTimeout: failTimeout,
// 			}
// 			r.ServeNodes = append(r.ServeNodes, node)
// 		}
// 	}
// }


func (r *WeightRoundRobinBalance) Add(addrs ...string) {
	if (len(addrs) > 0) {
		for _, addr := range addrs {
			node := &ServerNode {
				Addr: addr,
				Weight: 10,
				EffectiveWeight: 10,
				MaxFails: maxFails,
				FailTimeout: failTimeout,
			}
			r.ServeNodes = append(r.ServeNodes, node)
		}
	}
}

func (cs *WeightRoundRobinBalance) Get(addr string) string {
	return ""
}


func (cs *WeightRoundRobinBalance) Update() {
	if conf, ok := cs.conf.(*LoadBalanceZkConf); ok {
		for _, ip := range conf.GetConf() {
			is := strings.Split(ip, ",")
			cs.Add(is...)
		}
	}
}


func (cs *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	cs.conf = conf
}

func (r *WeightRoundRobinBalance) Next() (string, error){
	var effectiveTotal = 0
	var maxNode *ServerNode
	var index = 0

	for i := 0; i < len(r.ServeNodes); i++ {
		w := r.ServeNodes[i]

		if w.MaxFails <= 0 {
			w.MaxFails = maxFails - len(w.FailTimes)
			w.RefreshTime()
			if w.MaxFails <= 0 {
				fmt.Println("跳过该节点", w.Addr)
				continue
			}
		}

		w.CurrentWeight += w.EffectiveWeight
		effectiveTotal += w.EffectiveWeight

		if maxNode == nil || w.CurrentWeight > maxNode.CurrentWeight {
				index = i
				maxNode = w
		}
	}

	if maxNode != nil {
		r.CurIndex = index
		maxNode.CurrentWeight -= effectiveTotal
		return maxNode.Addr, nil
	}

	return "", errors.New("maxNode not found")
}

type ServerNode struct {
	Addr string
	// 初始化权重
	Weight int
	// 节点当前权重
	CurrentWeight int
	// 节点有效权重
	EffectiveWeight int
	// 在 failTimeout 时间内，最大失败次数
	MaxFails int
	// 指定时间段内的最大失败次数
	FailTimeout time.Duration
	// failTimeout 时间内，节点失败的时间
	FailTimes []time.Time
}



func (s *ServerNode) RefreshTime() {
	var i = 0
	var now = time.Now()
	for i = 0; i < s.MaxFails; i++ {
		t := s.FailTimes[i]
		deadline := t.Add(failTimeout)
		if deadline.After(now) || deadline.Equal(now) {
			break
		}
	}

	s.FailTimes = s.FailTimes[i:]
}
