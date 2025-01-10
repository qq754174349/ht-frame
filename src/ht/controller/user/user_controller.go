package user

import (
	"github.com/gin-gonic/gin"
	"ht-crm/src/ht/dto/req"
	"ht-crm/src/ht/service"
)

// WechatUserLogin 微信小程序登录
func WechatUserLogin(ctx *gin.Context) {
	loginReq := req.WechatProgramLoginReq{}
	ctx.BindJSON(&loginReq)

	service.Code2Session(loginReq.Code)

}

// WechatUserReg 微信小程序注册
func WechatUserReg(ctx *gin.Context) {

}

// Save 说
func Save(ctx *gin.Context) {
	//ctx.Writer.WriteString(service.GetAccToken())
}
