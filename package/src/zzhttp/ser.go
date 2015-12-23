package zzhttp

import (
	"fmt"
	"net/http"
	"strconv"
)

type Server_t struct {
	//	Ip   string
	//	Port uint16
}

func (p *Server_t) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}
func (p *Server_t) Run(ip string, port uint16) {
	httpAddr := ip + ":" + strconv.Itoa(int(port))
	err := http.ListenAndServe(httpAddr, nil)
	if nil != err {
		fmt.Println("######ListenAndServe: ", err, httpAddr)
	}
}
