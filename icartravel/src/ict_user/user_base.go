package ict_user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_bench_file"
	"strconv"
	"zzcliredis"
	"zzcommon"
)

var GregisterInfo registerInfo

////////////////////////////////////////////////////////////////////////////////
//用户注册信息

type registerInfo struct {
	//redis
	Redis          zzcliredis.ClientRedis
	RedisKeyPerfix string
}

//初始化
func (p *registerInfo) Init() (err error) {
	//redis
	p.Redis.RedisIp = ict_bench_file.GBenchFile.FileIni.Get("user_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(ict_bench_file.GBenchFile.FileIni.Get("user_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(ict_bench_file.GBenchFile.FileIni.Get("user_register", "redis_databases", " "))
	p.RedisKeyPerfix = ict_bench_file.GBenchFile.FileIni.Get("user_register", "redis_key_perfix", " ")

	//链接redis
	dialOption := redis.DialDatabase(p.Redis.RedisDatabases)
	var addrRedis = p.Redis.RedisIp + ":" + strconv.Itoa(int(p.Redis.RedisPort))
	p.Redis.Conn, err = redis.Dial("tcp", addrRedis, dialOption)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	//	defer conn.Close()
	return err

}

//生成redis的键值
func (p *registerInfo) GenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
