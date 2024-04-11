package handler

import (
	"go-gateway/model"
	"net"
	"net/http"
	"sync"
	"time"
)

var TransportorHandler *Transportor

type Transportor struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
	UpdateAt    time.Time
}


func NewTransportor() *Transportor {
	return &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}


func init() {
	TransportorHandler = NewTransportor()
	ServiceManagerHandler.Regist(TransportorHandler)
}

func (t *Transportor) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		t.GetTrans(service)
	}
	for _, service := range e.UpdateService {
		t.GetTrans(service)
	}
	newSlice := []*TransportItem{}
	for _, tItem := range t.TransportSlice {
		matched := false
		for _, service := range e.DeleteService {
			if tItem.ServiceName == service.Info.ServiceName {
				matched = true
			}
		}
		if matched {
			delete(t.TransportMap, tItem.ServiceName)
		} else {
			newSlice = append(newSlice, tItem)
		}
	}
	t.TransportSlice = newSlice
}


func (t *Transportor) GetTrans(service *model.ServiceDetail) (*http.Transport, error) {

	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName && transItem.UpdateAt == time.Time(service.Info.UpdateAt) {
			return transItem.Trans, nil
		}
	}

	connectTimeout := service.LoadBalance.UpstreamConnectTimeout
	headerTimeout := service.LoadBalance.UpstreamHeaderTimeout
	idleNum := service.LoadBalance.UpstreamMaxIdle
	idConnTimeout := service.LoadBalance.UpstreamIdleTimeout

	
	if connectTimeout == 0 {
		connectTimeout = 30
	}

	if headerTimeout == 0 {
		service.LoadBalance.UpstreamHeaderTimeout = 30
	}

	if idleNum == 0 {
		idleNum = 100
	}

	if idConnTimeout == 0 {
		idConnTimeout = 90
	}

	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(connectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext, //3次握手超时设置
		ForceAttemptHTTP2:     true,
		IdleConnTimeout: time.Duration(idConnTimeout)*time.Second,
		MaxIdleConns:          idleNum,
		WriteBufferSize:       1 << 18, //256m
		ReadBufferSize:        1 << 18, //256m
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(headerTimeout) * time.Second, //请求响应超时
	}

	// 能走到这里，说明服务有新增/更新
	matched := false

	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			matched = true
			transItem.Trans = trans
			transItem.UpdateAt = time.Time(service.Info.UpdateAt)
		}
	}
	
	if !matched {
		transItem := &TransportItem{
			Trans:       trans,
			ServiceName: service.Info.ServiceName,
			UpdateAt:    time.Time(service.Info.UpdateAt),
		}
		t.TransportSlice = append(t.TransportSlice, transItem)
		t.Locker.Lock()
		defer t.Locker.Unlock()
		t.TransportMap[service.Info.ServiceName] = transItem
	}

	return trans, nil
}

