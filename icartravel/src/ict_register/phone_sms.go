package ict_register

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"ict_cfg"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zzcommon"
	"zzredis"
)

var GphoneSms phoneSms

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

func PhoneSmsHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraNumber string = "number"

	var recNum string
	{ //解析手机号码
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}
		if len(req.Form[paraNumber]) > 0 {
			recNum = req.Form[paraNumber][0]
		}
		//手机号码长度
		const phoneNumberLen int = 11
		if phoneNumberLen != len(recNum) {
			fmt.Println("######PhoneRegisterHttpHandler", recNum)
			return
		}
		fmt.Println(recNum)
	}

	{ //检查是否有记录 来自redis
		isExist := GphoneSms.IsExist(recNum)
		if isExist {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_SENDING)))
			return
		}
	}

	{ //检查手机号是否绑定
		bind, err := Gphone.IsPhoneNumBind(recNum)
		if nil != err {
			return
		}
		if bind {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
			return
		}
	}

	var smsParamCode string
	{ //生成短信内容参数
		index := rand.Intn(smsParamCodeEnd)
		if index < smsParamCodeBegin {
			index += smsParamCodeBegin
		}
		smsParamCode = strconv.Itoa(index)
		fmt.Println(smsParamCode)
	}

	var smsParam = "{'code':'" + smsParamCode + "','product':'" + GphoneSms.SmsParamProduct + "'}"

	{ //设置到redis中
		err := GphoneSms.InsertSmsCode(recNum, smsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			return
		}
	}
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = GphoneSms.genSign(recNum, smsParam, timeStamp)

	reqUrl, err := GphoneSms.genReqUrl(strMd5, timeStamp, recNum, smsParam)
	if nil != err {
		fmt.Println("######gPhoneRegister.genReqUrl", err)
		return
	}
	fmt.Println(reqUrl)

	{ //发送消息到短信服务器
		resp, err := http.Get(reqUrl)
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler.Get err:", err, reqUrl)
			return
		}
		defer resp.Body.Close()
		fmt.Println(resp)
		//fmt.Println(resp.Body)
	}
}

type phoneSms struct {
	Pattern         string
	UrlPattern      string
	AppKey          string
	AppSecret       string
	Method          string
	SignMethod      string
	SmsFreeSignName string
	SmsTemplateCode string
	SmsType         string
	Versions        string
	SmsParamProduct string
	//redis
	redis          zzredis.Client
	redisKeyPerfix string
}

//初始化
func (p *phoneSms) Init() (err error) {
	p.Pattern = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "Pattern", " ")
	p.UrlPattern = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "UrlPattern", " ")
	p.AppKey = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "AppKey", " ")
	p.AppSecret = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "AppSecret", " ")
	p.Method = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "Method", " ")
	p.SignMethod = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "SignMethod", " ")
	p.SmsFreeSignName = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "SmsFreeSignName", " ")
	p.SmsTemplateCode = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "SmsTemplateCode", " ")
	p.SmsType = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "SmsType", " ")
	p.Versions = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "Versions", " ")
	p.SmsParamProduct = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "SmsParamProduct", " ")
	//redis
	ip := ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "redis_ip", " ")
	port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "redis_port", " "))
	redisDatabases := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "redis_databases", " "))
	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get("ict_register_phone_sms", "redis_key_perfix", " ")

	//链接redis
	err = p.redis.Connect(ip, port, redisDatabases)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	return err
}

//生成sign(MD5)
func (p *phoneSms) genSign(recNum string, smsParam string, timeStamp string) (value string) {
	var signSource = p.AppSecret +
		"app_key" + p.AppKey +
		"method" + p.Method +
		"rec_num" + recNum +
		"sign_method" + p.SignMethod +
		"sms_free_sign_name" + p.SmsFreeSignName +
		"sms_param" + smsParam +
		"sms_template_code" + p.SmsTemplateCode +
		"sms_type" + p.SmsType +
		"timestamp" + timeStamp +
		"v" + p.Versions +
		p.AppSecret
	strMd5 := zzcommon.GenMd5(signSource)
	strMd5 = strings.ToUpper(strMd5)
	fmt.Println(signSource)
	fmt.Println(strMd5)
	return strMd5
}

//生成短信请求url
func (p *phoneSms) genReqUrl(strMd5 string, timeStamp string, recNum string, smsParam string) (value string, err error) {
	var reqUrl = p.UrlPattern +
		"?sign=" + strMd5 +
		"&app_key=" + p.AppKey +
		"&method=" + p.Method +
		"&rec_num=" + recNum +
		"&sign_method=" + p.SignMethod +
		"&sms_free_sign_name=" + p.SmsFreeSignName +
		"&sms_param=" + smsParam +
		"&sms_template_code=" + p.SmsTemplateCode +
		"&sms_type=" + p.SmsType +
		"&timestamp=" + timeStamp +
		"&v=" + p.Versions

	url, err := url.Parse(reqUrl)
	if nil != err {
		fmt.Println("######PhoneRegister.genReqUrl err:", reqUrl, err)
		return reqUrl, err
	}
	reqUrl = p.UrlPattern + "?" + url.Query().Encode()
	return reqUrl, err
}

//生成redis的键值
func (p *phoneSms) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *phoneSms) IsExist(recNum string) (value bool) {
	{ //检查是否有记录
		commandName := "get"
		key := p.genRedisKey(recNum)
		reply, err := p.redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return false
		}
		if nil == reply {
			return false
		}
	}
	return true
}

func (p *phoneSms) InsertSmsCode(recNum string, smsParamCode string) (err error) {
	{ //设置到redis中
		commandName := "setex"
		key := p.genRedisKey(recNum)
		timeout := "3600" //60分钟
		_, err := p.redis.Conn.Do(commandName, key, timeout, smsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
		}
	}
	return err

}
func (p *phoneSms) IsExistSmsCode(recNum string, smsCode string) (exist bool, err error) {
	{ //检查是否有短信验证码记录
		commandName := "get"
		key := p.genRedisKey(recNum)
		reply, err := p.redis.Conn.Do(commandName, key)
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
	}
	return true, err
}

func (p *phoneSms) Del(recNum string) {
	{ //删除有短信验证码记录
		commandName := "del"
		key := p.genRedisKey(recNum)
		p.redis.Conn.Do(commandName, key)
	}
}
