package ict_account

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_cfg"
	"ict_common"
	"ict_user"
	"net/http"
	"strconv"
	"zzcommon"
)

var GphoneRegister phoneRegister_t

////////////////////////////////////////////////////////////////////////////////
//手机注册

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
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
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}

		//手机号码
		if len(req.Form[paraRecNumName]) > 0 {
			recNum = req.Form[paraRecNumName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		//原始密码
		if len(req.Form[paraPwdName]) > 0 {
			pwd = req.Form[paraPwdName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		//sms code
		if len(req.Form[paraSmsCodeName]) > 0 {
			smsCode = req.Form[paraSmsCodeName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}

		fmt.Println(recNum, pwd, smsCode)
	}

	{ //检查手机号是否绑定
		bind, err := GphoneRegister.IsPhoneNumBind(recNum)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			fmt.Println(err)
			return
		} else {
			if bind {
				w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
				fmt.Println(err)
				return
			}
		}
	}

	{ //检查是否有短信验证码记录
		exist, err := GphoneSmsRegister.IsExistSmsCode(recNum, smsCode)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			fmt.Println(err)
			return
		}
		if !exist {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
	}

	//生成uid
	uid, err := ict_user.GuidMgr.GenUid()
	if nil != err {
		w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
		fmt.Println(err)
		return
	}

	{ //插入用户数据
		err := GphoneRegister.Insert(recNum, uid)
		if nil != err {
			fmt.Println(err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}

	{
		err := ict_user.Gbase.Insert(uid, recNum, pwd)
		if nil != err {
			fmt.Println(err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}
	{ //删除有短信验证码记录 来自redis
		GphoneSmsRegister.Del(recNum)
	}

	w.Write([]byte(strconv.Itoa(zzcommon.SUCC)))
}

type phoneRegister_t struct {
	Pattern string
	//redis
	redisKeyPerfix string
}

//初始化
func (p *phoneRegister_t) Init() (err error) {
	const benchFileSection string = "ict_account"
	p.Pattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "PhoneRegisterHttpHandlerPattern", " ")
	//redis
	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_perfix_phone_register", " ")
	return err
}

//生成redis的键值
func (p *phoneRegister_t) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

//手机号是否绑定
func (p *phoneRegister_t) IsPhoneNumBind(recNum string) (bind bool, err error) {
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)

	if nil != err {
		fmt.Println("######IsPhoneNumBind err:", err)
		return false, err
	}
	if nil == reply {
		return false, err
	}
	return true, err
}

func (p *phoneRegister_t) Uid(recNum string) (uid string, err error) {
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)

	if nil != err {
		fmt.Println("######IsPhoneNumBind err:", err)
		return "", err
	}
	uid, err = redis.String(reply, err)
	return uid, err
}

func (p *phoneRegister_t) Insert(recNum string, uid string) (err error) {
	//插入用户数据
	commandName := "set"
	key := p.genRedisKey(recNum)
	_, err = ict_common.GRedisClient.Conn.Do(commandName, key, uid)

	if nil != err {
		fmt.Println("######gPhoneRegister err:", err, uid, recNum)
		return err
	}
	return err
}
