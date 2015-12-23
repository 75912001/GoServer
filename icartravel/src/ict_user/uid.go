package ict_user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_cfg"
	"ict_common"
	"strconv"
)

var GuidMgr uidMgr_t

//设置uid自增起始点100000   10w
const uidBegin int = 100000

////////////////////////////////////////////////////////////////////////////////
//USER ID 管理

type uidMgr_t struct {
	//redis
	redisKeyIncrUid string
}

//初始化
func (p *uidMgr_t) Init() (err error) {
	const benchFileSection string = "ict_user"
	//redis
	p.redisKeyIncrUid = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_incr_uid", " ")

	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.redisKeyIncrUid
		reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return err
		}
		if nil == reply {
			commandName := "set"
			key := p.redisKeyIncrUid
			_, err := ict_common.GRedisClient.Conn.Do(commandName, key, uidBegin)
			if nil != err {
				fmt.Println("######redis set err:", err)
				return err
			}
		}
	}
	return err
}

//生成uid
func (p *uidMgr_t) GenUid() (uid string, err error) {
	//检查是否有记录 来自redis
	commandName := "incr"
	key := p.redisKeyIncrUid
	reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)
	if nil != err {
		fmt.Println("######redis incr err:", err)
		return uid, err
	}
	if nil == reply {
		fmt.Println("######redis incr err:", err)
		return uid, err
	}
	uid64, err := redis.Int64(reply, err)

	if nil != err {
		fmt.Println("######redis String err:", err)
		return uid, err
	}

	uid = strconv.FormatInt(uid64, 10)
	return uid, err
}
