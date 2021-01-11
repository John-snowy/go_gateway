package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go_gateway/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
)

type AdminController struct {
}

// 登录接口
func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	// admin修改密码接口
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminlogin *AdminController) AdminInfo(c *gin.Context) {
	// １．读取sessionKey对应json转换为结构体
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	// 将获取到的sessionInfo进行数据类型转换，转换成string类型
	//sessionInfoStr := sessionInfo.(string)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// ２．取出数据然后封装输出结构体

	// 输出
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "Iam a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}


// ChangePwd godoc
// @Summary 管理员修改密码
// @Description 管理员修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminlogin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}

	// １．session读取用户信息到结构体　sessionInfo

	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	// 将获取到的sessionInfo进行数据类型转换，转换成string类型
	//sessionInfoStr := sessionInfo.(string)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// ２．sessionInfo.ID　读取数据库信息 adminInfo

	// 从数据库中读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err =adminInfo.Find(c, tx,
		(&dao.Admin{UserName:adminSessionInfo.UserName}))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// ３．params.password+adminInfo.salt sha256 saltPassword  生成新密码
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)

	// ４．saltPassword => adminInfo.password 执行数据保存
	adminInfo.Password = saltPassword
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}