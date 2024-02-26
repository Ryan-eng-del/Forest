package main

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)


type ConsistentHashBalance struct {
	hash Hash
	hashKeys UInt32Slice
	hashMap map[uint32]string
	replicas int
	mux sync.Mutex
}

func (cs *ConsistentHashBalance) Get(addr string) string {
	hash := cs.hash([]byte(addr))

	index := sort.Search(len(cs.hashKeys), func (i int) bool {
		return cs.hashKeys[i] >= hash
	})

	if index == len(cs.hashKeys) {
		index = 0
	}
	
	return cs.hashMap[cs.hashKeys[index]]
}
func (cs *ConsistentHashBalance) Add(addrs ...string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	for _, addr := range addrs {
		for i := 0; i < cs.replicas; i++ {
			hashValue := cs.hash([]byte(strconv.Itoa(i) + addr))
			cs.hashKeys = append(cs.hashKeys, hashValue)
			cs.hashMap[hashValue] = addr
		}
	}

	sort.Sort(cs.hashKeys)
}

func NewConsistentHashBalance (replicas int, hash Hash) *ConsistentHashBalance {
	cs := &ConsistentHashBalance{
		replicas: replicas,
		hash: hash,
	}

	if cs.hash == nil {
		cs.hash = crc32.ChecksumIEEE
	}

	return cs
}


type UInt32Slice []uint32

func (a UInt32Slice) Len() int {
	return len(a)
}

func (a UInt32Slice) Less (i, j int) bool {
	return a[i] < a[j]
}

func (a UInt32Slice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}	


type Hash func (data []byte) uint32
