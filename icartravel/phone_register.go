package main

import (
	//	"crypto/md5"
	//	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	//	"math/rand"
	"net/http"
	//	"net/url"
	"reflect"
	"strconv"
	//	"strings"
	//	"time"
	"zzcliredis"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//手机注册

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}

type PhoneRegister struct {
	Pattern string
	//redis
	Redis           zzcliredis.ClientRedis
	RedisKeyPerfix  string
	RedisKeyIncrUid string
}

//初始化
func (p *PhoneRegister) Init() (err error) {
	p.Pattern = gBenchFile.FileIni.Get("phone_register", "Pattern", " ")
	//redis
	p.Redis.RedisIp = gBenchFile.FileIni.Get("phone_register", "redis_ip", " ")
	p.Redis.RedisPort = zzcommon.StringToUint16(gBenchFile.FileIni.Get("phone_register", "redis_port", " "))
	p.Redis.RedisDatabases = zzcommon.StringToInt(gBenchFile.FileIni.Get("phone_register", "redis_databases", " "))
	p.RedisKeyPerfix = gBenchFile.FileIni.Get("phone_register", "redis_key_perfix", " ")
	p.RedisKeyIncrUid = gBenchFile.FileIni.Get("phone_register", "redis_key_incr_uid", " ")

	//链接redis
	dialOption := redis.DialDatabase(p.Redis.RedisDatabases)
	var addrRedis = p.Redis.RedisIp + ":" + strconv.Itoa(int(p.Redis.RedisPort))
	p.Redis.Conn, err = redis.Dial("tcp", addrRedis, dialOption)
	if nil != err {
		fmt.Println("######redis.Dial err:", err)
		return err
	}
	//	defer conn.Close()

	//设置uid自增起始点100000   10w

	{ //检查是否有记录 来自redis
		commandName := "get"
		key := p.RedisKeyIncrUid
		reply, err := p.Redis.Conn.Do(commandName, key)
		if nil != err {
			fmt.Println("######redis get err:", err)
			return err
		}
		if reflect.DeepEqual(reply, nil) {
			//设置uid自增起始点100000   10w
			const uidBegin int = 100000
			commandName := "set"
			key := p.RedisKeyIncrUid
			_, err := p.Redis.Conn.Do(commandName, key, uidBegin)
			if nil != err {
				fmt.Println("######redis set err:", err)
				return err
			}
		}
	}
	return err

}

//生成redis的键值
func (p *PhoneRegister) GenRedisKey(key string) (value string) {
	return p.RedisKeyPerfix + key
}
