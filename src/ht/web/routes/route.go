package routes

import (
	"github.com/gin-gonic/gin"
	"ht-crm/src/ht/controller/user"
)

func RegisterRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/wechat/Program/login", user.WechatUserLogin)
		}
		testGroup := apiGroup.Group("/test")
		{
			testGroup.POST("/say", user.Save)
		}
	}
}
