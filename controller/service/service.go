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
	"strings"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {}

func Register(group gin.IRoutes) {
	service := &ServiceController{}
	// todo 补充 Redis 服务统计功能
	group.GET("", service.ServiceList)
	group.DELETE("/:serviceId", service.ServiceDelete)
	group.GET("/:serviceId", service.ServiceDetail)
	group.POST("/http", service.ServiceCreateHttp)
	group.PATCH("/http/:serviceId", service.ServiceHttpUpdate)
	group.POST("/tcp", service.ServiceCreateTcp)
	group.PATCH("/tcp/:serviceId", service.ServiceUpdateTcp)
	group.POST("/grpc", service.ServiceCreateGrpc)
	group.PATCH("/grpc/:serviceId", service.ServiceUpdateGrpc)
}


// ServiceCreate godoc
// @Summary 创建 tcp 服务
// @Description 创建 tcp 服务
// @Tags Service
// @ID /service/tcp
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param info body serviceDto.ServiceAddTcpInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/tcp [post]
func (s *ServiceController) ServiceCreateTcp(c *gin.Context) {
	params := serviceDto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2004, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx, err := libMysql.GetGormPool("default")
	
	if err != nil {
		public.ResponseError(c, 2005, err)
	}
	serviceInfo := &model.Service{ServiceName: params.ServiceName}
	if _, err := serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		public.ResponseError(c, public.ResponseCode(2006), errors.New("服务已经存在"))
		return
	}

	tcpRuleSearch := &model.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.FindMust(c, tx, tcpRuleSearch); err == nil {
		public.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}


	grpcRuleSearch := &model.GrpcRule{
		Port: params.Port,
	}

	if _, err := grpcRuleSearch.FindMust(c, tx, grpcRuleSearch); err == nil {
		public.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}

	tx = tx.Begin()

	info := &model.Service{
		LoadType: libConst.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &model.LoadBalance{
		ServiceInfoID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}

	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2007, err)
		return
	}

	tcpRule := &model.TcpRule{
		ServiceInfoID: info.ID,
		Port:      params.Port,
	}
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2008, err)
		return
	}

	accessControl := &model.AccessControl{
		ServiceInfoID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}

	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	public.ResponseSuccess(c, "创建成功")
}

// ServiceHttpUpdate godoc
// @Summary 更新 tcp 服务
// @Description 更新 tcp 服务
// @Tags Service
// @ID /service/tcp/{service_id}
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param service_id path string true "服务id"
// @Param info body serviceDto.ServiceAddTcpInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/tcp/{service_id} [patch]
func (s *ServiceController) ServiceUpdateTcp(c *gin.Context) {
	serviceIdStr := c.Param("serviceId")
	serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	params := serviceDto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}


	serviceInfo := model.Service{}
	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2002), err)
		return
	}

	tx = tx.Begin()
	service, err := serviceInfo.Find(c, tx, &model.Service{
		AbstractModel: model.AbstractModel{ID: uint(serviceId)},
	})

	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	serviceDetail, err := service.ServiceDetail(c, tx)



	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc

	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2005, err)
		return
	}

	loadBalance := &model.LoadBalance{}
	if serviceDetail.LoadBalance != nil {
		loadBalance = serviceDetail.LoadBalance
	}
	loadBalance.ServiceInfoID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList

	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2004, err)
		return
	}

	tcpRule := &model.TcpRule{}
	if serviceDetail.TCPRule != nil {
		tcpRule = serviceDetail.TCPRule
	}
	tcpRule.ServiceInfoID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2005, err)
		return
	}

	accessControl := &model.AccessControl{}
	if serviceDetail.AccessControl != nil {
		accessControl = serviceDetail.AccessControl
	}
	accessControl.ServiceInfoID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit

	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2006, err)
		return
	}
	tx.Commit()
	public.ResponseSuccess(c, "更新成功")
}

// ServiceCreate godoc
// @Summary 创建 grpc 服务
// @Description 创建 grpc 服务
// @Tags Service
// @ID /service/grpc
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param info body serviceDto.ServiceAddGrpcInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/grpc [post]
func (s *ServiceController) ServiceCreateGrpc(c *gin.Context) {
	params := &serviceDto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2004, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx, err := libMysql.GetGormPool("default")
	
	if err != nil {
		public.ResponseError(c, 2005, err)
	}

	tcpRuleSearch := &model.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.FindMust(c, tx, tcpRuleSearch); err == nil {
		public.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}

	grpcRuleSearch := &model.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.FindMust(c, tx, grpcRuleSearch); err == nil {
		public.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

		//ip与权重数量一致
		if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
			public.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
			return
		}

	tx = tx.Begin()


	info := &model.Service{
		LoadType: libConst.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &model.LoadBalance{
		ServiceInfoID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}

	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2007, err)
		return
	}

	grpcRule := &model.GrpcRule{
		ServiceInfoID: info.ID,
		Port:      params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}

	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2008, err)
		return
	}

	accessControl := &model.AccessControl{
		ServiceInfoID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}

	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2009, err)
		return
	}

	tx.Commit()
	public.ResponseSuccess(c, "创建成功")
}

