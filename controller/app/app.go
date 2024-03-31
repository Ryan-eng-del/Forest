package appController

import (
	"errors"
	appDto "go-gateway/dto/app"
	"strconv"

	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)


type AppController struct {}

func Register (i gin.IRoutes) {
	appController := &AppController{}
	i.GET("", appController.AppList)
	i.GET("/:appId", appController.AppDetail)
	i.GET("/statics", appController.AppStatistics)
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
		outputList = append(outputList, appDto.APPListItemOutput{
			ID:       int64(item.ID),
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIps,
			Qpd:      int64(item.Qpd),
			Qps:      int64(item.Qps),
			RealQpd:  1,
			RealQps:  1,
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


func (i *AppController) AppAdd(ctx *gin.Context) {

}

func (i *AppController) AppUpdate(ctx *gin.Context) {

}