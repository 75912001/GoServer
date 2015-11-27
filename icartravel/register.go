package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	//	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册
const phoneRegisterPattern string = "/phoneRegister"

//?number=17721027200

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
	const sms_paramCode = "123456"
	const sms_paramProduct = "icartravel"
	const smsParam = "{'code':'" + sms_paramCode + "','product':'" + sms_paramProduct + "'}"

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

	fmt.Println("PhoneRegisterHttpHandler end")
}

type PhoneRegister struct {
	UrlPattern      string
	AppKey          string
	AppSecret       string
	Method          string
	SignMethod      string
	SmsFreeSignName string
	SmsTemplateCode string
	SmsType         string
	V               string
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
		"v" + p.V +
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
		"&v=" + p.V

	url, err := url.Parse(reqUrl)
	if nil != err {
		fmt.Println("######PhoneRegister.genReqUrl err:", reqUrl, err)
		return reqUrl, err
	}
	reqUrl = p.UrlPattern + "?" + url.Query().Encode()
	return reqUrl, err
}
