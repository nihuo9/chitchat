package main

import (
	"fmt"
	"log"
	"net/http"

	//"time"
)

func main() {
	fmt.Println("ChitChat", version, "started at", config.Address, ",dir:", config.Static)

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// 主页
	mux.HandleFunc("/", index)
	// 错误页面
	mux.HandleFunc("/err", err)

	// 登录
	mux.HandleFunc("/login", login)
	// 退出
	mux.HandleFunc("/logout", logout)
	// 注册
	mux.HandleFunc("/signup", signup)

	// thread
	// 创建一个话题
	mux.HandleFunc("/thread/new", routeNewThread)
	// 读取一个帖子的详情
	mux.HandleFunc("/thread/read", readThread)
	// 回复一个帖子
	mux.HandleFunc("/thread/post", postThread)

	// 开启服务
	server := &http.Server {
		Addr: config.Address,
		Handler: mux,
		/*
		ReadTimeout: time.Duration(config.ReadTimeout),
		WriteTimeout: time.Duration(config.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
		*/
	}

	log.Fatal(server.ListenAndServe())
	/*
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
	*/
}