package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nihuo9/chitchat/data"
)

const version = "0.1"
const _DEBUG = true

func printl(args ...interface{}) {
	fmt.Println(args...)
}

func info(args ...interface{}) {
	if _DEBUG {
		fmt.Print("INFO:")
		printl(args...)
		return
	}
	logger.SetPrefix("INFO")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	if _DEBUG {
		fmt.Print("WARNING:")
		printl(args...)
		return
	}
	logger.SetPrefix("WARNING")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	if _DEBUG {
		fmt.Print("ERROR:")
		printl(args...)
		return
	}
	logger.SetPrefix("ERROR")
	logger.Println(args...)
}

func sendError(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), http.StatusFound)
}

func generateHTML(writer io.Writer, data interface{}, filenames ...string) {
	var files []string
	for _, filename := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", filename))
	}
	t := template.Must(template.ParseFiles(files...))

	t.ExecuteTemplate(writer, "layout", data)
}

func sendFile(writer http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Read file failed!", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err == nil {
		fmt.Fprint(writer, string(data))
	}
}

func session(writer http.ResponseWriter, request *http.Request) (sess *data.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess, err = data.SessionByUUID(cookie.Value)
		if err == nil {
			if ok := sess.CheckValid(); !ok {
				err = errors.New("Invalid session")
			}
			printl("session is valid")
		}
	}
	return
}