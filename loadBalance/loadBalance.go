package load_balance

type LoadBalanceStrategy interface {
	Add(...string) error
	RemoveAll() error
	GetAll() ([]string, error)
	Get(string) (string, error)
}

var LoadBalanceStrategyHandlerMap map[string]LoadBalanceStrategyHandler

type LoadBalanceStrategyHandler func() LoadBalanceStrategy

func RegisterLoadBalanceStrategyHandler(name string, handler LoadBalanceStrategyHandler) {
	if LoadBalanceStrategyHandlerMap == nil {
		LoadBalanceStrategyHandlerMap = map[string]LoadBalanceStrategyHandler{}
	}
	LoadBalanceStrategyHandlerMap[name] = handler
}

func GetLoadBalanceStrategy(name string) LoadBalanceStrategy {
	if LoadBalanceStrategyHandlerMap == nil {
		return nil
	}
	handler := LoadBalanceStrategyHandlerMap[name]
	return handler()
}
