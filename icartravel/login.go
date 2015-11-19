package main

import (
	"crypto/md5"
	"encoding/hex"
	//	"fmt"

	"net/http"
)

const loginPattern string = "/login/"

func LoginHttpHandler(w http.ResponseWriter, req *http.Request) {
	var passWord string = "test md5 encrypto"
	var strMd5 string
	//1
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(passWord))
	cipherStr := md5Ctx.Sum(nil)
	//	fmt.Println(cipherStr)
	strMd5 = hex.EncodeToString(cipherStr)
	//	fmt.Println(strMd5)
	//2
	md5Ctx = md5.New()
	md5Ctx.Write([]byte(strMd5))
	cipherStr = md5Ctx.Sum(nil)
	//	fmt.Println(cipherStr)
	strMd5 = hex.EncodeToString(cipherStr)
	//	fmt.Println(strMd5)

	// 发送给login 服务器
	//异步返回给客户端，要么客户端主动请求服务器（ajax）；要么采用WebSocket连接服务器
}

type Login struct {
}

//登录服务器的返回包
func LoginCallBack() {
	//发送给请求登录的客户端
}