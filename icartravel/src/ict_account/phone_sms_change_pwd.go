package ict_account

import (
	"fmt"
	"ict_common"
	"ict_phone_sms"
	"net/http"
	"strconv"
	"zzcommon"
)

var GphoneSmsChangePwd phoneSmsChangePwd

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

	reqUrl, err := ict_phone_sms.GphoneSms.GenReqUrl(recNum, smsParam, ict_phone_sms.GphoneSms.SmsFreeSignNameChangePwd, ict_phone_sms.GphoneSms.SmsTemplateCodeChangePwd)
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

type phoneSmsChangePwd struct {
	Pattern string
	//redis
	redisKeyPerfix string
}

func (p *phoneSmsChangePwd) IsExist(recNum string) (value bool) {
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

func (p *phoneSmsChangePwd) InsertSmsCode(recNum string, smsParamCode string) (err error) {
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

//生成redis的键值
func (p *phoneSmsChangePwd) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}
