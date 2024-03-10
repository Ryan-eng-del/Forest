package middlewares

import (
	"bytes"
	libLog "go-gateway/lib/log"
	viperLog "go-gateway/lib/viper"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestInLog (c *gin.Context) {
	traceContext := libLog.NewTrace()

	if traceId := c.Request.Header.Get("com-header-rid"); traceId != "" {
		traceContext.TraceId = traceId
	}
	if spanId := c.Request.Header.Get("com-header-spanid"); spanId != "" {
		traceContext.SpanId = spanId
	}

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)
	libLog.Log.TagInfo(traceContext, "_com_request_in", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	})
}

func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")

	startExecTime, _ := st.(time.Time)
	libLog.ComLogNotice(c, "_com_request_out", map[string]interface{}{
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	})
}

func RequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if viperLog.ViperInstance.GetBoolConf("base.log.file_writer.on") {
			RequestInLog(ctx)
			defer RequestOutLog(ctx)
		}
		ctx.Next()
	}
}