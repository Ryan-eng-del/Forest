package model

import (
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HttpRule struct {
	ID             uint  `json:"id" gorm:"primary_key"`
	ServiceInfoID      uint  `json:"service_id" gorm:"comment:服务id"`
	Service *Service `json:"service,omitempty" gorm:"foreignKey:ServiceInfoID;references:ID"`
	RuleType       int    `json:"rule_type" gorm:"comment:匹配类型 domain=域名, url_prefix=url前缀"`
	Rule           string `json:"rule" gorm:"column:rule;type:varchar(255);" comment:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps      bool    `json:"need_https" gorm:"column:need_https" comment:"type=支持https 1=支持"`
	NeedWebsocket  bool    `json:"need_websocket" gorm:"column:need_websocket" comment:"启用websocket 1=启用"`
	NeedStripUri   bool    `json:"need_strip_uri" gorm:"column:need_strip_uri" comment:"启用strip_uri 1=启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite;type:varchar(5000);" comment:"url重写功能，每行一个	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor;type:varchar(5000);" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
}

func (t *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *HttpRule) Find(c *gin.Context, tx *gorm.DB) (*HttpRule, error) {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()))
	result := query.Find(t, "service_id = ?", t.ServiceInfoID).Error
	if t.ID == 0 {
		return nil, result
	}
	return t, result
}