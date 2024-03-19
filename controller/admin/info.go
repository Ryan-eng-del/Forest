package adminController

import (
	"errors"
	adminDto "go-gateway/dto/admin"
	lib "go-gateway/lib/func"
	mysqlLib "go-gateway/lib/mysql"

	"go-gateway/model"
	public "go-gateway/public"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {}



func RegisterAuth(group gin.IRoutes) {
	adminInfo := &AdminController{}
	group.POST("/change_pwd", adminInfo.AdminChangePwd)
	group.GET("/info", adminInfo.AdminInfo)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags Admin
// @ID /admin/info
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {object} public.Response{data=adminDto.AdminInfoOutput} "success"
// @Router /admin/info [get]
func (a *AdminController) AdminInfo(c *gin.Context) {
	tx, _ := mysqlLib.GetGormPool("default")
	admin, err := a.GetAdmin(c, tx)
	if err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	output :=  adminDto.AdminInfoOutput{
		ID: admin.ID,
		Name: admin.UserName,
		LoginTime:     time.Now(),
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	public.ResponseSuccess(c, output)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags Admin
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @security ApiKeyAuth
// @Param body body adminDto.ChangePwdInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) AdminChangePwd(c *gin.Context) {
	params := adminDto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2003, err)
	}

	tx, _ := mysqlLib.GetGormPool("default")
	admin, err := a.GetAdmin(c, tx)

	if err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	admin.Password = lib.GenSaltPassword(admin.Salt, params.Password)
	if err := admin.Save(c, tx); err != nil {
		public.ResponseError(c, 2001, err)
		return
	}
	
	public.ResponseSuccess(c, "")
}


func (a *AdminController) GetAdmin(c *gin.Context, tx *gorm.DB) (*model.Admin, error){
	userId, _:= c.Get("UserID")
	adminId, ok := userId.(uint); 
	if !ok {
		return nil, errors.New("not a valid user")
	}

	admin := &model.Admin{}
	
	admin, err := admin.FindById(c, tx, adminId)
	return admin, err
}