package main

import (
	"ict_register"
	"sync"
	"zzhttp"
	"zzser"
)

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.HttpClient
var gHttpServer zzhttp.HttpServer

var gzzserServer zzser.Server
