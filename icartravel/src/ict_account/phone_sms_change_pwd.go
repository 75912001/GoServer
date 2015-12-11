package ict_account

import (
	//	"fmt"
	//	"math/rand"
	"net/http"
	//"strconv"
	//"time"
	//"zzcommon"
)

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

	var smsParamCode string
	{ //生成短信内容参数
		index := rand.Intn(smsParamCodeEnd)
		if index < smsParamCodeBegin {
			index += smsParamCodeBegin
		}
		smsParamCode = strconv.Itoa(index)
		fmt.Println(smsParamCode)
	}

	var smsParam = "{'code':'" + smsParamCode + "','product':'" + ict_phone_sms.GphoneSms.SmsParamProduct + "'}"

	{ //设置到redis中
		err := GphoneSmsRegister.InsertSmsCode(recNum, smsParamCode)
		if nil != err {
			fmt.Println("######redis setex err:", err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zzcommon.StringSubstr(time.Now().String(), 19)

	var strMd5 string = GphoneSmsRegister.genSign(recNum, smsParam, timeStamp)

	reqUrl, err := GphoneSmsRegister.genReqUrl(strMd5, timeStamp, recNum, smsParam)
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
