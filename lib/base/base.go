package lib

import (
	"bytes"
	"fmt"
	libConf "go-gateway/lib/conf"
	"io"
	"os"

	"github.com/spf13/viper"
)

type BaseLib struct {
	ConfPath string
}


func (bL *BaseLib) ParseConfig() error {
	libConf.BaseConfInstance = &libConf.BaseConf{}
	file, err := os.Open(bL.ConfPath)
	if err != nil {
		return fmt.Errorf("open config %v fail, %v", bL.ConfPath, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail, %v", err)
	}

	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(data))
	
	if err := v.Unmarshal(libConf.BaseConfInstance); err != nil {
		return fmt.Errorf("parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
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
	return nil
}