package ict_register

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zzcliredis"
	"zzcommon"
)

var GPhoneSms phoneSms

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
		commandName := "get"
		key := GPhoneSms.genRedisKey(recNum)
		reply, err := GPhoneSms.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return
		}
		if nil != reply {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_SENDING)))
			return
		}
	}

	{ //检查手机号是否绑定
		hasUid, err := GPhoneRegister.IsPhoneNumBind(recNum)
		if nil != err {
			return
		} else {
			if hasUid {
				w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
				return
			}
		}
	}

	var SmsParamCode string
	{ //生成短信内容参数
		index := rand.Intn(SmsParamCodeEnd)
		if index < SmsParamCodeBegin {
			index += SmsParamCodeBegin
		}
		SmsParamCode = strconv.Itoa(index)
		fmt.Println(SmsParamCode)
	}

	var smsParam = "{'code':'" + SmsParamCode + "','product':'" + gPhoneSms.SmsParamProduct + "'}"

	{ //设置到redis中
		commandName := "setex"
		key := gPhoneSms.GenRedisKey(recNum)
		timeout := "300" //5分钟
		_, err := gPhoneSms.Redis.Conn.Do(commandName, key, timeout, SmsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			return
		}
	}
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = gPhoneSms.genSign(recNum, smsParam, timeStamp)

	reqUrl, err := gPhoneSms.genReqUrl(strMd5, timeStamp, recNum, smsParam)
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

type PhoneSms struct {
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
	Redis          zzcliredis.ClientRedis
	RedisKeyPerfix string
}

//初始化
func (p *PhoneSms) Init() (err error) {
	p.Pattern = gBenchFile.FileIni.Get("sms_phone_register", "Pattern", " ")
	p.UrlPattern = gBenchFile.FileIni.Get("sms_phone_register", "UrlPattern", " ")
	p.AppKey = gBenchFile.FileIni.Get("sms_phone_register", "AppKey", " ")
	p.AppSecret = gBenchFile.FileIni.Get("sms_phone_register", "AppSecret", " ")
	p.Method = gBenchFile.FileIni.Get("sms_phone_register", "Method", " ")
	p.SignMethod = gBenchFile.FileIni.Get("sms_phone_register", "SignMethod", " ")
	p.SmsFreeSignName = gBenchFile.FileIni.Get("sms_phone_register", "SmsFreeSignName", " ")
	p.SmsTemplateCode = gBenchFile.FileIni.Get("sms_phone_register", "SmsTemplateCode", " ")
	p.SmsType = gBenchFile.FileIni.Get("sms_phone_register", "SmsType", " ")
	p.Versions = gBenchFile.FileIni.Get("sms_phone_register", "Versions", " ")
	p.SmsParamProduct = gBenchFile.FileIni.Get("sms_phone_register", "SmsParamProduct", " ")
	//redis
	p.Redis.RedisIp = gBenchFile.FileIni.Get("sms_phone_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(gBenchFile.FileIni.Get("sms_phone_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(gBenchFile.FileIni.Get("sms_phone_register", "redis_databases", " "))
	p.RedisKeyPerfix = gBenchFile.FileIni.Get("sms_phone_register", "redis_key_perfix", " ")

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

//生成sign(MD5)
func (p *PhoneSms) genSign(recNum string, smsParam string, timeStamp string) (value string) {
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
	return strMd5
}

//生成短信请求url
func (p *PhoneSms) genReqUrl(strMd5 string, timeStamp string, recNum string, smsParam string) (value string, err error) {
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
func (p *PhoneSms) genRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}

func (p *PhoneSms) IsExist(recNum string) (value bool) {
	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.genRedisKey(recNum)
		reply, err := p.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return
		}
		if nil != reply {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_SENDING)))
			return
		}
	}
}
