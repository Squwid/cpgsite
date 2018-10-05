package main

import (
	"net/http"

	"github.com/Squwid/CPG-Site/cpgaws"
)

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		inquiry := cpgaws.CreateInquiry(req.FormValue("email"), req.FormValue("phone"), req.FormValue("name"), req.FormValue("details"))
		logger.Printf("new inquiry from %s\n", inquiry.Email)
		logger.Printf("sending inquiry %s\n", inquiry.UUID)
		go inquiry.Send() // TODO: This returns an error but we want to deal with it inside of send method

	}
	tpls.ExecuteTemplate(w, "index.html", struct{}{})
}
