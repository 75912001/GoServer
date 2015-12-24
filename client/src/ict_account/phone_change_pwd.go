package ict_account

import (
	"fmt"
	"ict_cfg"
	//	"ict_common"
	"ict_user"
	"net/http"
	"strconv"
	"zzcommon"
)

var GphoneChangePwd phoneChangePwd_t

////////////////////////////////////////////////////////////////////////////////
//手机修改密码

func PhoneChangePwdHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraRecNumName string = "number"
	const paraPwdName string = "pwd"
	const paraSmsCodeName string = "sms_code"

	var recNum string
	var pwd string
	var smsCode string

	{ //解析参数
		err := req.ParseForm()
		if nil != err {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}

		//手机号码
		if len(req.Form[paraRecNumName]) > 0 {
			recNum = req.Form[paraRecNumName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		//原始密码
		if len(req.Form[paraPwdName]) > 0 {
			pwd = req.Form[paraPwdName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}
		//sms code
		if len(req.Form[paraSmsCodeName]) > 0 {
			smsCode = req.Form[paraSmsCodeName][0]
		} else {
			fmt.Println("######PhoneRegisterHttpHandler")
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PARAM)))
			return
		}

		fmt.Println(recNum, pwd, smsCode)
	}

	{ //检查手机号是否绑定
		bind, err := GphoneRegister.IsPhoneNumBind(recNum)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			fmt.Println(err)
			return
		}
		if !bind {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_PHONE_NUM_BIND)))
			fmt.Println(err)
			return
		}
	}

	{ //检查是否有短信验证码记录
		exist, err := GphoneSmsChangePwd.IsExistSmsCode(recNum, smsCode)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			fmt.Println(err)
			return
		}
		if !exist {
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SMS_REGISTER_CODE)))
			return
		}
	}

	uid, err := GphoneRegister.Uid(recNum)
	if nil != err {
		w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
		fmt.Println(err)
		return
	}

	{
		err := ict_user.Gbase.Insert(uid, recNum, pwd)
		if nil != err {
			fmt.Println(err)
			w.Write([]byte(strconv.Itoa(zzcommon.ERROR_SYS)))
			return
		}
	}
	{ //删除有短信验证码记录 来自redis
		GphoneSmsChangePwd.Del(recNum)
	}
}

type phoneChangePwd_t struct {
	Pattern string
}

//初始化
func (p *phoneChangePwd_t) Init() (err error) {
	const benchFileSection string = "ict_account"
	p.Pattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "PhoneChangePwdHttpHandlerPattern", " ")
	return err
}
