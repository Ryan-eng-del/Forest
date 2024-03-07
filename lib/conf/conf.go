package lib

var BaseConfInstance = &BaseConf{}

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


