package lib

import (
	"bufio"
	"errors"
	lib "go-gateway/lib/func"
	"os"
	"path"
	"time"
)

var TimeFormatFuncMap map[byte]func (*time.Time) int

func init() {
	TimeFormatFuncMap = make(map[byte]func(*time.Time) int, 5)
	TimeFormatFuncMap['Y'] = lib.GetYear
	TimeFormatFuncMap['M'] = lib.GetMonth
	TimeFormatFuncMap['D'] = lib.GetDay
	TimeFormatFuncMap['H'] = lib.GetHour
	TimeFormatFuncMap['m'] = lib.GetMin
}

type FileWriter struct {
	logLevelFloor int
	logLevelCeil  int
	filename      string
	pathFmt       string
	file          *os.File
	fileBufWriter *bufio.Writer
	actions       []func(*time.Time) int
	variables     []interface{}
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (fw *FileWriter) CreateLogFile() error {
	if err := os.MkdirAll(path.Dir(fw.filename), 0755); err != nil {
		return err
	}

	if file, err := os.OpenFile(fw.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		fw.file = file
	}

	if fw.fileBufWriter = bufio.NewWriterSize(fw.file, 8192); fw.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed")
	}

	return nil
}

func (fw *FileWriter) Init() error {
	return fw.CreateLogFile()
}

func (fw *FileWriter) Write(r *Record) error {
	if r.level < fw.logLevelFloor || r.level > fw.logLevelCeil {
		return nil
	}	
	
	if fw.fileBufWriter == nil {
		return errors.New("no opened file")
	}

	if _, err := fw.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) SetFileName(filename string) {
	fw.filename = filename
}


func (w *FileWriter) SetLogLevelCeil(ceil int) {
	w.logLevelCeil = ceil
}

func (w *FileWriter) SetLogLevelFloor(floor int) {
	w.logLevelFloor = floor
}

func (fw *FileWriter) SetPathPattern(pattern string) error {
	n := 0
	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}

	if n == 0 {
		fw.pathFmt = pattern
		return nil
	}

	fw.actions = make([]func(*time.Time) int, 0, n)
	fw.variables = make([]interface{}, n)
	tmp := []byte(pattern)
	// 0 -> %  1 -> Y
	variable := 0

	for _, c := range tmp {
		if c == '%' {
			variable = 1
			continue
		}

		if variable == 1 {
			act, ok := TimeFormatFuncMap[c]
			if !ok {
				return errors.New("invalid time format str, expected like Y M D H m")
			}

			fw.actions = append(fw.actions, act)
			variable = 0
		}
	}

	for i, act := range fw.actions {
		now := time.Now()
		fw.variables[i] = act(&now)
	}

	fw.pathFmt = lib.ConvertPatternToFmt(tmp)
	return nil
}

