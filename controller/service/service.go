package serviceController

import (
	serviceDto "go-gateway/dto/service"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {

}

func Register(group gin.IRouter) {
	service := &ServiceController{}
	group.GET("/service", service.ServiceList)
	group.DELETE("/service/:serviceId", service.ServiceDelete)
	group.GET("/service/:serviceId", service.ServiceDetail)
}

func (s *ServiceController) ServiceCreate(c *gin.Context) {
	
}

func (s *ServiceController) ServiceList(c *gin.Context) {
	 serviceListInput := serviceDto.ServiceListInput{}
	 
	 if err := serviceListInput.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
	 }


	 tx, err := libMysql.GetGormPool("default")
	 if err != nil {
		public.ResponseError(c, 2001, err)
		return
	 } 

	 serviceModel := model.Service{}
	 _, _, err = serviceModel.PageList(c, tx, &serviceListInput)

	 if err != nil {
		public.ResponseError(c, 2002, err)
		return
	 } 

	//  outList := []serviceDto.ServiceListItemOutput{}

	//  for _, listItem := range list {
	// 	li

	//  }
	 
}


func (s *ServiceController) ServiceDelete(c *gin.Context) {

}

func (s *ServiceController) ServiceDetail(c *gin.Context) {

}