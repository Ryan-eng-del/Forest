package lib

type ConfFileWriter struct {
	On              bool   `toml:"On"`
	LogPath         string `toml:"LogPath"`
	RotateLogPath   string `toml:"RotateLogPath"`
	WfLogPath       string `toml:"WfLogPath"`
	RotateWfLogPath string `toml:"RotateWfLogPath"`
}

type ConfConsoleWriter struct {
	On    bool `toml:"On"`
	Color bool `toml:"Color"`
}

type LogConfig struct {
	Level string            `toml:"LogLevel"`
	FW    ConfFileWriter    `toml:"FileWriter"`
	CW    ConfConsoleWriter `toml:"ConsoleWriter"`
}

func SetupLog(lc LogConfig) error {
	NewSingleLoggerDefault()
	return loggerDefault.SetInstanceImpl(lc)
}
