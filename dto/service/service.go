package serviceDto

import (
	lib "go-gateway/lib/func"

	"github.com/gin-gonic/gin"
)
 
type ServiceListInput struct {
	Info string `json:"info" form:"info" comment:"关键词" en_comment:"keyword" validate:"" example:""`
	PageNo int `json:"page_no" form:"page_no" comment:"页数" en_comment:"page number" validate:"required"`
	PageSize int `json:"page_size" form:"page_size" comment:"页容量" en_comment:"page size" validate:"required"`
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return lib.ValidateParams(c, param)
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	Qps         int64  `json:"qps" form:"qps"`                   //qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   //qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
}

type ServiceListOutput struct {
	Total int `json:"total" form:"total" comment:"总数" example:"" validate:""`
	List []ServiceListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`
}