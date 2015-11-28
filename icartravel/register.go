package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册

//手机验证码个数(4位,[1000-9999] 一个9000个)
const SmsParamCodeBegin = 1000
const SmsParamCodeEnd = 9999
const SmsParamCodeCnt = SmsParamCodeEnd - SmsParamCodeBegin + 1

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
func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	//短信内容参数
	index := rand.Int31n(SmsParamCodeCnt)
	var SmsParamCode = gPhoneRegister.SmsParamCode[index]
	fmt.Println(index, gPhoneRegister.SmsParamCode[index])

	var smsParam = "{'code':'" + SmsParamCode + "','product':'" + gPhoneRegister.SmsParamProduct + "'}"

	//解析手机号码
	var recNum string
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

	//todo 检查是否有记录 来自redis
	//1.有记录就返回，短信已发出，请收到后重试

	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = gPhoneRegister.genSign(recNum, smsParam, timeStamp)

	reqUrl, err := gPhoneRegister.genReqUrl(strMd5, timeStamp, recNum, smsParam)
	if nil != err {
		fmt.Println("######gPhoneRegister.genReqUrl", err)
		return
	}
	fmt.Println(reqUrl)
	//发送消息到短信服务器
	resp, err := http.Get(reqUrl)
	if nil != err {
		fmt.Println("######PhoneRegisterHttpHandler.Get err:", err, reqUrl)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	fmt.Println(resp.Body)
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

type PhoneRegister struct {
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
	SmsParamCode    [SmsParamCodeCnt]string
}

func (p *PhoneRegister) Init() {
	p.Pattern = gBenchFile.FileIni.Get("phone_register", "Pattern", " ")
	p.UrlPattern = gBenchFile.FileIni.Get("phone_register", "UrlPattern", " ")
	p.AppKey = gBenchFile.FileIni.Get("phone_register", "AppKey", " ")
	p.AppSecret = gBenchFile.FileIni.Get("phone_register", "AppSecret", " ")
	p.Method = gBenchFile.FileIni.Get("phone_register", "Method", " ")
	p.SignMethod = gBenchFile.FileIni.Get("phone_register", "SignMethod", " ")
	p.SmsFreeSignName = gBenchFile.FileIni.Get("phone_register", "SmsFreeSignName", " ")
	p.SmsTemplateCode = gBenchFile.FileIni.Get("phone_register", "SmsTemplateCode", " ")
	p.SmsType = gBenchFile.FileIni.Get("phone_register", "SmsType", " ")
	p.Versions = gBenchFile.FileIni.Get("phone_register", "Versions", " ")
	p.SmsParamProduct = gBenchFile.FileIni.Get("phone_register", "SmsParamProduct", " ")

	p.genSmsParamCode()
}

func (p *PhoneRegister) genSmsParamCode() {
	for i := SmsParamCodeBegin; i <= SmsParamCodeEnd; i++ {
		p.SmsParamCode[i-SmsParamCodeBegin] = strconv.Itoa(i)
	}
}

//生成sign(MD5)
func (p *PhoneRegister) genSign(recNum string, smsParam string, timeStamp string) (value string) {
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

func (p *PhoneRegister) genReqUrl(strMd5 string, timeStamp string, recNum string, smsParam string) (value string, err error) {
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
