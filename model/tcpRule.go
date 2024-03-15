package model

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceInfoID int64 `json:"service_id" gorm:"comment:服务id"`
	Port      int   `json:"port" gorm:"comment:端口"`
}

func (t *TcpRule) TableName() string {
	return "gateway_tcp_rule"
}