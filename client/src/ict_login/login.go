package ict_login

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	//	"zzcommon"
)

var Glogin login_t

const LoginPattern string = "/login"

func LoginHttpHandler(w http.ResponseWriter, req *http.Request) {
	var passWord string = "test md5 encrypto"
	var strMd5 string
	//1
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(passWord))
	cipherStr := md5Ctx.Sum(nil)
	strMd5 = hex.EncodeToString(cipherStr)
	//2
	md5Ctx = md5.New()
	md5Ctx.Write([]byte(strMd5))
	cipherStr = md5Ctx.Sum(nil)
	strMd5 = hex.EncodeToString(cipherStr)

	_, err := w.Write([]byte(strMd5))
	if nil != err {
		fmt.Println("######LoginHttpHandler...err:", err)
	}

	time.Sleep(10 * time.Second)

	fmt.Println(strMd5)
	// 发送给login 服务器
	//异步返回给客户端，要么客户端主动请求服务器（ajax）；要么采用WebSocket连接服务器
}

type login_t struct {
}

//初始化
func (p *login_t) Init() (err error) {
	//	const benchFileSection string = "ict_account"
	//	p.Pattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "PhoneSmsRegisterHttpHandlerPattern", " ")
	//redis
	//	p.redisKeyPerfix = ict_cfg.Gbench.FileIni.Get(benchFileSection, "redis_key_perfix_phone_sms_register", " ")
	return err
}
