package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ht-crm/logger"
	"io"
	"net/http"
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

		// 获取路径参数并转化为 JSON
		params, err := json.Marshal(c.Params)
		if err != nil {
			params = []byte("{}")
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
			method, path, string(headers), string(params), body,
		)

		// 调用下一步处理
		c.Next()
	}
}
