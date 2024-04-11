package model

import (
	mysqlLib "go-gateway/lib/mysql"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoadBalance struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	ServiceInfoID      uint  `json:"service_id" gorm:"comment:服务id;column:service_id"`
	Service *Service `json:"service,omitempty" gorm:"foreignKey:ServiceInfoID;references:ID"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method;type:int;size:14" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  uint    `json:"check_timeout" gorm:"column:check_timeout;unsigned;size:15;not null;" description:"check超时时间"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}


func (t *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(t.WeightList, ",")
}


func (t *LoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}


func (t *LoadBalance) Create(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}

func (t *LoadBalance) Find(c *gin.Context, tx *gorm.DB) (*LoadBalance, error) {
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()))

	
	result := query.Find(t, "service_id = ?", t.ServiceInfoID).Error

	if t.ID == 0 {
		return nil, result
	}
	return t, result
}