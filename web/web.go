package web

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht-frame/config"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web/middlewares"
)

type AutoConfig struct{}

var cfg = config.WebConfig{}

func (AutoConfig) Init(webCfg interface{}) error {
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
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
