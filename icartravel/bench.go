package main

import (
	"runtime"
	"strconv"
	"zzcommon"
	"zzini"
)

//bench.ini配置文件
type BenchFile struct {
	FileIni         zzini.ZZIni //ini配置文件
	Ip              string
	Port            uint16
	PacketLengthMax int  //设置包最大
	Delay           bool //tcp中是否延迟,默认为true
	GoProcessMax    int  //并行执行的数量
}

//加载配置文件
func (p *BenchFile) Load() (err error) {
	err = p.FileIni.Load()
	if nil != err {
		return err
	}

	p.Delay = true
	p.Ip = p.FileIni.Get("server", "ip", "")
	p.Port = zzcommon.StringToUint16(p.FileIni.Get("server", "port", "0"))
	p.PacketLengthMax = zzcommon.StringToInt(p.FileIni.Get("common", "packet_length_max", "81920"))
	str_num_cpu := strconv.Itoa(runtime.NumCPU())
	p.GoProcessMax = zzcommon.StringToInt(p.FileIni.Get("common", "go_process_max", str_num_cpu))
	runtime.GOMAXPROCS(p.GoProcessMax)
	return err
}
