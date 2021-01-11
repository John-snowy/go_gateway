package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/public"
	"time"
)

// 登录session结构体
type AdminSessionInfo struct {
	ID			 int 	`json:"id"`
	UserName	 string `json:"username"`
	LoginTime	 time.Time 	`json:"login_time"`

}

// 登录输入参数
type AdminLoginInput struct {
	UserName string `json:"username" form:"username"comment:"姓名" example:"admin" validate:"required,valid_username"`// 管理员用户名
	Password string `json:"password" form:"password" comment:"姓名" example:"123456" validate:"required"`// 密码
}

func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

// 登录后输出结构
type AdminLoginOutput struct {
	Token string `json:"token" form:"token"comment:"token" example:"token" validate:""`// token
}
