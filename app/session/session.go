package session

import "github.com/gorilla/sessions"

var (
	key = []byte("super-secret")
	Store = sessions.NewCookieStore(key)
)
