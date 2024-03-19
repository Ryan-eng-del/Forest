package model

type App struct {
	AbstractModel
	AppID string `json:"app_id" gorm:"type:varchar(255)"`
	Name string `json:"name" gorm:"type:varchar(255)"`
	Secret string `json:"secret" gorm:"type:varchar(255)"`
	WhiteIps  string    `json:"white_ips"  gorm:"type:varchar(255);comment:ip白名单,支持前缀匹配"`
	Qpd       uint     `json:"qpd" gorm:"comment:日请求量限制"`
	Qps       uint     `json:"qps" gorm:"comment:每秒请求量限制"`
} 

func (t *App) TableName() string {
	return "gateway_app"
}
