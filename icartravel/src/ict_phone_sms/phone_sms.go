package ict_phone_sms

import (
	"fmt"
	"ict_cfg"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机短信
var GphoneSms phoneSms

type phoneSms struct {
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

	SmsFreeSignNameChangePwd string
	SmsTemplateCodeChangePwd string
}

//初始化
func (p *phoneSms) Init() (err error) {
	const benchFileSection string = "ict_phone_sms"

	p.UrlPattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "UrlPattern", " ")
	p.AppKey = ict_cfg.Gbench.FileIni.Get(benchFileSection, "AppKey", " ")
	p.AppSecret = ict_cfg.Gbench.FileIni.Get(benchFileSection, "AppSecret", " ")
	p.Method = ict_cfg.Gbench.FileIni.Get(benchFileSection, "Method", " ")
	p.SignMethod = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SignMethod", " ")
	p.SmsFreeSignName = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsFreeSignName", " ")
	p.SmsTemplateCode = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsTemplateCode", " ")
	p.SmsType = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsType", " ")
	p.Versions = ict_cfg.Gbench.FileIni.Get(benchFileSection, "Versions", " ")
	p.SmsParamProduct = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsParamProduct", " ")

	p.SmsFreeSignNameChangePwd = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsFreeSignNameChangePwd", " ")
	p.SmsTemplateCodeChangePwd = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsTemplateCodeChangePwd", " ")
	return err
}

func (p *phoneSms) GenSmsCode() (value string) {
	//手机验证码个数 5位,[10000-100000)
	//手机上5位数字 会有下划线，可以长按复制，方便用户使用
	const smsParamCodeBegin = 10000
	const smsParamCodeEnd = 99999 + 1
	{ //生成短信内容参数
		index := rand.Intn(smsParamCodeEnd)
		if index < smsParamCodeBegin {
			index += smsParamCodeBegin
		}
		value = strconv.Itoa(index)
		fmt.Println(value)
	}
	return value
}

//生成短信请求url
func (p *phoneSms) GenReqUrl(recNum string, smsParam string, SmsFreeSignName string, SmsTemplateCode string) (value string, err error) {
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = p.genSign(recNum, smsParam, timeStamp, SmsFreeSignName, SmsTemplateCode)

	var reqUrl = p.UrlPattern +
		"?sign=" + strMd5 +
		"&app_key=" + p.AppKey +
		"&method=" + p.Method +
		"&rec_num=" + recNum +
		"&sign_method=" + p.SignMethod +
		"&sms_free_sign_name=" + SmsFreeSignName +
		"&sms_param=" + smsParam +
		"&sms_template_code=" + SmsTemplateCode +
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

//生成sign(MD5)
func (p *phoneSms) genSign(recNum string, smsParam string, timeStamp string, SmsFreeSignName string, SmsTemplateCode string) (value string) {
	var signSource = p.AppSecret +
		"app_key" + p.AppKey +
		"method" + p.Method +
		"rec_num" + recNum +
		"sign_method" + p.SignMethod +
		"sms_free_sign_name" + SmsFreeSignName +
		"sms_param" + smsParam +
		"sms_template_code" + SmsTemplateCode +
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
