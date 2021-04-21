package main

import (
	"net/http"

	"github.com/nihuo9/chitchat/data"
)

// POST /login
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		//http.Redirect(writer, request, "/login", http.StatusSeeOther)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		warning(err, "Cannot find user")
	} else {
		if user.Password == request.PostFormValue("password") {
			session, err := user.CreateSession()
			if err != nil {
				danger(err, "Cannot create session")
			} else {
				cookie := http.Cookie {
					Name: "_cookie",
					Value: session.Uuid,
					HttpOnly: true,
				}
				http.SetCookie(writer, &cookie)
				http.Redirect(writer, request, "/", http.StatusFound)
				return
			}
		} 
	}

	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

// /login
func login(writer http.ResponseWriter, request *http.Request) {
	info("login")
	if request.Method == http.MethodGet {
		generateHTML(writer, nil, "login.layout", "public.navbar", "login")
	} else if request.Method == http.MethodPost {
		authenticate(writer, request)
	}
}

// GET /logout
func logout(writer http.ResponseWriter, request *http.Request) {
	info("logout")
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.Delete()
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

// POST /signup
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	info("signupAccount:" )
	err := request.ParseForm()
	if err != nil {
		warning(err, "ParseForm error")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err == nil && user.Id != 0 {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user = &data.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
		Uuid: data.CreateUUID(),
	}
	if err := user.Export(); err != nil {
		danger(err, "Cannot export user")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

// /signup
func signup(writer http.ResponseWriter, request *http.Request) {
	info("signup,", "meth:", request.Method)
	if request.Method == http.MethodGet {
		generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
	} else if request.Method == http.MethodPost {
		signupAccount(writer, request)
	}
}