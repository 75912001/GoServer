package main

import (
	"fmt"
	"net/http"
)

const pattern string = "/weather/"

func WeatherHttpHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write(gHttpClientWeather.Result)
	if nil != err {
		fmt.Println("######WeatherHttpHandler...err:", err)
	}
}

type Weather struct {
}
