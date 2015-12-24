package ict_account

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_cfg"
	"ict_common"
	"ict_phone_sms"
	"net/http"
	"strconv"
	"zzcommon"
)

var GphoneSmsChangePwd phoneSmsChangePwd_t

////////////////////////////////////////////////////////////////////////////////
//请求改变密码,短信码

func PhoneSmsChangePwdHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraNumber string = "number"

	var recNum string
	{ //解析手机号码
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneSmsChangePwdHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		if len(req.Form[paraNumber]) > 0 {
			recNum = req.Form[paraNumber][0]
		}
		//手机号码长度
		const phoneNumberLen int = 11
		if phoneNumberLen != len(recNum) {
			fmt.Println("######PhoneSmsChangePwdHttpHandler", recNum)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		fmt.Println(recNum)
	}

	{ //检查是否有记录 来自redis
		isExist := GphoneSmsChangePwd.IsExist(recNum)
		if isExist {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_SENDING)))
			return
		}
	}

	{ //检查手机号是否绑定
		bind, err := GphoneRegister.IsPhoneNumBind(recNum)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
		if !bind {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_NO_BIND)))
			return
		}
	}

	//生成短信内容参数
	var smsParamCode string = ict_phone_sms.GphoneSms.GenSmsCode()

	var smsParam = "{'code':'" + smsParamCode + "','product':'" + ict_phone_sms.GphoneSms.SmsParamProduct + "'}"

	{ //设置到redis中
		err := GphoneSmsChangePwd.InsertSmsCode(recNum, smsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}

	reqUrl, err := ict_phone_sms.GphoneSms.GenReqUrl(recNum, smsParam,
		ict_phone_sms.GphoneSms.SmsFreeSignNameChangePwd, ict_phone_sms.GphoneSms.SmsTemplateCodeChangePwd)
	if nil != err {
		fmt.Println("######GphoneSmsChangePwd.genReqUrl", err)
		w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
		return
	}
	fmt.Println(reqUrl)

	{ //发送消息到短信服务器
		resp, err := http.Get(reqUrl)
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler.Get err:", err, reqUrl)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
		defer resp.Body.Close()
		fmt.Println(resp)
		//fmt.Println(resp.Body)
	}
	w.Write([]byte(strconv.Itoa(zzcommon.SUCC)))
}

type phoneSmsChangePwd_t struct {
	Pattern string
	//redis
	redisKeyPerfix string
}

func (p *phoneSmsChangePwd_t) IsExist(recNum string) (value bool) {
	//检查是否有记录
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)
	if nil != err {
		fmt.Println("######redis get err:", err)
		return false
	}
	if nil == reply {
		return false
	}

	return true
}

func (p *phoneSmsChangePwd_t) InsertSmsCode(recNum string, smsParamCode string) (err error) {
	//设置到redis中
	commandName := "setex"
	key := p.genRedisKey(recNum)
	timeout := "3600" //60分钟
	_, err = ict_common.GRedisClient.Conn.Do(commandName, key, timeout, smsParamCode)
	if nil != err {
		fmt.Println("######redis setex err:", err)
	}

	return err
}

//初始化
func (p *phoneSmsChangePwd_t) Init() (err error) {
	const benchFileSection string = "ict_account"
	p.Pattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "PhoneSmsChangePwdHttpHandlerPattern", " ")
	//redis
	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_change_pwd_perfix", " ")
	return err
}

//生成redis的键值
func (p *phoneSmsChangePwd_t) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *phoneSmsChangePwd_t) IsExistSmsCode(recNum string, smsCode string) (exist bool, err error) {
	//检查是否有短信验证码记录
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := ict_common.GRedisClient.Conn.Do(commandName, key)
	if nil != err {
		return false, err
	}
	if nil == reply {
		return false, err
	}
	getRecNum, _ := redis.String(reply, err)
	if smsCode != getRecNum {
		fmt.Println("IsExistSmsCode,", recNum, smsCode, getRecNum)
		return false, err
	}

	return true, err
}

func (p *phoneSmsChangePwd_t) Del(recNum string) {
	//删除有短信验证码记录
	commandName := "del"
	key := p.genRedisKey(recNum)
	ict_common.GRedisClient.Conn.Do(commandName, key)
}
