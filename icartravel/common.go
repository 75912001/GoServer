package main

import (
	"zzcli"
	"zzhttp"
)

//用户
type USER_ID uint32
type CMD_ID uint32

var gHttpClientWeather zzhttp.HttpClient
var gHttpServer zzhttp.HttpServer
var gzzcliClient zzcli.Client
