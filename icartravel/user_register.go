package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	//	"net/http"
	//	"reflect"
	"strconv"
	"zzcliredis"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//用户注册

type UserRegister struct {
	//redis
	Redis          zzcliredis.ClientRedis
	RedisKeyPerfix string
}

//初始化
func (p *UserRegister) Init() (err error) {
	//redis
	p.Redis.RedisIp = gBenchFile.FileIni.Get("user_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(gBenchFile.FileIni.Get("user_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(gBenchFile.FileIni.Get("user_register", "redis_databases", " "))
	p.RedisKeyPerfix = gBenchFile.FileIni.Get("user_register", "redis_key_perfix", " ")

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
func (p *UserRegister) GenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
