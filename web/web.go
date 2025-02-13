package web

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht-frame/web/middlewares"
)

type Initializer struct{}

type Web struct {
	Port string
}

var config = Web{}

func (Initializer) Init(cfg interface{}) error {
	config = cfg.(Web)
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
	err := engine.Run(":" + config.Port)
	if err != nil {
		return err
	}
	return nil
}
