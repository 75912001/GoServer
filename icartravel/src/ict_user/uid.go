package ict_user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_bench_file"
	"strconv"
	"zzcliredis"
	"zzcommon"
)

var GUidMgr uidMgr

////////////////////////////////////////////////////////////////////////////////
//USER ID 管理

type uidMgr struct {
	//redis
	Redis           zzcliredis.ClientRedis
	RedisKeyIncrUid string
}

//初始化
func (p *uidMgr) Init() (err error) {
	//设置uid自增起始点100000   10w
	const uidBegin int = 100000

	//redis
	p.Redis.RedisIp = ict_bench_file.GBenchFile.FileIni.Get("uid", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(ict_bench_file.GBenchFile.FileIni.Get("uid", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(ict_bench_file.GBenchFile.FileIni.Get("uid", "redis_databases", " "))
	p.RedisKeyIncrUid = ict_bench_file.GBenchFile.FileIni.Get("uid", "redis_key_incr_uid", " ")

	//链接redis
	dialOption := redis.DialDatabase(p.Redis.RedisDatabases)
	var addrRedis = p.Redis.RedisIp + ":" + strconv.Itoa(int(p.Redis.RedisPort))
	p.Redis.Conn, err = redis.Dial("tcp", addrRedis, dialOption)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	//	defer conn.Close()

	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.RedisKeyIncrUid
		reply, err := p.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return err
		}
		if nil == reply {
			commandName := "set"
			key := p.RedisKeyIncrUid
			_, err := p.Redis.Conn.Do(commandName, key, uidBegin)
			if nil != err {
				fmt.Println("######redis set err:", err)
				return err
			}
		}
	}
	return err
}

//生成uid
func (p *uidMgr) GenUid() (uid string, err error) {
	{ //检查是否有记录 来自redis
		commandName := "incr"
		key := p.RedisKeyIncrUid
		reply, err := p.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis incr err:", err)
			return uid, err
		}
		if nil == reply {
			fmt.Println("######redis incr err:", err)
			return uid, err
		}
		uid, err = redis.String(reply, err)
	}
	return uid, err
}
