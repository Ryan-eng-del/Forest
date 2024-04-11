package load_balance

import (
	"errors"
	"log"
	"strconv"
	"time"
)

const (
	maxFails = 2;
	failTimeout = time.Hour * 1
)

// 加权轮训举例
// 4 2 1   4 2 1   7

// -3 2 1  4 2 1 -> a
// 1 -3 2   4 2 1 -> b
// -2 -1 3  4 2 1 -> a
// 5  -1 4   4 2 1 -> b
// 2 1 5    4 2 1 -> a
// 6 3 -1  4 2 1 -> c
// 3 5 0 4 2 1 -> a
type WeightRoundRobinStrategy struct {
	CurIndex int
	rss      []*WeightNode
	rsw      []int
	Conf     LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
	// 在 failTimeout 时间内，最大失败次数
	MaxFails int
	// 指定时间段内的最大失败次数
	FailTimeout time.Duration
	// failTimeout 时间内，节点失败的时间
	FailTimes []time.Time
}

func (s *WeightNode) RefreshTime() {
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

func (r *WeightRoundRobinStrategy) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinStrategy) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]

		if w.MaxFails <= 0 {
			w.MaxFails = maxFails - len(w.FailTimes)
			w.RefreshTime()
			if w.MaxFails <= 0 {
				log.Println("跳过该节点", w.addr)
				continue
			}
		}

		total += w.effectiveWeight
		w.currentWeight += w.effectiveWeight
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}


	if best == nil {
		return ""
	}
	best.currentWeight -= total
	return best.addr
}


func (r *WeightRoundRobinStrategy) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *WeightRoundRobinStrategy) GetAll() ([]string, error) {
	iplist := []string{}
	for _, item := range r.rss {
		iplist = append(iplist, item.addr)
	}
	return iplist, nil
}

func (r *WeightRoundRobinStrategy) RemoveAll() error {
	r.rss = []*WeightNode{}
	r.rsw = []int{}
	return nil
}

func (r *WeightRoundRobinStrategy) Callback(addr string, flag bool) {
	for _, node := range r.rss {
			if node.addr == addr {
				if flag {
					if node.effectiveWeight < node.weight {
						node.effectiveWeight++
					}
				} else {
					node.effectiveWeight--
					node.FailTimes = append(node.FailTimes, time.Now())
					node.RefreshTime()
					node.MaxFails = maxFails - len(node.FailTimes)
				}
			}
	}
}


func init() {
	RegisterLoadBalanceStrategyHandler("weight_round", func() LoadBalanceStrategy {
		return &WeightRoundRobinStrategy{}
	})
}


