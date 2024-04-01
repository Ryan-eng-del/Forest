package public

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)


type LocalTime time.Time


func (t LocalTime) Unix() int64 {
	tTime := time.Time(t)
	return tTime.Unix()
}


func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
			return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
			*t = LocalTime(value)
			return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}


type LoadType int

func (s LoadType) MarshalJSON() ([]byte, error) {
	loadTypeMap := map[int]string{
		0: "http 服务",
		1: "tcp 服务",
		2: "grpc 服务",
	}
	//! 这里有个坑，如果直接返回 []byte() 会报错，因为不符合 json 可以解析的格式 http, 正确的是 "http"
	if v, ok := loadTypeMap[int(s)]; !ok {
		return json.Marshal(v)
	} else {
		return json.Marshal(v)
	}
}
