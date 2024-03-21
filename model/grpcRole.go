package model

import (
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type GrpcRule struct {
	ID        uint `json:"id" gorm:"primary_key"`
	ServiceInfoID      uint  `json:"service_id" gorm:"comment:服务id"`
	Service *Service `json:"service,omitempty" gorm:"foreignKey:ServiceInfoID;references:ID"`
	Port      int   `json:"port" gorm:"comment:端口;not null;"`
	HeaderTransfor string `json:"header_transfor" gorm:"comment:header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue;type:varchar(5000);not null;"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *GrpcRule) Find(c *gin.Context, tx *gorm.DB) (*GrpcRule, error) {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()))
	result:= query.Find(t, "service_id = ?", t.ServiceInfoID).Error
	if t.ID == 0 {
		return nil, result
	}
	return t, result
}