package adminController

import (
	"encoding/json"
	adminDto "go-gateway/dto/admin"
	lib "go-gateway/lib/const"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	public "go-gateway/public"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)


type AdminLoginController struct {}

func Register(group gin.IRoutes) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.Login)
	group.GET("/logout", adminLogin.LoginOut)
}


// GO-Gateway godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags Admin
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body adminDto.AdminLoginInput true "body"
// @Success 200 {object} public.Response{data=commonDto.TokensOutput} "success"
// @Router /admin_login/login [post]	
func (adminLogin *AdminLoginController) Login (c *gin.Context) {
	params := adminDto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	tx, err := libMysql.GetGormPool("default")
	if err != nil {
		public.ResponseError(c, 2001, err)
		return
	}

	admin := &model.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}
	//设置session
	sessInfo := &adminDto.AdminSessionInfo{
		ID:        int(admin.ID),
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		public.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(lib.AdminSessionInfoKey, string(sessBts))
	sess.Save()
	out := &adminDto.AdminLoginOutput{Token: admin.UserName}
	public.ResponseSuccess(c, out)
	// token, err := _jwt.NewJWT().GenerateTokenWithUserID(admin.ID)
	// if err != nil {
	// 	pubic.ResponseError(c, 2004, err)
	// 	return
	// }

	// output := &commonDto.TokensOutput{
	// 	ExpiresIn: 3600 * 7,
	// 	TokenType:"Bearer",
	// 	AccessToken: token,
	// 	Scope: "all",
	// }
	// pubic.ResponseSuccess(c, output)
}


// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags Admin
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} public.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminLogin *AdminLoginController) LoginOut (c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(lib.AdminSessionInfoKey)
	sess.Save()
	public.ResponseSuccess(c, "")
}