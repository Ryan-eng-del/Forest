package lib

// 此包下，存放配置解析结构体，全局配置示例

import (
	"net"
	"time"
)


const TimeFormat = "2006-01-02 15:04:05"
const SecretKey = "Golang-Gateway-insecure-f4l_+=45qi4sru9&ep9trs_(nhombudb(^an36o%e)w#*b3pp8"
const TokenExpirePeriod = time.Hour * 24 * 7

var (
	LocalIP net.IP
)


// 全局变量
var TimeLocation *time.Location
var BaseConfInstance = &BaseConf{}
var MysqlConfInstance = &MysqlMapConf{}
var RedisConfInstance = &RedisMapConf{}
var ZooKeeperConfInstance = &ZooKeeperMapConf{}


type ZooKeeperMapConf struct {
	Zookeeper struct {
		Server    []string   `mapstructure:"server"`
		PathPrefix string `mapstructure:"path_prefix"`
	} `mapstructure:"zookeeper"`
}

type LogConfFileWriter struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type LogConfConsoleWriter struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

type LogConfig struct {
	LogLevel string `mapstructure:"log_level"`
	FW    LogConfFileWriter    `mapstructure:"file_writer"`
	CW    LogConfConsoleWriter `mapstructure:"console_writer"`
}

type BaseConf struct {
	DebugMode string `mapstructure:"debug_mode"`
	TimeLocation string `mapstructure:"time_location"`
	Base struct {
		DebugMode string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
	Log  LogConfig `mapstructure:"log"`
}

type MysqlMapConf struct {
	List map[string]*MySQLConf `mapstructure:"list"`
}

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}


type RedisMapConf struct {
	List map[string]*RedisConf `mapstructure:"list"`
}

type RedisConf struct {
	DataSourceName  string `mapstructure:"data_source_name"`
	ProxyList    []string `mapstructure:"proxy_list"`
	ConnTimeout  int      `mapstructure:"conn_timeout"`
	ReadTimeout  int      `mapstructure:"read_timeout"`
	WriteTimeout int      `mapstructure:"write_timeout"`
	MaxActive int `mapstructure:"max_active"`
	MaxIdle int `mapstructure:"max_idle"`
}