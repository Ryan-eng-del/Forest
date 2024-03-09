package lib

import (
	"context"
	"errors"
	"fmt"
	confLib "go-gateway/lib/conf"
	"time"

	logLib "go-gateway/lib/log"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

//mysql 日志打印类型
var DefaultGormLogger = MysqlGormLogger{
	LogLevel:logger.Info,
	SlowThreshold:200 * time.Millisecond,
}

type MysqlGormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (mgl *MysqlGormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	mgl.LogLevel = logLevel
	return mgl
}

func (mgl *MysqlGormLogger) Info(ctx context.Context, message string, values ...interface{}) {
	trace := logLib.GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	logLib.Log.TagInfo(trace, "_com_mysql_Info", params)
}

func (mgl *MysqlGormLogger) Warn(ctx context.Context, message string, values ...interface{}) {
	trace := logLib.GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	logLib.Log.TagInfo(trace, "_com_mysql_Warn", params)
}

func (mgl *MysqlGormLogger) Error(ctx context.Context, message string, values ...interface{}) {
	trace := logLib.GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	logLib.Log.TagInfo(trace, "_com_mysql_Error", params)
}

func (mgl *MysqlGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	trace := logLib.GetTraceContext(ctx)

	if mgl.LogLevel <= logger.Silent {
		return
	}

	sqlStr, rows := fc()
	currentTime := begin.Format(confLib.TimeFormat)
	elapsed := time.Since(begin)
	msg := map[string]interface{}{
		"FileWithLineNum": utils.FileWithLineNum(),
		"sql":             sqlStr,
		"rows":            "-",
		"proc_time":       float64(elapsed.Milliseconds()),
		"current_time":    currentTime,
	}
	switch {
	case err != nil && mgl.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		msg["err"] = err
		if rows == -1 {
			logLib.Log.TagInfo(trace, "_com_mysql_failure", msg)
		} else {
			msg["rows"] = rows
			logLib.Log.TagInfo(trace, "_com_mysql_failure", msg)
		}
	case elapsed > mgl.SlowThreshold && mgl.SlowThreshold != 0 && mgl.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
		msg["slowLog"] = slowLog
		if rows == -1 {
			logLib.Log.TagInfo(trace, "_com_mysql_success", msg)
		} else {
			msg["rows"] = rows
			logLib.Log.TagInfo(trace, "_com_mysql_success", msg)
		}
	case mgl.LogLevel == logger.Info:
		if rows == -1 {
			logLib.Log.TagInfo(trace, "_com_mysql_success", msg)
		} else {
			msg["rows"] = rows
			logLib.Log.TagInfo(trace, "_com_mysql_success", msg)
		}
	}
}
