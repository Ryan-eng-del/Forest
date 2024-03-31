package model

import (
	appDto "go-gateway/dto/app"
	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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


func (t *App) AppList (ctx *gin.Context, tx *gorm.DB, params *appDto.APPListInput) ([]App, int64, error){
	var list []App
	var count int64

	query := tx.Scopes(mysqlLib.WithContextAndTable(ctx, t.TableName()), mysqlLib.LogicalObjects(),mysqlLib.IDDesc())

	if params.Info != "" {
		query = query.Where("(name like ? or app_id like ?)", "%" + params.Info + "%",  "%" + params.Info + "%")
	}

	err := query.Scopes(mysqlLib.Paginate(params.PageNo, params.PageSize)).Find(&list).Count(&count).Error;

	if err != nil {
		return nil, 0, err
	}
	
	return list, count, nil
}