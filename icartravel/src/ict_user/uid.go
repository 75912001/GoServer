package ict_user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_cfg"
	"zzcommon"
	"zzredis"
)

var GuidMgr uidMgr

//设置uid自增起始点100000   10w
const uidBegin int = 100000

////////////////////////////////////////////////////////////////////////////////
//USER ID 管理

type uidMgr struct {
	//redis
	redis           zzredis.Client
	redisKeyIncrUid string
}

//初始化
func (p *uidMgr) Init() (err error) {
	//redis
	ip := ict_cfg.Gbench.FileIni.Get("ict_user_uid", "redis_ip", " ")
	port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("ict_user_uid", "redis_port", " "))
	redisDatabases := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("ict_user_uid", "redis_databases", " "))
	p.redisKeyIncrUid = ict_cfg.Gbench.FileIni.Get("ict_user_uid", "redis_key_incr_uid", " ")

	//链接redis
	err = p.redis.Connect(ip, port, redisDatabases)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}

	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.redisKeyIncrUid
		reply, err := p.redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return err
		}
		if nil == reply {
			commandName := "set"
			key := p.redisKeyIncrUid
			_, err := p.redis.Conn.Do(commandName, key, uidBegin)
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
		key := p.redisKeyIncrUid
		reply, err := p.redis.Conn.Do(commandName, key)
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
