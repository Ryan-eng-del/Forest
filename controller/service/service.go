package serviceController

import (
	"errors"
	serviceDto "go-gateway/dto/service"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {

}

func Register(group gin.IRoutes) {
	service := &ServiceController{}
	group.GET("", service.ServiceList)
	group.DELETE("/:serviceId", service.ServiceDelete)
	group.GET("/:serviceId", service.ServiceDetail)
}

func (s *ServiceController) ServiceCreate(c *gin.Context) {
	
}


// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags Service
// @ID /service
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} public.Response{data=serviceDto.ServiceListOutput} "success"
// @Router /service [get]
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
	 lists, count, err := serviceModel.PageList(c, tx, &serviceListInput)

	 if err != nil {
		public.ResponseError(c, 2002, err)
		return
	 } 

	 outList := []serviceDto.ServiceListItemOutput{}
	 for _, listItem := range lists {
			listItem.ServiceDetail(c, tx)
			serverAddr := "unknown"
			serviceListItem := serviceDto.ServiceListItemOutput{
				ID: int64(listItem.ID),
				ServiceAddr: serverAddr,
				ServiceName: listItem.ServiceName,
				ServiceDesc: listItem.ServiceDesc,
				LoadType: public.LoadType(listItem.LoadType),
				Qps: 0,
				Qpd: 1,
				TotalNode: 2,
				CreateAt: public.LocalTime(listItem.CreateAt),
				UpdateAt: public.LocalTime(listItem.UpdateAt),
			}

			outList = append(outList, serviceListItem)
	 }

	 out := &serviceDto.ServiceListOutput{
		Total: int(count),
		List: outList,
	 }
	 public.ResponseSuccess(c, out)
	 
}


func (s *ServiceController) ServiceDelete(c *gin.Context) {

}


// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务列表
// @Tags Service
// @ID /service/{service_id}
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param service_id path string true "服务id"
// @Success 200 {object} public.Response{data=model.ServiceDetail} "success"
// @Router /service/{service_id} [get]
func (s *ServiceController) ServiceDetail(c *gin.Context) {
	serviceId, err := strconv.ParseInt(c.Param("serviceId"), 10, 64)
	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
	}

	service := &model.Service{}
	tx, err :=libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2002), err)
		return
	}

	if err := tx.First(service, serviceId).Error; err != nil {
		public.ResponseError(c, public.ResponseCode(2003), errors.New("服务不存在"))
		return
	}

	 serviceDetail, err := service.ServiceDetail(c, tx); 

	 if err != nil {
		public.ResponseError(c, public.ResponseCode(2004), err)
		return
	}
	
	public.ResponseSuccess(c, serviceDetail)
}