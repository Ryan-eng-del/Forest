package serviceController

import (
	"errors"
	"fmt"
	serviceDto "go-gateway/dto/service"
	libConst "go-gateway/lib/const"
	libLog "go-gateway/lib/log"
	libMysql "go-gateway/lib/mysql"
	libViper "go-gateway/lib/viper"
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
	group.POST("/http", service.ServiceCreateHttp)
}

// ServiceList godoc
// @Summary 创建 http 服务
// @Description 创建 http 服务
// @Tags Service
// @ID /service
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param body body serviceDto.ServiceAddHttp "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service [get]
func (s *ServiceController) ServiceCreateHttp(c *gin.Context) {
	
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
		serviceDetail, _ := listItem.ServiceDetail(c, tx)
		log := libLog.GetLogger()
		log.Info(serviceDetail.Info.ServiceName)
		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := libViper.ViperInstance.GetStringConf("base.cluster.cluster_ip")
		clusterPort :=  libViper.ViperInstance.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort :=  libViper.ViperInstance.GetStringConf("base.cluster.cluster_ssl_port")

		if serviceDetail.Info.LoadType == libConst.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == libConst.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == libConst.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == libConst.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == libConst.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == libConst.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == libConst.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == libConst.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}

		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		serviceListItem := serviceDto.ServiceListItemOutput{
			ID: int64(listItem.ID),
			ServiceAddr: serviceAddr,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType: public.LoadType(listItem.LoadType),
			// todo 完善 Qps Qpd
			Qps: 0,
			Qpd: 1,
			TotalNode: len(ipList),
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

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags Service
// @ID /service/{service_id}
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param service_id path string true "服务ID"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/{service_id} [delete]
func (s *ServiceController) ServiceDelete(c *gin.Context) {
	serviceIdStr := c.Param("serviceId")
	serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)
	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), errors.New("not a valid service id"))
		return
	}

	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	serviceModel := &model.Service{}
	serviceInstance, err := serviceModel.FindById(c, tx, int(serviceId))
	if err != nil {
		public.ResponseError(c, public.ResponseCode(2002), err)
		return
	}

	serviceInstance.IsDelete = 1
	if err := serviceInstance.Save(c, tx); err != nil {
		public.ResponseError(c, public.ResponseCode(2003), err)
		return
	}

	public.ResponseSuccess(c, "删除成功")
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

	if err := tx.Scopes(libMysql.LogicalObjects()).First(service, serviceId).Error; err != nil {
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