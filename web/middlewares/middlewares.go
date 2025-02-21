package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web/prometheus"
	"io"
	"net/http"
	"time"
)

// GenerateTraceID 生成链路 ID
func GenerateTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("traceID", uuid.New().String())
		c.Next()
	}
}

// RequestInfoLogger 请求参数打印中间件
func RequestInfoLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的基本信息
		method := c.Request.Method
		path := c.Request.URL.Path

		// 获取请求头并转化为 JSON
		headers, err := json.Marshal(c.Request.Header)
		if err != nil {
			headers = []byte("{}")
		}

		// 保存请求体数据，并且不影响后续中间件
		var body string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.WithTraceID(c.GetString("traceID")).Errorf("Failed to read request body: %v", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}
			body = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// 获取 traceID
		traceID, _ := c.Get("traceID")

		// 记录请求信息
		logger.WithTraceID(traceID.(string)).Infof(
			"Request: method=%s, path=%s, headers=%s, params=%s, body=%s",
			method, path, string(headers), c.Request.URL.RawQuery, body,
		)

		// 调用下一步处理
		c.Next()

	}
}

// Prometheus 普罗米修斯监控请求
func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 设置路由
		route := c.FullPath()

		// 继续处理请求
		c.Next()

		// 记录请求持续时间
		prometheus.Duration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())

		// 更新请求计数器
		prometheus.Requests.WithLabelValues(c.Request.Method, route).Inc()
	}
}
