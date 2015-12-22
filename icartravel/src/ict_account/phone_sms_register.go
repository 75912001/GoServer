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

var GphoneSmsRegister phoneSmsRegister

////////////////////////////////////////////////////////////////////////////////
//手机短信注册(发送手机号,接收验证码)

//手机验证码个数 5位,[10000-100000)
//手机上5位数字 会有下划线，可以长按复制，方便用户使用
const smsParamCodeBegin = 10000
const smsParamCodeEnd = 99999 + 1

/*
	http://gw.api.taobao.com/router/rest
	?sign=5A3BF0982B182890A900852CA6076CA9
	&app_key=23273583
	&method=alibaba.aliqin.fc.sms.num.send
	&rec_num=17721027200
	&sign_method=md5
	&sms_free_sign_name=%E6%B3%A8%E5%86%8C%E9%AA%8C%E8%AF%81
	&sms_param=%7B%27code%27%3A%27123%27%2C%27product%27%3A%27%E7%88%B1%E8%BD%A6%E6%97%85%27%7D
	&sms_template_code=SMS_2515091
	&sms_type=normal
	&timestamp=2015-11-26+19%3A29%3A56
	&v=2.0
*/

func PhoneSmsRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraNumber string = "number"

	var recNum string
	{ //解析手机号码
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneSmsRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		if len(req.Form[paraNumber]) > 0 {
			recNum = req.Form[paraNumber][0]
		}
		//手机号码长度
		const phoneNumberLen int = 11
		if phoneNumberLen != len(recNum) {
			fmt.Println("######PhoneRegisterHttpHandler", recNum)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		fmt.Println(recNum)
	}

	{ //检查是否有记录 来自redis
		isExist := GphoneSmsRegister.IsExist(recNum)
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
		if bind {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
			return
		}
	}
	//生成短信内容参数
	var smsParamCode string = ict_phone_sms.GphoneSms.GenSmsCode()

	var smsParam = "{'code':'" + smsParamCode + "','product':'" + ict_phone_sms.GphoneSms.SmsParamProduct + "'}"

	{ //设置到redis中
		err := GphoneSmsRegister.InsertSmsCode(recNum, smsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}

	reqUrl, err := ict_phone_sms.GphoneSms.GenReqUrl(recNum, smsParam, ict_phone_sms.GphoneSms.SmsFreeSignName, ict_phone_sms.GphoneSms.SmsTemplateCode)
	if nil != err {
		fmt.Println("######gPhoneRegister.genReqUrl", err)
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

type phoneSmsRegister struct {
	Pattern string
	//redis
	redisKeyPerfix string
}

//初始化
func (p *phoneSmsRegister) Init() (err error) {
	const benchFileSection string = "ict_account"
	p.Pattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "PhoneSmsRegisterHttpHandlerPattern", " ")
	//redis
	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_perfix_phone_sms_register", " ")
	return err
}

//生成redis的键值
func (p *phoneSmsRegister) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *phoneSmsRegister) IsExist(recNum string) (value bool) {
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

func (p *phoneSmsRegister) InsertSmsCode(recNum string, smsParamCode string) (err error) {
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

func (p *phoneSmsRegister) IsExistSmsCode(recNum string, smsCode string) (exist bool, err error) {
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

func (p *phoneSmsRegister) Del(recNum string) {
	//删除有短信验证码记录
	commandName := "del"
	key := p.genRedisKey(recNum)
	ict_common.GRedisClient.Conn.Do(commandName, key)
}
