package web

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qq754174349/ht-frame/config"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web/middlewares"
)

type AutoConfig struct{}

var cfg = config.WebConfig{}

func (AutoConfig) Init(webCfg *config.AppConfig) error {
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	if webCfg.Active == "pro" {
		gin.SetMode(gin.ReleaseMode)
	}
	cfg = webCfg.Web
	return nil
}

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.Default(opts...)
	engine.Use(middlewares.GenerateTraceID(), middlewares.RequestInfoLogger(), middlewares.Prometheus())
	return engine
}

func Run(regRoutes func(engine *gin.Engine), opts ...gin.OptionFunc) error {
	engine := Default(opts...)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.GET("/health", func(c *gin.Context) {})
	regRoutes(engine)
	err := engine.Run(":" + cfg.Port)
	if err != nil {
		return err
	}
	return nil
}
