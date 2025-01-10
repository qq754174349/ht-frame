package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ht-crm/src/ht/config/log"
	"io"
)

// GenTraceId 生成链路id
func GenTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("traceId", uuid.New().String())
		c.Next()
	}
}

// ReqInfoLogger 请求参数打印中间件
func ReqInfoLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的基本信息
		method := c.Request.Method
		path := c.Request.URL.Path

		// 获取请求头并转化为 JSON
		headers, err := json.Marshal(c.Request.Header)
		if err != nil {
			headers = []byte("")
		}

		// 获取路径参数并转化为 JSON
		params, err := json.Marshal(c.Params)
		if err != nil {
			params = []byte("")
		}

		// 保存请求体数据，并且不影响后续中间件
		var body string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Error(err)
			} else {
				body = string(bodyBytes)
			}
			c.Request.Body = io.NopCloser(io.Reader(bytes.NewReader(bodyBytes)))
		}

		// 获取 traceId
		traceId, _ := c.Get("traceId")

		// 记录请求信息
		log.TraceInfof(traceId.(string),
			"\nMethod: %s\nPath: %s\nHeaders: %s\nParams: %s\nBody: %s\n", method, path, string(headers), string(params), body)

		// 调用下一步处理
		c.Next()
	}
}
