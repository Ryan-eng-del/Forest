package dashboardController

import (
	"errors"
	appDto "go-gateway/dto/app"
	dashboardDto "go-gateway/dto/dashboard"
	serviceDto "go-gateway/dto/service"
	"go-gateway/handler"
	libConf "go-gateway/lib/conf"
	libConst "go-gateway/lib/const"
	lib "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func DashboardRegister(group gin.IRoutes) {
	service := &DashboardController{}
	group.GET("/panel_group_data", service.PanelGroupData)
	group.GET("/flow_stat", service.FlowStat)
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags Dashboard
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} public.Response{data=dashboardDto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (service *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		public.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &model.Service{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &serviceDto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}
	app := &model.App{}

	_, appNum, err := app.AppList(c, tx, &appDto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}

	counter, err := handler.ServerCountHandler.GetCounter(libConst.FlowTotal)
	if err != nil {
		public.ResponseError(c, 2003, err)
		return
	}
	out := &dashboardDto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: counter.TotalCount,
		CurrentQPS:      counter.QPS,
	}
	public.ResponseSuccess(c, out)
}


// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags Dashboard
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} public.Response{data=dashboardDto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (service *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		public.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &model.Service{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := libConst.LoadTypeMap[item.LoadType]
		if !ok {
			public.ResponseError(c, 2003, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}

	out := &dashboardDto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	public.ResponseSuccess(c, out)
}


// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags Dashboard
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (service *DashboardController) FlowStat(c *gin.Context) {
	counter, err := handler.ServerCountHandler.GetCounter(libConst.FlowTotal)
	if err != nil {
		public.ResponseError(c, 2001, err)
		return
	}

	todayList := []int64{}
	currentTime := time.Now()

	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, libConf.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, libConf.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	public.ResponseSuccess(c, &serviceDto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}