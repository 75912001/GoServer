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
	redis          zzcliredis.ClientRedis
	redisKeyPerfix string
}

//初始化
func (p *base) Init() (err error) {
	//redis
	ip := ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_ip", " ")
	port := zzcommon.StringToUint16(ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_port", " "))
	redisDatabases := zzcommon.StringToInt(ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_databases", " "))
	p.redisKeyPerfix = ict_bench_file.GBenchFile.FileIni.Get("ict_user_base", "redis_key_perfix", " ")

	//链接redis
	err = p.redis.Connect(ip, port, redisDatabases)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}

	return err
}

//生成redis的键值
func (p *base) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *base) Insert(uid string, recNum string, pwd string) (err error) {
	{ //注册用户。。。
		//md5
		var pwd1 string = pwd + "icartravel"
		var pwd2 string = pwd + "ict"
		pwd1 = zzcommon.GenMd5(pwd1)
		pwd2 = zzcommon.GenMd5(pwd2)

		commandName := "hmset"
		key := p.genRedisKey(uid)
		_, err := p.redis.Conn.Do(commandName, key, "pid", recNum, "pwd1", pwd1, "pwd2", pwd2)
		if nil != err {
			fmt.Println("######gUserRegister hmset err:", err, uid, recNum, pwd1, pwd2)
			return err
		}
	}
	return err
}
