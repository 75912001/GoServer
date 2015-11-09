// mygo project main.go
package main

import (
	"fmt"
	"net"
	//"time"
	//	"zzcommon"
	//"strconv"
	//"sync"
	//"zzcli"
	//"zzser"
	//	"common_msg"
	//	"game_msg"
	"bytes"
	"encoding/binary"
	"proto"
)

var G_user_mgr User_mgr_t

type User_t struct {
	Conn    *net.TCPConn
	Account string
	Uid     uint32
}

func (user *User_t) Send(cmd_id uint32, req proto.Message) (ret int, err error) {

	req_buf, err := proto.Marshal(req)
	if err != nil {
		fmt.Printf("######proto.Marshal err:", err)
		return ret, err
	}
	var req_buf_len = len(req_buf)
	var send_buf_all_len uint32
	send_buf_all_len = uint32(req_buf_len + 8)

	head_buf := new(bytes.Buffer)
	var data = []interface{}{
		send_buf_all_len,
		cmd_id,
	}
	for _, v := range data {
		err := binary.Write(head_buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}

	//todo
	ret, err = user.Conn.Write(head_buf.Bytes())
	ret, err = user.Conn.Write(req_buf)
	if err != nil {
		fmt.Printf("######user.Conn.Write err:", err)
		return ret, err
	}
	fmt.Println("Send body len:", req_buf_len)
	return ret, err
}

type USER_MAP map[*net.TCPConn]User_t

type User_mgr_t struct {
	User_map USER_MAP
}

func (user_mgr *User_mgr_t) Init() {
	user_mgr.User_map = make(USER_MAP)
}
