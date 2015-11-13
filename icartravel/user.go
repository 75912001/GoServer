package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
)

var GUserMgr UserMgr

type User struct {
	Conn    *net.TCPConn
	Account string
	Uid     USER_ID
}

//todo
func (p *User) Send(cmdId CMD_ID, req proto.Message) (ret int, err error) {
	reqBuf, err := proto.Marshal(req)
	if nil != err {
		fmt.Printf("######proto.Marshal err:", err)
		return ret, err
	}
	var reqBufLen = len(reqBuf)
	var send_buf_all_len uint32
	send_buf_all_len = uint32(reqBufLen + 8)

	head_buf := new(bytes.Buffer)
	var data = []interface{}{
		send_buf_all_len,
		cmdId,
	}
	for _, v := range data {
		err := binary.Write(head_buf, binary.LittleEndian, v)
		if nil != err {
			fmt.Println("binary.Write failed:", err)
		}
	}

	//todo
	ret, err = p.Conn.Write(head_buf.Bytes())
	ret, err = p.Conn.Write(reqBuf)
	if nil != err {
		fmt.Printf("######user.Conn.Write err:", err)
		return ret, err
	}
	fmt.Println("Send body len:", reqBufLen)
	return ret, err
}

type USER_MAP map[*net.TCPConn]User

type UserMgr struct {
	UserMap USER_MAP
}

func (user_mgr *UserMgr) Init() {
	user_mgr.UserMap = make(USER_MAP)
}
