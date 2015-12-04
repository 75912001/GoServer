package ict_register

import (
	"fmt"
	//	"github.com/garyburd/redigo/redis"
	//	"ict_register"
	"ict_bench_file"
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

	{ //检查是否有短信验证码记录
		err := GphoneSms.IsExistSmsCode(recNum, smsCode)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
	}

	//生成uid
	uid, err := ict_user.GuidMgr.GenUid()
	if nil != err {
		w.Write([]byte(strconv.Itoa(zzcommon.ERROR)))
		return
	}

	{ //插入用户数据
		err := Gphone.Insert(recNum, uid)
		if nil != err {
			return
		}
	}

	{
		err := ict_user.Gbase.Insert(uid, recNum, pwd)
		if nil != err {
			return
		}
	}

	{ //删除有短信验证码记录 来自redis
		GphoneSms.Del(recNum)
	}

}

type phone struct {
	Pattern string
	//redis
	redis          zzcliredis.ClientRedis
	redisKeyPerfix string
}

//初始化
func (p *phone) Init() (err error) {
	p.Pattern = ict_bench_file.GBenchFile.FileIni.Get("phone_register", "Pattern", " ")
	//redis
	ip := ict_bench_file.GBenchFile.FileIni.Get("phone_register", "redis_ip", " ")
	port := zzcommon.StringToUint16(ict_bench_file.GBenchFile.FileIni.Get("phone_register", "redis_port", " "))
	redisDatabases := zzcommon.StringToInt(ict_bench_file.GBenchFile.FileIni.Get("phone_register", "redis_databases", " "))
	p.redisKeyPerfix = ict_bench_file.GBenchFile.FileIni.Get("phone_register", "redis_key_perfix", " ")

	//链接redis
	err = p.redis.Connect(ip, port, redisDatabases)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	//	defer conn.Close()
	return err
}

//生成redis的键值
func (p *phone) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

//手机号是否绑定
func (p *phone) IsPhoneNumBind(recNum string) (bind bool, err error) {
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := p.redis.Conn.Do(commandName, key)

	if nil != err {
		fmt.Println("######HasUid err:", err)
		return false, err
	}
	if nil == reply {
		return false, err
	}
	return true, err
}

func (p *phone) Insert(recNum string, uid string) (err error) {
	//插入用户数据
	commandName := "set"
	key := p.genRedisKey(recNum)
	_, err = p.redis.Conn.Do(commandName, key, uid)

	if nil != err {
		fmt.Println("######gPhoneRegister err:", err, uid, recNum)
		return err
	}
	return err
}
