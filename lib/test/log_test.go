package test

import (
	lib "go-gateway/lib/log"
	"testing"
	"time"
)

//测试日志打点
func TestDefaultLog(t *testing.T) {
	SetUp()
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagMySqlSuccess, map[string]interface{}{
		"sql": "sql",
	})
	time.Sleep(time.Second)
	TearDown()
}

//测试日志实例打点
func TestLogInstance(t *testing.T) {
	nlog:= lib.NewLogger()
	logConf:= lib.LogConfig{
		Level:"trace",
		FW: lib.ConfFileWriter{
			On:true,
			LogPath:"./log_test.log",
			RotateLogPath:"./log_test.log",
			WfLogPath:"./log_test.wf.log",
			RotateWfLogPath:"./log_test.wf.log",
		},
		CW: lib.ConfConsoleWriter{
			On:true,
			Color:true,
		},
	}
	nlog.SetInstanceImpl(logConf)
	nlog.Info("test message")
	nlog.Close()
	time.Sleep(time.Second)
}