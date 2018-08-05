package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

func main() {
	addr := flag.String("addr", ":8080", "the application address")
	flag.Parse()

	gomniauth.SetSecurityKey(gomniauthSecret)
	gomniauth.WithProviders(
		google.New(googleClientKey, googleClientSecret, googleCallbackURL),
	)

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Printf("start listening and serving. port: %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}
