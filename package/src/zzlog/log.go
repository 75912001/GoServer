/*
//使用方法
*/
package zzlog

import (
	"log"
	"os"
)

//例子
func exampleFun() (err error) {
	var mylog Log_t
	err = mylog.Init("xxx.log", TRACE)
	if nil != err {
		return
	}
	mylog.Debug("hello word!", 1, 1.2345, "show log")
	mylog.DeInit()
	return

}

//日志等级
const (
	EMERG   int = 1
	CRIT    int = 2
	ERROR   int = 3
	WARNING int = 4
	NOTICE  int = 5
	INFO    int = 6
	DEBUG   int = 7
	TRACE   int = 8
)

type Log_t struct {
	level  int      //日志等级
	file   *os.File //日志文件
	logger *log.Logger
}

//初始化
func (this *Log_t) Init(name string, level int) (err error) {
	this.level = level
	this.file, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if nil != err {
		return
	}
	this.logger = log.New(this.file, "", log.Ldate|log.Ltime|log.Llongfile)

	this.Trace("sfsfsdfds", 4, "dsfdf")
	return
}

//反初始化
func (this *Log_t) DeInit() {
	this.file.Close()
}

/////////////////////////////////////////////////////////////////////////////
//日志方法
//踪迹日志
func (this *Log_t) Trace(format string, v ...interface{}) {
	if this.level < TRACE {
		return
	}
	this.logger.Printf(format, v)
}

//调试日志
func (this *Log_t) Debug(format string, v ...interface{}) {
	if this.level < DEBUG {
		return
	}
	this.logger.Printf(format, v)
}

//报告日志
func (this *Log_t) Info(format string, v ...interface{}) {
	if this.level < INFO {
		return
	}
	this.logger.Printf(format, v)
}

//公告日志
func (this *Log_t) Notice(format string, v ...interface{}) {
	if this.level < NOTICE {
		return
	}
	this.logger.Printf(format, v)
}

//警告日志
func (this *Log_t) Warning(format string, v ...interface{}) {
	if this.level < WARNING {
		return
	}
	this.logger.Printf(format, v)
}

//错误日志
func (this *Log_t) Error(format string, v ...interface{}) {
	if this.level < ERROR {
		return
	}
	this.logger.Printf(format, v)
}

//临界日志
func (this *Log_t) Crit(format string, v ...interface{}) {
	if this.level < CRIT {
		return
	}
	this.logger.Printf(format, v)
}

//不可用日志
func (this *Log_t) Emerg(format string, v ...interface{}) {
	if this.level < EMERG {
		return
	}
	this.logger.Printf(format, v)
}
