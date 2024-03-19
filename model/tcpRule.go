package model

import (
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceInfoID uint `json:"service_id"`
	Port      int   `json:"port" gorm:"comment:端口"`
}

func (t *TcpRule) TableName() string {
	return "gateway_tcp_rule"
}


func (t *TcpRule) Find(c *gin.Context, tx *gorm.DB) error {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()),mysqlLib.LogicalObjects())
	return query.Find(t).Error
}