package main

import (
	"sync"
	"zzhttp"
	"zzser"
)

type PACKET_LENGTH uint32
type MESSAGE_ID uint32
type SESSION_ID uint32
type USER_ID uint32
type RESULT_ID uint32

//消息包头
type ProtoHead struct {
	PacketLength PACKET_LENGTH //总包长度
	MessageId    MESSAGE_ID    //消息号
	SessionId    SESSION_ID    //会话id
	UserId       USER_ID       //用户id
	ResultId     RESULT_ID     //结果id
}

var gLock = &sync.Mutex{}
var gHttpClientWeather zzhttp.HttpClient
var gHttpServer zzhttp.HttpServer

var gzzserServer zzser.Server
