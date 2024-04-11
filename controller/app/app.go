package appController

import (
	"errors"
	appDto "go-gateway/dto/app"
	"go-gateway/handler"
	lib "go-gateway/lib/const"
	libFunc "go-gateway/lib/func"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"
	"strconv"

	"github.com/gin-gonic/gin"
)


type AppController struct {}

func Register (i gin.IRoutes) {
	appController := &AppController{}
	i.GET("", appController.AppList)
	i.GET("/:appId", appController.AppDetail)
	i.GET("/statistics", appController.AppStatistics)
	i.DELETE("/:appId", appController.AppDelete)
	i.PATCH("/:appId", appController.AppUpdate)
	i.POST("", appController.AppAdd)
}


// APPList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags App
// @ID /app
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} public.Response{data=appDto.APPListOutput} "success"
// @Router /app [get]
func (i *AppController) AppList(ctx *gin.Context) {
	params := appDto.APPListInput{}
	if err := params.GetValidParams(ctx);  err  != nil {
		public.ResponseError(ctx, 2001, err)
		return
	}


	appInfo := model.App{}
	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(ctx, 2001, err)
		return
	}

	appList, count, err := appInfo.AppList(ctx, tx, &params)

	if err != nil {
		public.ResponseError(ctx, 2002, err)
		return
	}

	outputList := []appDto.APPListItemOutput{}

	for _, item := range appList {

		appCounter, err := handler.ServerCountHandler.GetCounter(lib.FlowAppPrefix  + item.AppID)

		if err != nil {
			public.ResponseError(ctx, 2003, err)
			ctx.Abort()
			return
		}

		outputList = append(outputList, appDto.APPListItemOutput{
			ID:       int64(item.ID),
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIps,
			Qpd:      int64(item.Qpd),
			Qps:      int64(item.Qps),
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
			CreatedAt: item.CreateAt,
			UpdatedAt: item.UpdateAt,
		})
	}

	output := appDto.APPListOutput{}
	output.List = outputList
	output.Total = count
	public.ResponseSuccess(ctx, output)
}


// APPDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags App
// @ID /app/{app_id}/get
// @Accept  json
// @Produce  json
// @Param app_id path string true "服务id"
// @Success 200 {object} public.Response{data=model.App} "success"
// @Router /app/{app_id} [get]
func (i *AppController) AppDetail(ctx *gin.Context) {
	appIdStr := ctx.Param("appId")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2001), errors.New("not a valid app id"))
		return
	}

	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2002), err)
		return
	}

	search := model.App{
		AbstractModel: model.AbstractModel{
			ID: uint(appId),
		},
	}

	app, err := search.Find(ctx, tx, &search)

	if err != nil {
		public.ResponseError(ctx, 2002, err)
		return
	}

	public.ResponseSuccess(ctx, app)
}

func (i *AppController) AppStatistics(ctx *gin.Context) {

}


// AppDelete godoc
// @Summary 删除租户
// @Description 删除租户
// @Tags App
// @ID /app/{app_id}/delete
// @Accept  json
// @Produce  json
// @Param app_id path string true "服务id"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /app/{app_id} [delete]
func (i *AppController) AppDelete(ctx *gin.Context) {
	appIdStr := ctx.Param("appId")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2001), errors.New("not a valid app id"))
		return
	}

	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2002), err)
		return
	}

	search := model.App{
		AbstractModel: model.AbstractModel{
			ID: uint(appId),
		},
	}

	app, err := search.Find(ctx, tx, &search)

	if err != nil {
		public.ResponseError(ctx, 2002, err)
		return
	}
	app.IsDelete = 1
	app.Save(ctx, tx)
	public.ResponseSuccess(ctx, "删除成功")

}


// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags App
// @ID /app/ post
// @Accept  json
// @Produce  json
// @Param body body appDto.APPAddHttpInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /app [post]
func (i *AppController) AppAdd(ctx *gin.Context) {
	params := appDto.APPAddHttpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		public.ResponseError(ctx, 2001, err)
	}


	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2002), err)
		return
	}

	search := model.App{
		AppID: params.AppID,
	}

	_, err = search.Find(ctx, tx, &search)

	if err == nil {
		public.ResponseError(ctx, public.ResponseCode(2003), errors.New("租户ID被占用，请重新输入"))
		return
	}

	if params.Secret == "" {
		params.Secret = libFunc.MD5(params.AppID)
	}

	info := &model.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIps: params.WhiteIPS,
		Qps:      uint(params.Qps),
		Qpd:      uint(params.Qpd),
	}

	if err := info.Save(ctx, tx); err != nil {
		public.ResponseError(ctx, 2003, err)
		return
	}

	public.ResponseSuccess(ctx, "添加成功")
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags App
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body appDto.APPAddHttpInput true "body"
// @Param app_id path string true "app id"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /app/{app_id} [patch]
func (i *AppController) AppUpdate(ctx *gin.Context) {
	appIdStr := ctx.Param("appId")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2001), errors.New("not a valid app id"))
		return
	}

	params := appDto.APPAddHttpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		public.ResponseError(ctx, 2001, err)
	}


	tx, err := libMysql.GetGormPool("default")

	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2002), err)
		return
	}

	search := model.App{
		AbstractModel: model.AbstractModel{
			ID: uint(appId),
		},
	}

	info, err := search.Find(ctx, tx, &search)
	if err != nil {
		public.ResponseError(ctx, public.ResponseCode(2003), err)
		return
	}

	if params.Secret == "" {
		params.Secret = libFunc.MD5(info.AppID)
	}

	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIps = params.WhiteIPS
	info.Qps = uint(params.Qps)
	info.Qpd = uint(params.Qpd)

	if err := info.Save(ctx, tx); err != nil {
		public.ResponseError(ctx, public.ResponseCode(2003), err)
		return
	}
	
	public.ResponseSuccess(ctx, "更新成功")
}