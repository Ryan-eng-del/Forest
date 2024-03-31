package oauthController

import "github.com/gin-gonic/gin"


type OauthController struct {
}


func Register(i gin.IRoutes) {
	oauthController := &OauthController{}
	i.POST("/tokens", oauthController.AuthToken)
}

// AuthTokenList godoc
// @Summary 租户获取 token
// @Description 租户获取 token
// @Tags Oauth














func (o *OauthController) AuthToken(c *gin.Context) {

}

// @ID /oauth/token
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/token [post]
