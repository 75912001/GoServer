package zzredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

//己方作为客户端
type Client_t struct {
	Conn      redis.Conn
	ip        string
	port      uint16
	dataBases int
}

//连接
func (p *Client_t) Connect(ip string, port uint16, dataBases int) (err error) {
	p.ip = ip
	p.port = port
	p.dataBases = dataBases

	var addr = ip + ":" + strconv.Itoa(int(port))
	dialOption := redis.DialDatabase(dataBases)

	p.Conn, err = redis.Dial("tcp", addr, dialOption)
	if nil != err {
		fmt.Println("######redis.Dial err:", err, ip, port, dataBases)
		return err
	}
	//	defer conn.Close()
	return err
}
