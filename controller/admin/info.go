package adminController

import (
	"encoding/json"
	"fmt"
	adminDto "go-gateway/dto/admin"
	constLib "go-gateway/lib/const"
	lib "go-gateway/lib/func"
	mysqlLib "go-gateway/lib/mysql"
	"go-gateway/model"
	public "go-gateway/public"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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
// @Success 200 {object} public.Response{data=adminDto.AdminInfoOutput} "success"
// @Router /admin/info [get]
func (a *AdminController) AdminInfo(c *gin.Context) {
	// anyAdmin, exist := c.Get("admin")
	// if !exist {
	// 	public.ResponseError(c, 2000, errors.New("admin not found"))
	// 	return
	// }

	// admin, ok := anyAdmin.(*model.Admin); 
	// if !ok {
	// 	public.ResponseError(c, 2000, errors.New("not a model admin"))
	// 	return
	// }

	sess := sessions.Default(c)
	sessInfo := sess.Get(constLib.AdminSessionInfoKey)
	adminSessionInfo := &adminDto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}
	
	output :=  adminDto.AdminInfoOutput{
		ID: uint(adminSessionInfo.ID),
		Name: adminSessionInfo.UserName,
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
// @Param body body adminDto.ChangePwdInput true "body"
// @Success 200 {object} public.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) AdminChangePwd(c *gin.Context) {
	params := adminDto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2003, err)
	}


	sess := sessions.Default(c)
	sessInfo := sess.Get(constLib.AdminSessionInfoKey)
	adminSessionInfo := &adminDto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}
	
	tx, _ := mysqlLib.GetGormPool("default")

	adminInfo := &model.Admin{}
	adminInfo, err := adminInfo.FindById(c, tx, uint(adminSessionInfo.ID))
	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}

	adminInfo.Password = lib.GenSaltPassword(adminInfo.Salt, params.Password)
	if err := adminInfo.Save(c, tx); err != nil {
		public.ResponseError(c, 2001, err)
		return
	}
	public.ResponseSuccess(c, "")
}

