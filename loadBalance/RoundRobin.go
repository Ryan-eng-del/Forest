// 轮训算法
package main

type RoundRobinBalance struct {
	Addrs []string `json:"addrs"`
	CurIndex int `json:"cur_index"`
}

func (r *RoundRobinBalance) Next() string{
	addr := r.Addrs[r.CurIndex]
	r.CurIndex = (r.CurIndex + 1) % len(r.Addrs) 
	return addr
}

func (r *RoundRobinBalance) Add(addrs ...string) {
	if len(addrs) > 0 {
		r.Addrs = append(r.Addrs, addrs...)
	}
}