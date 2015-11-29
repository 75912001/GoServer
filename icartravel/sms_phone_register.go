package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"zzcliredis"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册

//手机验证码个数 5位,[10000-100000)
//手机上5位数字 会有下划线，可以长按复制，方便用户使用
const SmsParamCodeBegin = 10000
const SmsParamCodeEnd = 99999 + 1

//const SmsParamCodeCnt = SmsParamCodeEnd - SmsParamCodeBegin + 1

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
func SmsPhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {

	var recNum string
	{ //解析手机号码
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler")
			return
		}
		if len(req.Form["number"]) > 0 {
			recNum = req.Form["number"][0]
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
		key := gSmsPhoneRegister.SmsGenRedisKey(recNum)
		reply, err := gSmsPhoneRegister.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return
		}
		if !reflect.DeepEqual(reply, nil) {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_SENDING)))
			return
		}
	}

	var SmsParamCode string
	{ //短信内容参数
		index := rand.Intn(SmsParamCodeEnd)
		if index < SmsParamCodeBegin {
			index += SmsParamCodeBegin
		}
		SmsParamCode = strconv.Itoa(index)
		fmt.Println(SmsParamCode)
	}

	var smsParam = "{'code':'" + SmsParamCode + "','product':'" + gSmsPhoneRegister.SmsParamProduct + "'}"

	{ //设置到redis中
		commandName := "setex"
		key := gSmsPhoneRegister.SmsGenRedisKey(recNum)
		timeout := "3600" //一小时超时时间
		_, err := gSmsPhoneRegister.Redis.Conn.Do(commandName, key, timeout, SmsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			return
		}
	}
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = gSmsPhoneRegister.smsGenSign(recNum, smsParam, timeStamp)

	reqUrl, err := gSmsPhoneRegister.smsGenReqUrl(strMd5, timeStamp, recNum, smsParam)
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

	//	result := make([]byte, resp.ContentLength)

	//	result, err = ioutil.ReadAll(resp.Body)
	//	if nil != err {
	//		fmt.Println("######PhoneRegisterHttpHandler.Get err:", err, resp.Body)
	//		return
	//	}
	//	fmt.Println(resp)
	//	fmt.Println(result)
	//	_, err = w.Write(result)
	//	if nil != err {
	//		fmt.Println("######PhoneRegisterHttpHandler...err:", err)
	//	}

	//	fmt.Println("PhoneRegisterHttpHandler end")
}

type SmsPhoneRegister struct {
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
func (p *SmsPhoneRegister) Init() (err error) {
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
func (p *SmsPhoneRegister) smsGenSign(recNum string, smsParam string, timeStamp string) (value string) {
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

	var strMd5 string
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signSource))
	cipherStr := md5Ctx.Sum(nil)
	strMd5 = hex.EncodeToString(cipherStr)
	strMd5 = strings.ToUpper(strMd5)
	return strMd5
}

//生成短信请求url
func (p *SmsPhoneRegister) smsGenReqUrl(strMd5 string, timeStamp string, recNum string, smsParam string) (value string, err error) {
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
func (p *SmsPhoneRegister) SmsGenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
