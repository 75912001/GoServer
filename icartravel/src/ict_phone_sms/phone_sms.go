package ict_phone_sms

import (
	"fmt"
	"ict_cfg"
	"net/url"
	"strings"
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
