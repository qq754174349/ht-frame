package user

import (
	"github.com/gin-gonic/gin"
	"ht-crm/src/ht/service"
)

// WechatProgramLogin 微信小程序登录
func WechatProgramLogin(ctx *gin.Context) {
	code := ctx.Query("code")
	service.Code2Session(code)

}

// Save 说
func Save(ctx *gin.Context) {
	//ctx.Writer.WriteString(service.GetAccToken())
}
