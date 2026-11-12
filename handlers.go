package main

import (
	_ "embed"
	"net/http"
)

//go:embed public/index.html
var indexPage []byte

var handlers map[string]http.HandlerFunc = map[string]http.HandlerFunc{
	"GET /":           GetHandle,
	"GET /index.php":  GetIndexPhp,
	"POST /index.php": PostHandle,
}

var phpHeaders map[string]string = map[string]string{
	"Server":       "apache/2.4.68",
	"Content-Type": "text/html",
	"X-Powered-By": "PHP/8.4.26",
}

func GetHandle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.php", http.StatusSeeOther)
}

func GetIndexPhp(w http.ResponseWriter, r *http.Request) {
	for k, v := range phpHeaders {
		w.Header().Add(k, v)
	}
	w.Write(indexPage)
}

func PostHandle(w http.ResponseWriter, r *http.Request) {}
