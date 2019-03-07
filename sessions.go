package main

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

var sessions = map[string]session{}

type session struct {
	User         string
	lastActivity time.Time
}

func logOn(w http.ResponseWriter, email, password string) error {
	// getting the user first to make sure that it doesnt error out after putting the user in the map
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	idString := id.String()
	cookie := &http.Cookie{
		Name:  "cpgsess",
		Value: idString,
	}
	http.SetCookie(w, cookie)
	sessions[idString] = session{
		User:         email,
		lastActivity: time.Now(),
	}
	return nil
}

// func checkLoggedIn(w http.ResponseWriter, req *http.Request) {
// 	if !loggedIn(req) {
// 		http.Redirect(w, req, "/login", http.StatusSeeOther)
// 		return
// 	}
// 	// the user is loggedi

// }

func tryLogin(email, password string) bool {
	if email == "ianwhitelaw@att.net" && password == "ben" {
		return true
	}
	return false
}

// loggedIn checks to see if a player is currently logged in
func loggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("cpgsess")
	if err != nil {
		return false
	}
	// fmt.Printf("SESSIONS: %v\tMY COOKIEVAL: %s\n", sessions, cookie.Value)
	var ok bool
	if session, ok := sessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		return true
	}
	return ok
}
