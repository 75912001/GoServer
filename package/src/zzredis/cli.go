package zzredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

//己方作为客户端
type Client struct {
	Conn           redis.Conn
	redisIp        string
	redisPort      uint16
	redisDatabases int
}

//连接
func (p *Client) Connect(ip string, port uint16, redisDatabases int) (err error) {
	p.redisIp = ip
	p.redisPort = port
	p.redisDatabases = redisDatabases

	var addr = ip + ":" + strconv.Itoa(int(port))
	dialOption := redis.DialDatabase(redisDatabases)

	p.Conn, err = redis.Dial("tcp", addr, dialOption)
	if nil != err {
		// handle error
		fmt.Println("######redis.Dial err:", err, ip, port, redisDatabases)
		return err
	}
	//	defer conn.Close()
	return err
}
