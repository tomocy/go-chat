package main

import "net/http"

type authHandler struct {
	next http.Handler
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		if err == http.ErrNoCookie {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		panic(err)
	}

	a.next.ServeHTTP(w, r)
}
