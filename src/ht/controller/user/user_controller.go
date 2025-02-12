package user

import (
	"github.com/gin-gonic/gin"
	"ht-crm/src/ht/common/result"
	"ht-crm/src/ht/dto/req"
	"ht-crm/src/ht/service/user"
)

// WechatUserLogin 微信小程序登录
func WechatUserLogin(ctx *gin.Context) {
	code := ctx.Query("code")
	traceId := ctx.GetString("traceId")

	jwt, err := user.WechatUserLogin(ctx, code)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
	} else {
		ctx.Writer.WriteString(result.NewSuccessResult(traceId, jwt).ToString())
	}
}

// WechatUserReg 微信小程序注册
func WechatUserReg(ctx *gin.Context) {
	req := req.WechatUserRegReq{}
	ctx.BindJSON(&req)
	user.WechatUserReg(ctx, req)
}
