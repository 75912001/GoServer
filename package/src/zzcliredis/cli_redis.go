package zzcliredis

import (
	//	"fmt"
	"github.com/garyburd/redigo/redis"
	//	"net"
	"strconv"
	//	"zzcommon"
)

//己方作为客户端
type ClientRedis struct {
	conn redis.Conn
}

//连接
func (p *ClientRedis) Connect(ip string, port uint16, recvBufMax int) (err error) {

	var addr = ip + ":" + strconv.Itoa(int(port))
	p.conn, err = redis.Dial("tcp", addr)
	if nil != err {
		// handle error
	}
	//	defer conn.Close()

	return err
}
