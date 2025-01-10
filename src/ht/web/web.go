package web

import (
	"github.com/gin-gonic/gin"
	"ht-crm/src/ht/config"
	"ht-crm/src/ht/config/log"
	"ht-crm/src/ht/web/middlewares"
	"ht-crm/src/ht/web/routes"
)

func Start() {
	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()
	r := gin.Default()
	r.Use(middlewares.GenTraceId(), middlewares.ReqInfoLogger())

	// 路由注册
	routes.RegisterRoutes(r)
	port := config.GetEnvCfg().Port

	err := r.Run(":" + port)
	if err != nil {
		log.Error(err)
		return
	}
}
