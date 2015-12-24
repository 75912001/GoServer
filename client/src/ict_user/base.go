package ict_user

import (
	"fmt"
	"ict_cfg"
	"ict_common"
	"zzcommon"
)

var Gbase base_t

const (
	md5_pwd1_suffix string = "icartravel"
	md5_pwd2_suffix string = "ict"
	field_phone_id  string = "pid"
	field_pwd1      string = "pwd1"
	field_pwd2      string = "pwd2"
)

////////////////////////////////////////////////////////////////////////////////
//用户注册信息

type base_t struct {
	//redis
	redisKeyPerfix string
}

//初始化
func (p *base_t) Init() (err error) {
	const benchFileSection string = "ict_user"
	//redis
	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_perfix_base", " ")

	return err
}

//生成redis的键值
func (p *base_t) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *base_t) Insert(uid string, recNum string, pwd string) (err error) {
	fmt.Println(uid, recNum, pwd)
	{ //注册用户。。。
		//md5
		var pwd1 string = pwd + md5_pwd1_suffix
		var pwd2 string = pwd + md5_pwd2_suffix
		pwd1 = zzcommon.GenMd5(pwd1)
		pwd2 = zzcommon.GenMd5(pwd2)

		commandName := "hmset"
		key := p.genRedisKey(uid)

		_, err = ict_common.GRedisClient.Conn.Do(commandName, key, field_phone_id, recNum, field_pwd1, pwd1, field_pwd2, pwd2)
		if nil != err {
			fmt.Println("######gUserRegister hmset err:", err, uid, recNum, pwd1, pwd2)
			return err
		}
	}
	return err
}
