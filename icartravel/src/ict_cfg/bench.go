package ict_cfg

import (
	"zzini"
)

var Gbench bench

//bench.ini配置文件
type bench struct {
	FileIni zzini.Ini //ini配置文件
}

//加载配置文件
func (p *bench) Load(path string) (err error) {
	err = p.FileIni.Load(path)
	if nil != err {
		return err
	}
	return err
}
