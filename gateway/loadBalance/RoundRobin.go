// 轮训算法
package loadbalance

import (
	"log"
)

type RoundRobinBalance struct {
	Addrs []string `json:"addrs"`
	CurIndex int `json:"cur_index"`
	conf LoadBalanceConf
}

func (r *RoundRobinBalance) Next() string{
	addr := r.Addrs[r.CurIndex]
	r.CurIndex = (r.CurIndex + 1) % len(r.Addrs) 
	return addr
}


func (cs *RoundRobinBalance) Get(addr string) string {
	return cs.Next()
}

func (r *RoundRobinBalance) Add(addrs ...string) {
		r.Addrs = addrs
}

func (cs *RoundRobinBalance) Update() {
	if conf, ok := cs.conf.(*LoadBalanceZkConf); ok {
		// for _, ip := range {
		// 	cs.Add(strings.Split(ip, ",")...)
		// }
		cs.Add(conf.GetConf()...)
	}
	log.Println("更新后队列 ", cs.Addrs)
}

func (cs *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	cs.conf = conf
}