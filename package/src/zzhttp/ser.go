package zzhttp

import (
	"fmt"
	"net/http"
	"strconv"
)

type HttpServer struct {
	Ip   string
	Port uint16
}

func (p *HttpServer) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}
func (p *HttpServer) Run() {
	httpAddr := p.Ip + ":" + strconv.Itoa(int(p.Port))
	err := http.ListenAndServe(httpAddr, nil)
	if nil != err {
		fmt.Println("######ListenAndServe: ", err, httpAddr)
	}
}