// ServiceGRPCUpdate godoc
// @Summary 更新 grpc 服务
// @Description 更新 grpc 服务
// @Tags Service
// @ID /service/grpc/{service_id}
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param service_id path string true "服务id"
// @Param info body serviceDto.ServiceAddGrpcInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/grpc/{service_id} [patch]
func (s *ServiceController) ServiceUpdateGrpc(c *gin.Context) {
	serviceIdStr := c.Param("serviceId")
	serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	params := serviceDto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}


	serviceInfo := model.Service{}
	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2002), err)
		return
	}

	tx = tx.Begin()
	service, err := serviceInfo.Find(c, tx, &model.Service{
		AbstractModel: model.AbstractModel{ID: uint(serviceId)},
	})

	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	serviceDetail, err := service.ServiceDetail(c, tx)



	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc

	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2005, err)
		return
	}

	loadBalance := &model.LoadBalance{}
	if serviceDetail.LoadBalance != nil {
		loadBalance = serviceDetail.LoadBalance
	}
	loadBalance.ServiceInfoID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList

	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2004, err)
		return
	}

	grpcRule := &model.GrpcRule{}
	if serviceDetail.GRPCRule != nil {
		grpcRule = serviceDetail.GRPCRule
	}

	grpcRule.ServiceInfoID = info.ID
	grpcRule.HeaderTransfor = params.HeaderTransfor

	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2005, err)
		return
	}

	accessControl := &model.AccessControl{}
	if serviceDetail.AccessControl != nil {
		accessControl = serviceDetail.AccessControl
	}
	accessControl.ServiceInfoID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit

	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2006, err)
		return
	}
	tx.Commit()
	public.ResponseSuccess(c, "更新成功")
}

// ServiceHttpUpdate godoc
// @Summary 更新 http 服务
// @Description 更新 http 服务
// @Tags Service
// @ID /service/http/{service_id}
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param service_id path string true "服务id"
// @Param info body serviceDto.ServiceAddHttpInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/http/{service_id} [patch]
func (s *ServiceController) ServiceHttpUpdate(c *gin.Context) {
	serviceIdStr := c.Param("serviceId")
	serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	params := serviceDto.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	serviceInfo := model.Service{}
	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(c, public.ResponseCode(2002), err)
		return
	}

	tx = tx.Begin()
	service, err := serviceInfo.Find(c, tx, &model.Service{
		AbstractModel: model.AbstractModel{ID: uint(serviceId)},
	})

	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	serviceDetail, err := service.ServiceDetail(c, tx)

	if err != nil {
		tx.Rollback()
		public.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc

	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2005, err)
		return
	}


	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor

	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2006, err)
		return
	}


	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2007, err)
		return
	}



	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = params.RoundType
	loadbalance.IpList = params.IpList
	loadbalance.WeightList = params.WeightList
	loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadbalance.Save(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2008, err)
		return
	}

	tx.Commit()
	public.ResponseSuccess(c, "更新成功")
}
 
// ServiceCreate godoc
// @Summary 创建 http 服务
// @Description 创建 http 服务
// @Tags Service
// @ID /service/http
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param info body serviceDto.ServiceAddHttpInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /service/http [post]
func (s *ServiceController) ServiceCreateHttp(c *gin.Context) {
	params := &serviceDto.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, public.ResponseCode(2001), err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		public.ResponseError(c, 2004, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx, err := libMysql.GetGormPool("default")
	
	if err != nil {
		public.ResponseError(c, 2005, err)
	}

	tx = tx.Begin()

	serviceInfo := &model.Service{ServiceName: params.ServiceName}
	if _, err := serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		public.ResponseError(c, public.ResponseCode(2006), errors.New("服务已经存在"))
		return
	}

	httpRuleInfo := &model.HttpRule{
		RuleType: params.RuleType,
		Rule: params.Rule,
	}


	if _, err := httpRuleInfo.FindMust(c, tx, httpRuleInfo); err == nil {
		tx.Rollback()
		public.ResponseError(c, public.ResponseCode(2007), errors.New("服务接入前缀或域名已存在"))
		return
	}


	serviceModel := &model.Service{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := serviceModel.Create(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2008, err)
		return
	}

	httpRule := &model.HttpRule{
		ServiceInfoID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}

	if err := httpRule.Create(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2009, err)
		return
	}

	accessControl := &model.AccessControl{
		ServiceInfoID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}

	if err := accessControl.Create(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2010, err)
		return
	}

	loadbalance := &model.LoadBalance{
		ServiceInfoID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}

	if err := loadbalance.Create(c, tx); err != nil {
		tx.Rollback()
		public.ResponseError(c, 2011, err)
		return
	}

	tx.Commit()
	public.ResponseSuccess(c, "创建成功")
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
		serviceDetail, err := listItem.ServiceDetail(c, tx)

		if err != nil {
			public.ResponseError(c, public.ResponseCode(2003), err)
			return
		}

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
	serviceInstance, err := serviceModel.Find(c, tx, &model.Service{
		AbstractModel: model.AbstractModel{
			ID: uint(serviceId),
		},
		},
	)

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