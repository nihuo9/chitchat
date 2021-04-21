package main

import (
	//"fmt"
	"net/http"

	data "github.com/nihuo9/chitchat/data"
)

// GET /
func index(writer http.ResponseWriter, request *http.Request) {
	threads, err := data.GetAllThreads(0)
	if err != nil {
		sendError(writer, request, "获取话题失败，请稍后再试")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "index")
		}
	}
}

// GET /err?msg=
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}