package controller

import (
	"encoding/json"
	"github.com/go_gateway/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
	"time"
)

type AdminLoginController struct {

}

// 登录接口
func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	// 登录接口
	group.POST("/login", adminLogin.AdminLogin)
	// 登出接口
	group.GET("/logout", adminLogin.AdminLogout)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	// 登录方法 参数是上下文的一个指针
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	// １．params.UserName 取得管理员信息 admininfo
	// ２．admininfo.salt + params.Password sha256 => saltPassword
	// ３．saltPassword == admininfo.password
	// 对用户名以及密码进行校验
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 设置登录session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	session := sessions.Default(c)
	session.Set(public.AdminSessionInfoKey, string(sessBts))
	// 保存至Redis中
	session.Save()

	// 输出
	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}



// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/logout [get]
func (adminlogin *AdminLoginController) AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	// 对该session进行删除处理，以实现登出功能
	session.Delete(public.AdminSessionInfoKey)
	// 保存至Redis中
	session.Save()

	middleware.ResponseSuccess(c, "")
}

