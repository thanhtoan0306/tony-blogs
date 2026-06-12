package main

import (
	"crypto/subtle"
	"net/http"
	"strings"
	"time"
)

const (
	uploadPassword     = "admin"
	uploadAuthCookie   = "upload_auth"
	uploadAuthToken    = "granted"
	uploadAuthMaxAge   = 7 * 24 * 3600
)

func uploadAuthed(r *http.Request) bool {
	if c, err := r.Cookie(uploadAuthCookie); err == nil &&
		subtle.ConstantTimeCompare([]byte(c.Value), []byte(uploadAuthToken)) == 1 {
		return true
	}
	pw := strings.TrimSpace(r.Header.Get("X-Upload-Password"))
	return pw != "" && checkUploadPassword(pw)
}

func setUploadAuth(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     uploadAuthCookie,
		Value:    uploadAuthToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   uploadAuthMaxAge,
		Expires:  time.Now().Add(uploadAuthMaxAge * time.Second),
	})
}

func checkUploadPassword(password string) bool {
	return subtle.ConstantTimeCompare([]byte(password), []byte(uploadPassword)) == 1
}
