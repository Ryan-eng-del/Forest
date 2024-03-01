package loadbalance

import (
	"math/rand"
	"strings"
)

type RandomBalance struct {
	Addrs []string `json:"addrs"`
	CurIndex int `json:"cur_index"`
	conf LoadBalanceConf
}

func (r *RandomBalance) Add(addrs ...string) {
	if (len(addrs) > 0) {
		r.Addrs = append(r.Addrs, addrs...)
	}
} 

func (cs *RandomBalance) Get(addr string) string {
	return ""
}


func (r *RandomBalance) Next() string {
	r.CurIndex = rand.Intn(len(r.Addrs))
	return r.Addrs[r.CurIndex]
}

func (cs *RandomBalance) Update() {
	if conf, ok := cs.conf.(*LoadBalanceZkConf); ok {
		for _, ip := range conf.GetConf() {
			cs.Add(strings.Split(ip, ",")...)
		}
	}
}

func (cs *RandomBalance) SetConf(conf LoadBalanceConf) {
	cs.conf = conf
}