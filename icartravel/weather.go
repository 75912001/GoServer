package main

import (
	"fmt"
	"net/http"
	"time"
	//	"zztimer"
)

const weatherPattern string = "/weather"

func WeatherHttpHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	_, err := w.Write(gHttpClientWeather.Result)
	if nil != err {
		fmt.Println("######WeatherHttpHandler...err:", err)
	}

	time.Sleep(10 * time.Second)
}

type Weather_t struct {
}
