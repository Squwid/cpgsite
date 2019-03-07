package main

import (
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "index.html", nil)
}

func comingsoon(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "comingsoon.html", nil)
}

func aboutUs(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "aboutus.html", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	exeTpl := func(incorrectPass bool) {
		tpls.ExecuteTemplate(w, "login.html", struct {
			IncorrectPassword bool
		}{
			IncorrectPassword: incorrectPass,
		})
	}

	if loggedIn(req) {
		http.Redirect(w, req, "/admin/featured", http.StatusSeeOther)
		return
	}

	// if the user is trying to login
	if req.Method == "POST" {
		reqEmail := req.FormValue("login_email")
		reqPass := req.FormValue("login_password")

		correctLogin := tryLogin(reqEmail, reqPass)
		if !correctLogin {
			// incorrect password == true
			exeTpl(true)
			return
		}
		// the password is correct
		err := logOn(w, reqEmail, reqPass)
		if err != nil {
			logger.Fatalf("error loggin user on %v\n", err)
			return
		}
		http.Redirect(w, req, "/admin/featured", http.StatusSeeOther)
		return
	}
	// incorrect password == false
	exeTpl(false)
}
