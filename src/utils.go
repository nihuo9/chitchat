package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const version = "0.1"


func sendError(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), http.StatusFound)
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, filename := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", filename))
	}
	t := template.Must(template.ParseFiles(files...))
	t.ExecuteTemplate(writer, "index", data)
}

