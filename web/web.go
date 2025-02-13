package web

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht-frame/web/middlewares"
)

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.Default(opts...)
	engine.Use(middlewares.GenerateTraceID(), middlewares.RequestInfoLogger())
	return engine
}
