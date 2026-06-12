package main

import (
	"crypto/subtle"
	"net/http"
	"time"
)

const (
	uploadPassword     = "admin"
	uploadAuthCookie   = "upload_auth"
	uploadAuthToken    = "granted"
	uploadAuthMaxAge   = 7 * 24 * 3600
)

func uploadAuthed(r *http.Request) bool {
	c, err := r.Cookie(uploadAuthCookie)
	return err == nil && subtle.ConstantTimeCompare([]byte(c.Value), []byte(uploadAuthToken)) == 1
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
