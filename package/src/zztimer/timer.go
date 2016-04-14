/*
////////////////////////////////////////////////////////////////////////////////
//使用方法
import (
	"zztimer"
)
func main() {
	zztimer.Second(1, timerSecondTest)
}

//定时器,秒,测试
func timerSecondTest() {
	//todo

	//继续循环该定时器
	zztimer.Second(1, timerSecondTest)
}
*/

package zztimer

import (
	"time"
)

//定时器,秒
func Second(value uint32, f func()) *time.Timer {
	v := time.Duration(value)
	return time.AfterFunc(v*time.Second, f)
}
