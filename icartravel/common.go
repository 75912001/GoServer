package main

import (
	"sync"
	"zzhttp"
	"zztcp"
)

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.Client_t
var gHttpServer zzhttp.Server_t

var gTcpServer zztcp.Server_t
