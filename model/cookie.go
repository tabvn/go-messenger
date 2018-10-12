package model

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, name string, value string) {

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)

}

func GetCookie(r *http.Request, name string) (string) {

	cookie, err := r.Cookie(name)
	if err != nil || cookie == nil {
		return ""
	}

	return cookie.Value

}
