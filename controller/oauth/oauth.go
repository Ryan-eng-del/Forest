package oauthController

import (
	"encoding/base64"
	"errors"
	oauthDto "go-gateway/dto/oauth"
	"go-gateway/handler"
	lib "go-gateway/lib/conf"
	"go-gateway/lib/jwt"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)


type OauthController struct {
}


func Register(i gin.IRoutes) {
	oauthController := &OauthController{}
	i.POST("/tokens", oauthController.AuthToken)
}


// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body oauthDto.TokensInput true "body"
// @Success 200 {object} public.Response{data=oauthDto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (o *OauthController) AuthToken(c *gin.Context) {
	params := oauthDto.TokensInput{}

	if err := params.BindValidParam(c); err != nil {
		public.ResponseError(c, 2000, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
	  public.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])


	if err != nil {
		public.ResponseError(c, 2002, err)
		return
	}

	parts := strings.Split(string(appSecret), ":")

	if len(parts) != 2 {
		public.ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	appList := handler.AppManagerHandler.GetAppList()


	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			jwt := jwt.NewJWT()
			token, err := jwt.GenerateTokenWithAppID(parts[0])

			if err != nil {
				public.ResponseError(c, 2004, err)
				return
			}

			output := &oauthDto.TokensOutput{
				ExpiresIn: int(lib.TokenExpirePeriod),
				TokenType:"Bearer",
				AccessToken:token,
				Scope:"read_write",
			}
			public.ResponseSuccess(c, output)
		}
	}
}


