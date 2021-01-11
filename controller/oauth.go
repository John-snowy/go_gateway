package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/golang_common/lib"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
	"strings"
	"time"
)

type OAuthController struct {

}

// 登录接口
func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	// token接口
	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	// 登录方法 参数是上下文的一个指针
	params := &dto.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	fmt.Println(string(appSecret))

	// 取出app_id secret
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}
	//appID := parts[0]
	// 生成app_list
	appList := dao.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		// 匹配app_id
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			// 基于jwt生成token
			claims := jwt.StandardClaims{
				Issuer:appInfo.AppID,
				ExpiresAt:time.Now().Add(public.JwtExpires*time.Second).In(lib.TimeLocation).Unix(),
			}
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2004, err)
				return
			}
			// 生成output
			output := &dto.TokensOutput{
				AccessType: token,
				ExpiresIn:  public.JwtExpires,
				TokenType:  "Bearer",
				Scope:      "read_write",
			}
			middleware.ResponseSuccess(c, output)
			return
		}
	}
	middleware.ResponseError(c, 2005, errors.New("未匹配正确APP信息"))
}


