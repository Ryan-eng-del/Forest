package main

import "math/rand"

type RandomBalance struct {
	Addrs []string `json:"addrs"`
	CurIndex int `json:"cur_index"`
}


func (r *RandomBalance) Add(addrs ...string) {
	if (len(addrs) > 0) {
		r.Addrs = append(r.Addrs, addrs...)
	}
} 


func (r *RandomBalance) Next() string {
	r.CurIndex = rand.Intn(len(r.Addrs))
	return r.Addrs[r.CurIndex]
}