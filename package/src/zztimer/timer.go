package zztimer

import (
	"time"
)

type Timer struct {
}

//定时器,秒
func Second(value uint32, f func()) *time.Timer {
	v := time.Duration(value)
	return time.AfterFunc(v*time.Second, f)
}
