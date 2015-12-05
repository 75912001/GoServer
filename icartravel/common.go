package main

import (
	"sync"
	"zzhttp"
	"zztcp"
)

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.Client
var gHttpServer zzhttp.Server

var gTcpServer zztcp.Server
