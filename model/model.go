package model

import (
	"errors"
	lib "go-gateway/lib/const"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)


type AbstractModel struct {
	ID        uint `gorm:"primarykey;comment:自增主键"`
	CreateAt public.LocalTime  `json:"create_at" gorm:"comment:创建时间;autoCreateTime"`
	UpdateAt public.LocalTime `json:"update_at" gorm:"comment:更新时间;autoUpdateTime"`
	IsDelete int8  `gorm:"comment:是否删除"`
}

type ServiceDetail struct {
	Info          *Service   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule,omitempty" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule,omitempty" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule,omitempty" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance,omitempty" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control,omitempty" description:"access_control"`
}


func GetServiceDetailFromGinContext(c *gin.Context) (*ServiceDetail, error) {
	serverInterface, ok := c.Get(lib.ServiceDetailContextKey)
	if !ok {
		return nil, errors.New("service not found")
	}
	return serverInterface.(*ServiceDetail), nil
}
