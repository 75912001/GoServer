package main

import (
	"fmt"
	"net/http"
	"zzhttp"
)

const pattern string = "/weather/"

var httpClientWeather zzhttp.HttpClient

func WeatherHttpHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write(httpClientWeather.Result)
	if nil != err {
		fmt.Println("######WeatherHttpHandler...err:", err)
	}
}

type Weather struct {
}

func (p *Weather) Register(httpServer *zzhttp.HttpServer) {
	httpServer.AddHandler(pattern, WeatherHttpHandler)
}
