package ict_cfg

import (
	"zzini"
)

var Gbench bench_t

//bench.ini配置文件
type bench_t struct {
	FileIni zzini.Ini_t //ini配置文件
}

//加载配置文件
func (p *bench_t) Load(path string) (err error) {
	err = p.FileIni.Load(path)
	if nil != err {
		return err
	}
	return err
}
