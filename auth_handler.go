package main

import "net/http"

type authHandler struct {
	next http.Handler
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		if err == http.ErrNoCookie {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		panic(err)
	}

	h.next.ServeHTTP(w, r)
}
