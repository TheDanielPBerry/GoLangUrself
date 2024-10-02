package models


import (
	"github.com/gorilla/sessions"
	"net/http"
)

const SESSKEY = "SESSKEY"
var cookieStore *sessions.CookieStore
var session *sessions.Session


func GetCookieStore(req *http.Request) *sessions.CookieStore {
	if cookieStore == nil {
		sessionKey, err := Config("SESSION", "SESSION_KEY")
		if err != nil {
			panic(err)
		}

		key := []byte(sessionKey)
		cookieStore = sessions.NewCookieStore(key)
	}
	return cookieStore
}


func GetSession(req *http.Request) *sessions.Session {
	if session == nil {
		store := GetCookieStore(req)
		var err error
		session, err = store.Get(req, SESSKEY)
		if err != nil {
			panic(err)
		}
	}

	return session
}

