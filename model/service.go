package model

import (
	dashboardDto "go-gateway/dto/dashboard"
	serviceDto "go-gateway/dto/service"
	mysqlLib "go-gateway/lib/mysql"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	// gorm.Model
	AbstractModel
	LoadType    public.LoadType       `json:"load_type" gorm:"comment:负载类型 0=http 1=tcp 2=grpc;"`
	ServiceName string    `json:"service_name" gorm:"comment:服务名称;type:varchar(255);"`
	ServiceDesc string    `json:"service_desc" gorm:"comment:服务描述;type:varchar(255);"`
}


func (t *Service) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dashboardDto.DashServiceStatItemOutput, error) {
	list := []dashboardDto.DashServiceStatItemOutput{}
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()), mysqlLib.LogicalObjects())

	if err := query.Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
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
 
func (t *Service) ServiceDetail (c *gin.Context, tx *gorm.DB)  (*ServiceDetail, error) {
	httpInfo := &HttpRule{}
	httpRule, err := httpInfo.Find(c, tx, &HttpRule{ServiceInfoID: t.ID})

	// 外键查询
	// log.Printf("httpRule Preload Before %+v",httpRule)
	// tx.Table(httpRule.TableName()).Preload("Service").Find(&httpRule)
	// log.Printf("httpRule Preload After %+v", httpRule)
	if err != nil {
		return nil, err
	}

	grpcRule := &GrpcRule{ServiceInfoID: t.ID}
	grpcRule, err = grpcRule.Find(c, tx)

	if err != nil {
		return nil, err
	}

	tcpRule := &TcpRule{ServiceInfoID: t.ID}
	tcpRule, err = tcpRule.Find(c, tx)

	if err != nil {
		return nil, err
	}
	
	accessControl := &AccessControl{ServiceInfoID: t.ID}
	accessControl,err = accessControl.Find(c, tx)

	if err != nil {
		return nil, err
	}

	loadBalance := &LoadBalance{ServiceInfoID: t.ID}
	loadBalance, err = loadBalance.Find(c, tx)

	if err != nil {
		return nil, err
	}

	detail := &ServiceDetail{
		Info: t,
		HTTPRule: httpRule,
		GRPCRule: grpcRule,
		AccessControl: accessControl,
		TCPRule: tcpRule,
		LoadBalance: loadBalance,
	}

	return detail, nil
}


func (t *Service) TableName() string {
	return "gateway_service_info"
}


func (t *Service) Find(c *gin.Context, tx *gorm.DB, queryStruct *Service) (*Service, error) {
	out := &Service{}
	query := tx.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()), mysqlLib.LogicalObjects())
	return out, query.Where(queryStruct).First(out).Error
}


func (t *Service) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}

func (t *Service) Create(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Create(t).Error
}