package lib

import (
	libConf "go-gateway/lib/conf"
	libFunc "go-gateway/lib/func"
	log "go-gateway/lib/log"
)


type BaseLib struct {
	ConfPath string
}

var BaseLibInstance *BaseLib

func (bL *BaseLib) ParseConfig() error {
	return libFunc.ParseConfigFromFile(bL.ConfPath, libConf.BaseConfInstance)
}


func (bL *BaseLib) InitConf () error {
	if err := bL.ParseConfig(); err != nil {
		return err
	}

	if libConf.BaseConfInstance.DebugMode == "" {
		if libConf.BaseConfInstance.Base.DebugMode == "" {
			 libConf.BaseConfInstance.DebugMode = "debug"
		} else {
			libConf.BaseConfInstance.DebugMode = libConf.BaseConfInstance.Base.DebugMode
		}
	}


	if libConf.BaseConfInstance.TimeLocation == "" {
		if libConf.BaseConfInstance.TimeLocation == "" {
			libConf.BaseConfInstance.TimeLocation = "Asia/Shanghai"
		} else {
			libConf.BaseConfInstance.TimeLocation = libConf.BaseConfInstance.Base.TimeLocation
		}
	}

	logConf := log.LogConfig{
		Level: libConf.BaseConfInstance.Log.LogLevel,
		FW: log.ConfFileWriter{
			On:              libConf.BaseConfInstance.Log.FW.On,
			LogPath:         libConf.BaseConfInstance.Log.FW.LogPath,
			RotateLogPath:   libConf.BaseConfInstance.Log.FW.RotateLogPath,
			WfLogPath:       libConf.BaseConfInstance.Log.FW.WfLogPath,
			RotateWfLogPath: libConf.BaseConfInstance.Log.FW.RotateWfLogPath,
		},
		CW: log.ConfConsoleWriter{
			On:    libConf.BaseConfInstance.Log.CW.On,
			Color: libConf.BaseConfInstance.Log.CW.Color,
		},
	}

	if err := log.SetupLog(logConf); err != nil {
		return err
	}

	log.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func (bL *BaseLib) SetPath(fileName string, ConfEnvPath string)  {
	 bL.ConfPath = ConfEnvPath + "/" + fileName + ".toml"
}