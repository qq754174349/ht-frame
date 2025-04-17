package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/config"
	_ "github.com/qq754174349/ht-frame/consul"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web/middlewares"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AutoConfig struct{}

var cfg = config.WebConfig{}

func init() {
	autoconfigure.Register(AutoConfig{})
}

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
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	regRoutes(engine)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("服务器错误: %v", err)
		}
	}()

	// 处理优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务器关闭错误: %v", err)
	}

	return nil
}
