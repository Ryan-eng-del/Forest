package model


type Service struct {
	AbstractModel
	LoadType    int       `json:"load_type" gorm:"comment:负载类型 0=http 1=tcp 2=grpc;"`
	ServiceName string    `json:"service_name" gorm:"comment:服务名称;type:varchar(255);"`
	ServiceDesc string    `json:"service_desc" gorm:"comment:服务描述;type:varchar(255);"`
	AccessControl AccessControl  `gorm:"foreignKey:ServiceInfoID;references:ID"`
	HttpRule HttpRule `gorm:"foreignKey:ServiceInfoID;references:ID"`
	GrpcRole GrpcRule `gorm:"foreignKey:ServiceInfoID;references:ID"`
	LoadBalance LoadBalance `gorm:"foreignKey:ServiceInfoID;references:ID"`
}


func (t *Service) TableName() string {
	return "gateway_service"
}
