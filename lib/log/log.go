package lib

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	setUp = false
	loggerDefault *Logger
)

const tunnelSizeDefault = 1024

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

var LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Flusher interface {
	Flush() error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string{
	return fmt.Sprintf("[%5s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

func (l *Logger) SetInstanceImpl(lc LogConfig) (err error) {
	if lc.FW.On {
		if len(lc.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.SetFileName(lc.FW.LogPath)
			w.SetPathPattern(lc.FW.RotateWfLogPath)
			w.SetLogLevelFloor(TRACE)

			// 设置日志级别层级，分文件打印
			// WfLogPath 有，warn 和 error 打印在 WfLogPath, 其余级别打印在 LogPath
			// WfLogPath 没有，都打印在 LogPath
			if len(lc.FW.WfLogPath) > 0 {
				w.SetLogLevelCeil(INFO)
			} else {
				w.SetLogLevelCeil(ERROR)
			}

			l.RegisterWriter(w)
		}

		if len(lc.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(lc.FW.WfLogPath)
			wfw.SetPathPattern(lc.FW.RotateWfLogPath)
			// 设置日志级别层级，分文件打印
			// warn 和 error 打印在 WfLogPath
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			l.RegisterWriter(wfw)
		}
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetColor(lc.CW.Color)
		l.RegisterWriter(w)
	}

	switch lc.Level {
	case "trace":
		l.SetLevel(TRACE)

	case "debug":
		l.SetLevel(DEBUG)

	case "info":
		l.SetLevel(INFO)

	case "warning":
		l.SetLevel(WARNING)

	case "error":
		l.SetLevel(ERROR)

	case "fatal":
		l.SetLevel(FATAL)

	default:
		err = errors.New("invalid log level")
	}

	return
}


func (l *Logger) RegisterWriter(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}


func (l *Logger) Bootstrap()  {
	var (
		r *Record
		ok bool
	)

	if r, ok = <- l.tunnel; !ok {
		l.c <- true
		return
	}


	for _, w := range l.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second  * 10)

	for {
		select {
		case r, ok = <- l.tunnel:
			if !ok {
				l.c <- true
				return
			}

			for _, w := range l.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}
			l.recordPool.Put(r)
		case <- flushTimer.C:
			for _, w := range l.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}

			flushTimer.Reset(time.Second)

		case <-rotateTimer.C:
			for _, w := range l.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}
 
func NewLogger() *Logger {
	if loggerDefault != nil && setUp == false {
		setUp = true
		return loggerDefault
	}

	l := new(Logger)
	l.writers = []Writer{}
	l.tunnel = make(chan *Record, tunnelSizeDefault)
	l.c = make(chan bool, 2)
	l.level = DEBUG
	l.layout = "2006/01/02 15:04:05"
	l.recordPool = &sync.Pool{
		New:  func() interface{} {
			return &Record{}
		},
	}

	go l.Bootstrap()
	return l
}


func NewSingleLoggerDefault() {
	if !setUp {
	 loggerDefault = NewLogger()
	}
}

func SetLayout(layout string) {
	NewSingleLoggerDefault()
	loggerDefault.layout = layout
}