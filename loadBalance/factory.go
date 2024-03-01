package loadbalance


type LoadBalance interface {
	Add(...string) 
	Get(string)(string)
	SetConf(LoadBalanceConf)
	Update()
}

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbTYpe LbType) LoadBalance {
	var lb LoadBalance = nil

	switch lbTYpe {
	case LbRandom:
		lb = &RandomBalance{}

	case LbRoundRobin:
		lb = &RoundRobinBalance{}

	case LbWeightRoundRobin:
		lb = &WeightRoundRobinBalance{}

	case LbConsistentHash:
		lb = NewConsistentHashBalance(5, nil)
	default:
		lb = &RoundRobinBalance{}
	}
	return lb
}


func LoadBalanceWithConfFactory(lbTYpe LbType, mconf LoadBalanceConf) LoadBalance {
	var lb LoadBalance = nil

	switch lbTYpe {
	case LbRandom:
		lb = &RandomBalance{}
		initLoadBalance(lb, mconf)
	case LbRoundRobin:
		lb = &RoundRobinBalance{}
		initLoadBalance(lb, mconf)
	case LbWeightRoundRobin:
		lb = &WeightRoundRobinBalance{}
		initLoadBalance(lb, mconf)
	case LbConsistentHash:
		lb = NewConsistentHashBalance(5, nil)
		initLoadBalance(lb, mconf)
	default:
		initLoadBalance(lb, mconf)
	}
	return lb
}


func initLoadBalance(lb LoadBalance, mconf LoadBalanceConf) {
	lb.SetConf(mconf)
	mconf.WatchConf(lb)
	mconf.Attach(lb)
	lb.Update()
}