package handler

import (
	serviceDto "go-gateway/dto/service"
	lib "go-gateway/lib/conf"
	libLog "go-gateway/lib/log"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var ServiceManagerHandler *ServiceManager = NewServiceManager()

//通知事件
type ServiceEvent struct {
	DeleteService []*model.ServiceDetail
	AddService    []*model.ServiceDetail
	UpdateService []*model.ServiceDetail
}

//观察者接口
type ServiceObserver interface {
	Update(*ServiceEvent)
}

//被观察者接口
type ServiceSubject interface {
	Regist(ServiceObserver)
	Deregist(ServiceObserver)
	Notify(*ServiceEvent)
}


type ServiceManager struct {
	Locker sync.RWMutex
	err error
	ServiceSlice []*model.ServiceDetail
	UpdateAt time.Time
	ServiceMap map[string]*model.ServiceDetail
	Observers map[ServiceObserver]bool
	sync.RWMutex
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap: map[string]*model.ServiceDetail{},
		ServiceSlice: []*model.ServiceDetail{},
		Observers:    map[ServiceObserver]bool{},
	}
}

func (s *ServiceManager) Regist(ob ServiceObserver) {
	s.Lock()
	defer s.Unlock()
	s.Observers[ob] = true
}

func (s *ServiceManager) Deregist(ob ServiceObserver) {
	s.Lock()
	defer s.Unlock()
	delete(s.Observers, ob)
}

func (s *ServiceManager) Notify(e *ServiceEvent) {
	s.RLock()
	defer s.RUnlock()
	for ob := range s.Observers {
		ob.Update(e)
	}
}

func (m *ServiceManager) LoadService() *ServiceManager {
	sm := NewServiceManager()

	defer func () {
		if sm.err != nil {
			log := libLog.GetLogger()
			log.Error("load service config error:%v", sm.err)
		}
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	tx, err := libMysql.GetGormPool("default")
	if err != nil {
		sm.err = err
		return sm
	}

	params := serviceDto.ServiceListInput{PageNo: 1, PageSize: 9999999}

	serviceInfo := model.Service{}

	services, _, err := serviceInfo.PageList(c, tx, &params)

	if err != nil {
		sm.err = err
		return sm
	}

	for _, service := range services {

		serviceDetail, err := service.ServiceDetail(c, tx);
		if err != nil {
			sm.err = err
			return sm
		}

		m.ServiceMap[service.ServiceName] = serviceDetail
		m.ServiceSlice = append(m.ServiceSlice, serviceDetail)

		if m.UpdateAt.Unix() < service.UpdateAt.Unix() {
			m.UpdateAt = time.Time(service.UpdateAt)
		}
	}
	
	return m
}

func (s *ServiceManager) LoadAndWatch() error {
	libLog.NewSingleLoggerDefault()
	log := libLog.GetLogger()
	log.Info("watching load service config from resource")

	ns := s.LoadService()

	if ns.err != nil {
		return ns.err
	}

	s.ServiceSlice = ns.ServiceSlice
	s.ServiceMap = ns.ServiceMap
	s.UpdateAt = ns.UpdateAt

	e := &ServiceEvent{AddService: ns.ServiceSlice}
	s.Notify(e)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			ns := s.LoadService()
			if ns.err != nil {
				log := libLog.GetLogger()
				log.Error("load service config error:%v", ns.err)
			}

			if ns.UpdateAt != s.UpdateAt || len(ns.ServiceSlice) != len(s.ServiceSlice)  {
				e := &ServiceEvent{}
				for _, service := range s.ServiceSlice {
					matched := false
					for _, newService := range ns.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName {
							matched = true
						}
					}
					if !matched {
						e.DeleteService = append(e.DeleteService, service)
					}
				}

				for _, newService := range ns.ServiceSlice {
					matched := false
					for _, service := range s.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName {
							matched = true
						}
					}
					if !matched {
						e.AddService = append(e.AddService, newService)
					}
				}

				for _, newService := range ns.ServiceSlice {
					matched := false
					for _, service := range s.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName && service.Info.UpdateAt != newService.Info.UpdateAt {
							matched = true
						}
					}
					if matched {
						e.UpdateService = append(e.UpdateService, newService)
					}
				}

				for _, item := range e.DeleteService {
					log.Info("found config delete service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))
				}

				for _, item := range e.AddService {
					log.Info("found config add service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))
				}

				for _, item := range e.UpdateService {
					log.Info("found config update service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))
				}

				s.ServiceSlice = ns.ServiceSlice
				s.ServiceMap = ns.ServiceMap
				s.UpdateAt = ns.UpdateAt
				s.Notify(e)
			}
		}
	}()

	return ns.err
}
