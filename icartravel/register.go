package main

import (
	"fmt"
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册
const phoneRegisterPattern string = "/phoneRegister"

//?number=17721027200

//手机号码长度
const phoneNumberLen int = 11

const app_key = "23273583"
const method = "alibaba.aliqin.fc.sms.num.send"
const sms_free_sign_name = "注册验证"
const sms_template_code = "SMS_2515091"
const sms_type = "normal"
const v = "2.0"

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	var rec_num string
	req.ParseForm()
	if len(req.Form["number"]) > 0 {
		rec_num = req.Form["number"][0]
	}
	if phoneNumberLen != len(rec_num) {
		return
	}

	var sms_paramCode string = "123"
	var sms_paramProduct string = "爱车旅"
	var sms_param string = "{'code':'" + sms_paramCode + "','product':'" + sms_paramProduct + "'}"
	fmt.Println(rec_num)
	fmt.Println(sms_param)

	//const timestamp=2015-11-26+19%3A29%3A56

	/*sign_method 	String 	是 	签名的摘要算法，可选值为：hmac，md5。


	//		http://gw.api.taobao.com/router/rest
			?sign=5A3BF0982B182890A900852CA6076CA9
	//		&timestamp=2015-11-26+19%3A29%3A56
	//		&v=2.0
	//		&app_key=23273583
	//		&method=alibaba.aliqin.fc.sms.num.send
	//		&sms_type=normal
	//		&rec_num=17721027200
	//		&sms_free_sign_name=%E6%B3%A8%E5%86%8C%E9%AA%8C%E8%AF%81
	//		&sms_template_code=SMS_2515091
	//		&sms_param=%7B%27code%27%3A%27123%27%2C%27product%27%3A%27%E7%88%B1%E8%BD%A6%E6%97%85%27%7D


	//		http://gw.api.taobao.com/router/rest
	//		?app_key=23273583
	//		&method=alibaba.aliqin.fc.sms.num.send
	//		&rec_num=17721027200
	//		&sms_free_sign_name=%E6%B3%A8%E5%86%8C%E9%AA%8C%E8%AF%81
	//		&sms_param=%7B%27code%27%3A%27123%27%2C%27product%27%3A%27%E7%88%B1%E8%BD%A6%E6%97%85%27%7D
	//		&sms_template_code=SMS_2515091
	//		&sms_type=normal
	//		&timestamp=2015-11-26+19%3A29%3A56
	//		&v=2.0
	*/
}

type PhoneRegister struct {
}
