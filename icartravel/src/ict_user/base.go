package ict_user

import (
	"fmt"
	"ict_bench_file"
	"zzcliredis"
	"zzcommon"
)

var Gbase base

////////////////////////////////////////////////////////////////////////////////
//用户注册信息

type base struct {
	//redis
	Redis          zzcliredis.ClientRedis
	RedisKeyPerfix string
}

//初始化
func (p *base) Init() (err error) {
	//redis
	ip := ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_ip", " ")
	port := zzcommon.StringToUint16(ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_port", " "))
	redisDatabases := zzcommon.StringToInt(ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_databases", " "))
	p.RedisKeyPerfix = ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_key_perfix", " ")

	//链接redis
	err = p.Redis.Connect(ip, port, redisDatabases)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}

	return err

}

//生成redis的键值
func (p *base) genRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
