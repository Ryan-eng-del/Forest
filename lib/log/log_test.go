package lib

import (
	"testing"
	"time"
)
func TestLogInstance(t *testing.T) {
	nlog:= NewLogger()
	logConf:= LogConfig{
		Level:"trace",
		FW: ConfFileWriter{
			On:true,
			LogPath:"./log_test.log",
			RotateLogPath:"./log_test.log.%Y%M%D%H",
			WfLogPath:"./log_test.wf.log",
			RotateWfLogPath:"./log_test.wf.log.%Y%M%D%H",
		},
		CW: ConfConsoleWriter{
			On:true,
			Color:true,
		},
	}
	nlog.SetInstanceImpl(logConf)

	nlog.Info("test message")
	nlog.Error("test message")
	nlog.Debug("test message")
	nlog.Trace("test message")
	nlog.Warn("test message")
	nlog.Fatal("test message")
	time.Sleep(time.Second * 20)
	nlog.Close()
}