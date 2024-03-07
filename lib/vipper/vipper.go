package lib

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
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
