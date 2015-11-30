package main

import (
	"sync"
	"zzhttp"
	"zzser"
)

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.HttpClient
var gHttpServer zzhttp.HttpServer
var gSmsPhoneRegister SmsPhoneRegister
var gPhoneRegister PhoneRegister
var gUserRegister UserRegister
var gUid Uid
var gzzserServer zzser.Server
