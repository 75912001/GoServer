package zzcliredis

import (
	//	"fmt"
	"github.com/garyburd/redigo/redis"
	//	"net"
	//	"strconv"
	//	"zzcommon"
)

//己方作为客户端
type ClientRedis struct {
	Conn           redis.Conn
	RedisIp        string
	RedisPort      uint16
	RedisDatabases int
}

//连接
/*
func (p *ClientRedis) Connect(ip string, port uint16) (err error) {

	var addr = ip + ":" + strconv.Itoa(int(port))
	p.conn, err = redis.Dial("tcp", addr)
	if nil != err {
		// handle error
	}
	//	defer conn.Close()

	return err
}
*/
