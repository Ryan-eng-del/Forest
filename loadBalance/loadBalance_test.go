package main

import (
	"log"
	"testing"
)


func TestRoundRobin( t *testing.T ) {
	r := &RoundRobinBalance{}
	r.Add("127.0", "127.1", "127.2")

	for i :=  0; i < 10; i++ {
		log.Println(r.Next())
	}
}