package models


import (
	"github.com/gorilla/sessions"
	"net/http"
)

const SESSKEY = "SESSKEY"


func GetCookieStore(req *http.Request) *sessions.CookieStore {
	sessionKey, err := Config("SESSION", "SESSION_KEY")
	if err != nil {
		panic(err)
	}

	key := []byte(sessionKey)
	cookieStore := sessions.NewCookieStore(key)
	return cookieStore
}


func GetSession(req *http.Request) *sessions.Session {
	store := GetCookieStore(req)
	session, err := store.Get(req, SESSKEY)
	if err != nil {
		panic(err)
	}

	return session
}

