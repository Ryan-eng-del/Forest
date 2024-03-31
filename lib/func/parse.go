package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"github.com/spf13/viper"
)

func ParseConfigFromFile(path string, instance interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %v fail, %v", path, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail, %v", err)
	}

	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(data))

	if err := v.Unmarshal(instance); err != nil {
		return fmt.Errorf("parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}