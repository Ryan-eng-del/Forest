package lib

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	confLib "go-gateway/lib/conf"
	funcLib "go-gateway/lib/func"

	"github.com/gin-gonic/gin"
)

type Trace struct {
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}

const (
	DLTagUndefind      = "_undef"
	DLTagMySqlFailed   = "_com_mysql_failure"
	DLTagRedisFailed   = "_com_redis_failure"
	DLTagMySqlSuccess  = "_com_mysql_success"
	DLTagRedisSuccess  = "_com_redis_success"
	DLTagThriftFailed  = "_com_thrift_failure"
	DLTagThriftSuccess = "_com_thrift_success"
	DLTagHTTPSuccess   = "_com_http_success"
	DLTagHTTPFailed    = "_com_http_failure"
	DLTagTCPFailed     = "_com_tcp_failure"
	DLTagRequestIn     = "_com_request_in"
	DLTagRequestOut    = "_com_request_out"
)

const (
	_dlTag          = "dltag"
	_traceId        = "traceid"
	_spanId         = "spanid"
	_childSpanId    = "cspanid"
	_dlTagBizPrefix = "_com_"
	_dlTagBizUndef  = "_com_undef"
)


var Log *DLLogger
type TraceString string

type TraceContext struct {
	Trace
	CSpanId string
}


type DLLogger struct {}

func (l *DLLogger) TagInfo(trace *TraceContext, dltag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(dltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	NewSingleLoggerDefault()
	loggerDefault.Info(parseParams(m))
}


func (l *DLLogger) TagWarn(trace *TraceContext, dltag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(dltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	NewSingleLoggerDefault()
	loggerDefault.Warn(parseParams(m))
}

func (l *DLLogger) TagError(trace *TraceContext, dltag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(dltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	NewSingleLoggerDefault()
	loggerDefault.Error(parseParams(m))
}

func (l *DLLogger) TagTrace(trace *TraceContext, dltag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(dltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	NewSingleLoggerDefault()
	loggerDefault.Trace(parseParams(m))
}

func (l *DLLogger) TagDebug(trace *TraceContext, dltag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(dltag)

	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	NewSingleLoggerDefault()
	loggerDefault.Debug(parseParams(m))
}

func checkDLTag(dltag string) string {
	if strings.HasPrefix(dltag, _dlTagBizPrefix) {
		return dltag
	}

	if strings.HasPrefix(dltag, "_com_") {
		return dltag
	}

	if dltag == DLTagUndefind {
		return dltag
	}
	return dltag
}

// 生成业务dltag
func CreateBizDLTag(tagName string) string {
	if tagName == "" {
		return _dlTagBizUndef
	}

	return _dlTagBizPrefix + tagName
}

func parseParams(m map[string]interface{}) string {
	var dlTag string = DLTagUndefind

	if _dlTag, isHave := m["dltag"]; isHave {
		if _value, ok := _dlTag.(string); ok {
			dlTag = _value
		}
	}

	for property, value := range m {
		if property == "dltag" {
			continue
		}

		dlTag = dlTag + "||" + fmt.Sprintf("%v=%+v", property, value)
	}

	dlTag = strings.Trim(dlTag, "\"")
	return dlTag
}

func NewTrace() *TraceContext {
	trace := &TraceContext{}
	trace.TraceId = funcLib.GetTraceId()
	trace.SpanId = funcLib.NewSpanId()
	return trace
}


// 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *TraceContext {
	// 防御
	if c == nil {
		return NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*TraceContext); ok {
			return tc
		}
	}
	return NewTrace()
}

func GetTraceContext(ctx context.Context) *TraceContext {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		traceIntraceContext, exists := ginCtx.Get("trace")
		if !exists {
			return NewTrace()
		}
		traceContext, ok := traceIntraceContext.(*TraceContext)
		if ok {
			return traceContext
		}
		return NewTrace()
	}
	traceContext, ok := ctx.Value("trace").(*TraceContext)
	if ok {
		return traceContext
	}

	return NewTrace()
}


func SetGinTraceContext(c *gin.Context, trace *TraceContext) error {
	if trace == nil || c == nil {
		return errors.New("context is nil")
	}
	c.Set("trace", trace)
	return nil
}



func SetTraceContext(ctx context.Context, trace *TraceContext) context.Context {
	if trace == nil {
		return ctx
	}
	return context.WithValue(ctx, TraceString("trace"), trace)
}


func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(confLib.LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

//错误日志
func ContextWarning(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*TraceContext)
	if !ok{
		traceContext = NewTrace()
	}
	Log.TagWarn(traceContext, dltag, m)
}

//错误日志
func ContextError(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*TraceContext)
	if !ok{
		traceContext = NewTrace()
	}
	Log.TagError(traceContext, dltag, m)
}

//普通日志
func ContextNotice(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*TraceContext)
	if !ok{
		traceContext = NewTrace()
	}
	Log.TagInfo(traceContext, dltag, m)
}

//错误日志
func ComLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	Log.TagError(traceContext, dltag, m)
}

//普通日志
func ComLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	Log.TagInfo(traceContext, dltag, m)
}
