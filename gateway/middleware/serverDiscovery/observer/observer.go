package observer


type Observer interface {
	Update()
}

type ConcreteSubject struct {
	observers []Observer
	conf []string
	name string
}

func (s *ConcreteSubject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *ConcreteSubject) Notify() {
	for _, o := range s.observers {
		o.Update()
	}
}

func (s *ConcreteSubject) GetConf() []string {
	return s.conf
}

func (s *ConcreteSubject) UpdateConf(conf []string) {
	s.conf = conf
	for _, o := range s.observers {
		o.Update()
	}
}

func NewConCreateSubject(name string) *ConcreteSubject {
	return &ConcreteSubject{
		name: name,
	}
}
