package main

import (
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "index.html", struct{}{})
}

func comingsoon(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "comingsoon.html", nil)
}

func aboutUs(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "aboutus.html", nil)
}
