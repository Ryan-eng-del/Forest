package lib

import (
	"fmt"
	"os"
)

type ColorRecord Record

func (r *ColorRecord) String() string {
	switch r.level {
	case TRACE:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%-5s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)
	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%-5s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)
	case INFO:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[32m%-5s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[33m%-5s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[31m%-5s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[35m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)
	}
	return ""
}

type ConsoleWriter struct {
	Color bool
}

func NewConsoleWriter() *ConfConsoleWriter {
	return &ConfConsoleWriter{}
}


func (w *ConfConsoleWriter) Write(r *Record) error {
	if w.Color {
		fmt.Fprint(os.Stdout, ((*ColorRecord)(r)).String())
	} else {
		fmt.Fprint(os.Stdout, r.String())
	}
	return nil
}


func (w *ConfConsoleWriter) Init() error {
	return nil
}


func (w *ConfConsoleWriter) SetColor(color bool)  {
	w.Color = color
}