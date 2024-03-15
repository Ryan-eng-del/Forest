package model


type GrpcRule struct {
	ID        uint `json:"id" gorm:"primary_key"`
	ServiceInfoID int64 `json:"service_id" gorm:"comment:服务id;not null;"`
	Port      int   `json:"port" gorm:"comment:端口;not null;"`
	HeaderTransfor string `json:"header_transfor" gorm:"comment:header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue;type:varchar(5000);not null;"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_grpc_rule"
}
