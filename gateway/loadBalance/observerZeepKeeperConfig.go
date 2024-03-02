package loadbalance

import (
	"go-gateway/gateway/middleware/serverDiscovery/observer"
	zookeeper "go-gateway/gateway/middleware/serverDiscovery/zooKeeper"
	"log"
)

type LoadBalanceObserver struct {
	zkConf *LoadBalanceZkConf
}

func NewLoadBalanceObserver (zkConf *LoadBalanceZkConf) (*LoadBalanceObserver) {
	return &LoadBalanceObserver{
		zkConf: zkConf,
	}
}

func (o *LoadBalanceObserver) Update() {
	log.Printf("Observer Receive Update: %+v", o.zkConf.activeList)
}

type LoadBalanceConf interface {
	Attach (o observer.Observer)
	GetConf () []string
	WatchConf(lb LoadBalance)
	UpdateConf(conf []string)
}

type LoadBalanceZkConf struct {
	observers []observer.Observer
	path string
	zkHost []string
	confIPWeight map[string]string
	activeList []string
	format string
}

func (s *LoadBalanceZkConf) UpdateConf (conf []string) {
	s.activeList = conf
	log.Println(s.observers, "Observers")
	for _, o := range s.observers {
		o.Update()
	}
}

func (s *LoadBalanceZkConf) Attach (o observer.Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceZkConf) GetConf () []string {
	return s.activeList
}

func (s *LoadBalanceZkConf) WatchConf(lb LoadBalance) {
	zk := zookeeper.NewZkManager(s.zkHost)
	zk.GetConnection()
	chanList, chanErr := zk.WatchServerListByPath(s.path)

	go func () {
		defer zk.Close()
		for {
			select {
			case changeErr := <- chanErr:
				log.Println("changeErr:", changeErr)
			case changeList := <- chanList:
				log.Println("watch node changed", changeList)
				s.UpdateConf(changeList)
			}
		} 
	}()
}

func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string ) (*LoadBalanceZkConf, error) {
	zk := zookeeper.NewZkManager(zkHosts)
	zk.GetConnection()

	defer zk.Close()
	zList, err := zk.GetServerListPath(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	mConf := &LoadBalanceZkConf{
		format: format,
		activeList: zList,
		confIPWeight: conf,
		zkHost: zkHosts,
		path: path,
	}

	return mConf, nil
}


// example observer | zkconf
func Example() {
	conf, err := NewLoadBalanceZkConf("%s", zookeeper.NodeName, []string{"localhost:2181", "localhost:2182", "localhost:2183"}, map[string]string{
		"localhost:8002": "10",
		"localhost:8001": "20",
	})

	if err != nil {
		log.Println(err)
		return
	}

	observer := NewLoadBalanceObserver(conf)
	conf.Attach(observer)
	conf.UpdateConf([]string{"localhost:8087"})
	select{}

}