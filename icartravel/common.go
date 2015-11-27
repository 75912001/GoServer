package main

import (
	"sync"
	"zzhttp"
	"zzser"
)

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.HttpClient
var gHttpServer zzhttp.HttpServer
var gPhoneRegister PhoneRegister
var gzzserServer zzser.Server
