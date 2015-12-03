package ict_register

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_user"
	"net/http"
	"strconv"
	"zzcliredis"
	"zzcommon"
)

var Gphone phone

////////////////////////////////////////////////////////////////////////////////
//手机注册

func PhoneHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraRecNumName string = "number"
	const paraPwdName string = "pwd"
	const paraSmsCodeName string = "sms_code"

	var recNum string
	var pwd string
	var smsCode string

	{ //解析参数
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}

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

	{ //检查手机号是否绑定
		bind, err := Gphone.IsPhoneNumBind(recNum)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
			return
		} else {
			if bind {
				w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
				return
			}
		}
	}

	{ //检查是否有短信验证码记录 来自redis
		commandName := "get"
		key := GphoneSms.GenRedisKey(recNum)
		reply, err := GPhoneSms.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return
		}
		if nil == reply {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
		getRecNum, _ := redis.String(reply, err)
		if smsCode != getRecNum {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
	}

	//生成uid
	uid, err := ict_user.GUid.GenUid()
	if nil != err {
		w.Write([]byte(strconv.Itoa(zzcommon.ERROR)))
		return
	}

	{ //插入用户数据
		commandName := "set"
		key := GPhone.GenRedisKey(recNum)
		_, err := GPhone.Redis.Conn.Do(commandName, key, uid)
		if nil != err {
			fmt.Println("######gPhoneRegister err:", err, uid, recNum)
			return
		}
	}

	{ //注册用户。。。
		//md5
		var pwd1 string = pwd + "icartravel"
		var pwd2 string = pwd + "ict"
		pwd1 = zzcommon.GenMd5(pwd1)
		pwd2 = zzcommon.GenMd5(pwd2)

		commandName := "hmset"
		key := ict_user.GregisterInfo.GenRedisKey(recNum)
		_, err := ict_user.GregisterInfo.Redis.Conn.Do(commandName, key, "pid", recNum, "pwd1", pwd1, "pwd2", pwd2)
		if nil != err {
			fmt.Println("######gUserRegister hmset err:", err, uid, recNum, pwd1, pwd2)
			return
		}
	}

	{ //删除有短信验证码记录 来自redis
		commandName := "del"
		key := gUserRegister.GenRedisKey(uid)
		gUserRegister.Redis.Conn.Do(commandName, key)
	}

}

type phone struct {
	Pattern string
	//redis
	Redis          zzcliredis.ClientRedis
	RedisKeyPerfix string
}

//初始化
func (p *phone) Init() (err error) {
	p.Pattern = gBenchFile.FileIni.Get("phone_register", "Pattern", " ")
	//redis
	p.Redis.RedisIp = gBenchFile.FileIni.Get("phone_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(gBenchFile.FileIni.Get("phone_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(gBenchFile.FileIni.Get("phone_register", "redis_databases", " "))
	p.RedisKeyPerfix = gBenchFile.FileIni.Get("phone_register", "redis_key_perfix", " ")

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
func (p *phone) GenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}

//手机号是否绑定
func (p *phone) IsPhoneNumBind(recNum string) (bind bool, err error) {
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
