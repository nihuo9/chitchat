package main

import (
	"net/http"
	"fmt"
	"github.com/nihuo9/chitchat/data"
)

// GET /threads/new
// 得到创建话题的页面
func newThread(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "new.thread")
	}
}

// POST /threads/new
// 创建一个话题
func createThread(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		warning("createThread:", err)
		return
	}

	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user, err := sess.User()
	if err != nil {
		danger("Cannot get user from session:", err)
		return
	}
	topic := request.PostFormValue("topic")
	thread, err := user.CreateThread(topic); 
	if err != nil {
		danger("Cannot create thread:", err)
		return
	}
	url := fmt.Sprintf("/thread/read?id=%s", thread.Uuid)
	http.Redirect(writer, request, url, http.StatusSeeOther)
}

// /threads/new
func routeNewThread(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		newThread(writer, request)
	} else if request.Method == http.MethodPost {
		createThread(writer, request)
	}
}

// GET /thread/read
// 读取一个话题详情
func readThread(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	printl("uuid:", uuid)
	thread, err := data.ThreadByUUID(uuid)

	if err != nil {
		warning("cannot read thread, UUID:", uuid, ":", err)
		sendError(writer, request, "cannot read thread")
		return
	}

	_, err = session(writer, request)
	if err != nil {
		generateHTML(writer, &thread, "layout", "public.navbar", "public.thread")
	} else {
		generateHTML(writer, &thread, "layout", "private.navbar", "private.thread")
	}
}

// POST /thread/post
// 回复一个话题
func postThread(writer http.ResponseWriter, request *http.Request) {
	println("postThread")
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	err = request.ParseForm()
	if err != nil {
		warning("Cannot parse form:", err)
		return
	}
	user, err := sess.User()
	if err != nil {
		warning("Cannot get user from session:", err)
		return
	}
	body := request.PostFormValue("body")
	uuid := request.PostFormValue("uuid")
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		sendError(writer, request, "get thread failed")
		return
	}

	if _, err := user.CreatePost(thread, body); err != nil {
		warning("Cannot create post:", err)
		return
	}
	err = thread.UpdateNumReplies(thread.NumReplies + 1)
	if err != nil {
		warning("Cannot update threads table:", err)
		return
	}
	url := fmt.Sprintf("/thread/read?id=%s", uuid)
	http.Redirect(writer, request, url, http.StatusSeeOther)
}