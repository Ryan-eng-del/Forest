package lib

import (
	"bytes"
	"time"
)


func GetYear(now *time.Time) int {
	return now.Year()
}

func GetMonth(now *time.Time) int {
	return int(now.Month())
}

func GetDay(now *time.Time) int {
	return now.Day()
}

func GetHour(now *time.Time) int {
	return now.Hour()
}

func GetMin(now *time.Time) int {
	return now.Minute()
}

func ConvertPatternToFmt(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	return string(pattern)
}