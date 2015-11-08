package zzcommon

import (
	"strconv"
)

const (
	//成功
	SUCC int = 0
	//错误
	ERROR int = -1
	//断开对方的连接
	ERROR_DISCONNECT_PEER int = -2
)

func StringToUint32(s string) (value uint32) {
	vaule, err := strconv.ParseUint(s, 10, 32)
	if nil != err {
		return 0
	}
	return uint32(vaule)
}

func StringToInt(s string) (value int) {
	vaule, err := strconv.ParseInt(s, 10, 0)
	if nil != err {
		return 0
	}
	return int(vaule)
}

func StringToUint16(s string) (value uint16) {
	vaule, err := strconv.ParseUint(s, 10, 16)
	if nil != err {
		return 0
	}
	return uint16(vaule)
}
