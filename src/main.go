package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("ChitChat", version, "started at", config.Address)

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// 路由
	mux.HandleFunc("/", index)
	mux.HandleFunc("/err", err)

	// 开启服务
	server := &http.Server {
		Addr: config.Address,
		Handler: mux,
		ReadTimeout: time.Duration(config.ReadTimeout),
		WriteTimeout: time.Duration(config.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}