package web

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht-frame/config"
	"github.com/qq754174349/ht-frame/web/middlewares"
)

type Initializer struct{}

var cfg = config.WebConfig{}

func (Initializer) Init(webCfg interface{}) error {
	cfg = webCfg.(config.WebConfig)
	return nil
}

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.Default(opts...)
	engine.Use(middlewares.GenerateTraceID(), middlewares.RequestInfoLogger())
	return engine
}

func Run(regRoutes func(engine *gin.Engine), opts ...gin.OptionFunc) error {
	engine := Default(opts...)
	regRoutes(engine)
	err := engine.Run(":" + cfg.Port)
	if err != nil {
		return err
	}
	return nil
}
