package model

import (
	serviceDto "go-gateway/dto/service"
	mysqlLib "go-gateway/lib/mysql"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type Service struct {
	AbstractModel
	LoadType    int       `json:"load_type" gorm:"comment:负载类型 0=http 1=tcp 2=grpc;"`
	ServiceName string    `json:"service_name" gorm:"comment:服务名称;type:varchar(255);"`
	ServiceDesc string    `json:"service_desc" gorm:"comment:服务描述;type:varchar(255);"`
	AccessControl *AccessControl  `gorm:"foreignKey:ServiceInfoID;references:ID"`
	HttpRule *HttpRule `gorm:"foreignKey:ServiceInfoID;references:ID"`
	GrpcRole *GrpcRule `gorm:"foreignKey:ServiceInfoID;references:ID"`
	LoadBalance *LoadBalance `gorm:"foreignKey:ServiceInfoID;references:ID"`
	TcpRole *TcpRule `gorm:"foreignKey:ServiceInfoID;references:ID"`
}



func (t *Service) PageList(c *gin.Context, tx *gorm.DB, params *serviceDto.ServiceListInput) ([]Service, int64, error) {
	var total int64
	list := []Service{}
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()), mysqlLib.LogicalObjects(), mysqlLib.Paginate(params.PageNo, params.PageSize), mysqlLib.IDDesc())
	
	if params.Info != "" {
		query = query.Where("service_name like ? or service_desc like ?", "%" + params.Info + "%", "%" + params.Info + "%")
	}
	err := query.Find(&list).Count(&total).Error
	return list, total, err
}
 
func (t *Service) ServiceDetail (c *gin.Context, tx *gorm.DB) (error) {

	httpRule := &HttpRule{ServiceInfoID: t.ID}
	err := httpRule.Find(c, tx)

	if err != nil {
		return err
	}

	grpcRule := &GrpcRule{ServiceInfoID: int64(t.ID)}
	err = grpcRule.Find(c, tx)

	if err != nil {
		return err
	}

	tcpRule := &TcpRule{ServiceInfoID: t.ID}
	err = tcpRule.Find(c, tx)

	if err != nil {
		return err
	}
	
	accessControl := &AccessControl{ServiceInfoID: t.ID}
	err = accessControl.Find(c, tx)

	if err != nil {
		return err
	}

	log.Println(httpRule, "HttpRule")
	log.Println(tcpRule, "tcpRule")
	log.Println(grpcRule, "grpcRule")
	log.Println(accessControl, "accessRule")

	t.HttpRule = httpRule
	t.GrpcRole = grpcRule
	t.AccessControl = accessControl
	t.TcpRole = tcpRule
	return nil
}


func (t *Service) TableName() string {
	return "gateway_service"
}
