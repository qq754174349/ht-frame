package web

import (
	"github.com/gin-gonic/gin"
	"ht-crm/autoconfigure"
	"ht-crm/logger"
	"ht-crm/src/ht/web/middlewares"
	"ht-crm/src/ht/web/routes"
)

func Start() {
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	r := gin.Default()
	r.Use(middlewares.GenTraceId(), middlewares.ReqInfoLogger())

	// 路由注册
	routes.RegisterRoutes(r)
	port := autoconfigure.GetAppCig().Web.Port

	err := r.Run(":" + port)
	if err != nil {
		logger.Error(err)
		return
	}
}
