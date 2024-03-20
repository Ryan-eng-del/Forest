package model

import (
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessControl struct {
	AbstractModel
	ServiceInfoID   uint
	OpenAuth          bool  `json:"open_auth" gorm:"column:open_auth;comment:是否开启权限 1=开启" `
	BlackList         string `json:"black_list" gorm:"column:black_list;type:varchar(1000);comment:黑名单ip" `
	WhiteList         string `json:"white_list" gorm:"column:white_list;type:varchar(1000);comment:白名单ip" `
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name;type:varchar(1000);comment:白名单主机	"`
	ClientIPFlowLimit int  `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit;comment:客户端ip限流" `
	ServiceFlowLimit  int  `json:"service_flow_limit" gorm:"column:service_flow_limit;comment:服务端限流" `
}

func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}


func (t *AccessControl) Find(c *gin.Context, tx *gorm.DB) error {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()),mysqlLib.LogicalObjects())
	return query.Find(t).Error
}