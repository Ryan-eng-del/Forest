package model

import (
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceInfoID      uint  `json:"service_id" gorm:"comment:服务id;column:service_id"`
	Service *Service `json:"service,omitempty" gorm:"foreignKey:ServiceInfoID;references:ID"`
	Port      int   `json:"port" gorm:"comment:端口"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TcpRule) Find(c *gin.Context, tx *gorm.DB) (*TcpRule, error) {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()))
	result := query.Find(t, "service_id = ?", t.ServiceInfoID).Error
	if t.ID == 0 {
		return nil, result
	}
	return t, result
}


func (t *TcpRule) FindMust(c *gin.Context, tx *gorm.DB, queryStruct *TcpRule) (*TcpRule, error) {
	out := &TcpRule{}
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()))
	result := query.Where(queryStruct).First(t).Error
	return out, result
}

func (t *TcpRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}
