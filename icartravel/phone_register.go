package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"reflect"
	"strconv"
	"zzcliredis"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	var recNum string
	var pwd string
	var smsCode string

	{ //解析参数
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}
		const paraRecNumName string = "number"
		const paraPwdName string = "pwd"
		const paraSmsCodeName string = "sms_code"

		//手机号码
		if len(req.Form[paraRecNumName]) > 0 {
			recNum = req.Form[paraRecNumName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}
		//原始密码
		if len(req.Form[paraPwdName]) > 0 {
			pwd = req.Form[paraPwdName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}
		//sms code
		if len(req.Form[paraSmsCodeName]) > 0 {
			smsCode = req.Form[paraSmsCodeName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}

		fmt.Println(recNum, pwd, smsCode)
	}

	{ //检查是否有短信验证码记录 来自redis
		commandName := "get"
		key := gSmsPhoneRegister.SmsGenRedisKey(recNum)
		reply, err := gSmsPhoneRegister.Redis.Conn.Do(commandName, key)

		if nil != err {
			fmt.Println("######redis get err:", err)
			return
		}
		if nil == reply {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
		getRecNum, err := redis.String(reply, err)
		if smsCode != getRecNum {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
	}
	{ //检查手机号是否绑定
		hasUid, err := gPhoneRegister.IsPhoneNumBind(recNum)
		if nil != err {
			return
		} else {
			if hasUid {
				w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
				return
			}
		}
	}
	//todo
	//注册用户。。。
	//生成uid
	//插入用户数据

}

type PhoneRegister struct {
	Pattern string
	//redis
	Redis           zzcliredis.ClientRedis
	RedisKeyPerfix  string
	RedisKeyIncrUid string
}

//初始化
func (p *PhoneRegister) Init() (err error) {
	p.Pattern = gBenchFile.FileIni.Get("phone_register", "Pattern", " ")
	//redis
	p.Redis.RedisIp = gBenchFile.FileIni.Get("phone_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(gBenchFile.FileIni.Get("phone_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(gBenchFile.FileIni.Get("phone_register", "redis_databases", " "))
	p.RedisKeyPerfix = gBenchFile.FileIni.Get("phone_register", "redis_key_perfix", " ")
	p.RedisKeyIncrUid = gBenchFile.FileIni.Get("phone_register", "redis_key_incr_uid", " ")

	//链接redis
	dialOption := redis.DialDatabase(p.Redis.RedisDatabases)
	var addrRedis = p.Redis.RedisIp + ":" + strconv.Itoa(int(p.Redis.RedisPort))
	p.Redis.Conn, err = redis.Dial("tcp", addrRedis, dialOption)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	//	defer conn.Close()

	//设置uid自增起始点100000   10w

	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.RedisKeyIncrUid
		reply, err := p.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return err
		}
		if reflect.DeepEqual(reply, nil) {
			//设置uid自增起始点100000   10w
			const uidBegin int = 100000
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

//手机号是否绑定
//todo 前缀 ， 获取方法
func (p *PhoneRegister) IsPhoneNumBind(recNum string) (bind bool, err error) {
	commandName := "get"
	key := p.GenRedisKey(recNum)
	reply, err := p.Redis.Conn.Do(commandName, key)

	if nil != err {
		fmt.Println("######HasUid err:", err)
		return false, err
	}
	if nil == reply {
		return false, err
	}
	return true, err
}

//生成redis的键值
func (p *PhoneRegister) GenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
