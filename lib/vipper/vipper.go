package lib

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	ViperInstance *ViperLib
)

type ViperLib struct {
	ConfEnvPath string
	ConfEnv string
	ConfMap map[string]*viper.Viper
}

func (v *ViperLib) ParseConfPath(config string)  {
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path) - 1], "/")
	v.ConfEnvPath = prefix
	v.ConfEnv = path[len(path) - 2]
}

func (v *ViperLib) LogErr (pos string, err error) {
	log.Printf("[ERROR] %s failed: %v",pos, err)
}

func (v *ViperLib) InitConfig () error {
	f, err := os.Open(v.ConfEnvPath + "/")
	if err != nil {
		v.LogErr("ViperConf.InitConfig", err)
		return err
	}
	
	fileList, err := f.Readdir(1024)
	if err != nil {
		v.LogErr("ViperConf.InitConfig", err)
		return err
	}

	for _, file := range fileList {
		if !file.IsDir() {
			bts, err := os.ReadFile(v.ConfEnvPath + "/" + file.Name())
			if err != nil {
				v.LogErr("ViperConf.InitConfig", err)
				return err
			}
			vI := viper.New()
			vI.SetConfigType("toml")
			vI.ReadConfig(bytes.NewBuffer(bts))

			pathAddr := strings.Split(file.Name(), ".")

			if v.ConfMap == nil {
				v.ConfMap = make(map[string]*viper.Viper)
			}

			v.ConfMap[pathAddr[0]] = vI
		}
	}
	return nil
} 



func (v *ViperLib) GetStringConf (key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	viper, ok := v.ConfMap[keys[0]]

	if !ok {
		return ""
	}
	return viper.GetString(strings.Join(keys[1:], "."))
}

func (v *ViperLib) GetIntConf (key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	viper, ok := v.ConfMap[keys[0]]

	if !ok {
		return 0
	}
	return viper.GetInt(strings.Join(keys[1:], "."))
}

func  (v *ViperLib) GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetStringMap(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func  (v *ViperLib) GetConf(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.Get(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func  (v *ViperLib) GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetBool(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func  (v *ViperLib) GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetFloat64(strings.Join(keys[1:], "."))
	return conf
}


//获取get配置信息
func  (v *ViperLib) GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetStringMapString(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func  (v *ViperLib) GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetStringSlice(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func   (v *ViperLib) GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetTime(strings.Join(keys[1:], "."))
	return conf
}

//获取时间阶段长度
func  (v *ViperLib) GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.GetDuration(strings.Join(keys[1:], "."))
	return conf
}

//是否设置了key
func  (v *ViperLib) IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	viper := v.ConfMap[keys[0]]
	conf := viper.IsSet(strings.Join(keys[1:], "."))
	return conf
}
