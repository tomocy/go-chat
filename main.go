package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

func main() {
	setUpGomniauthProviders()
	setUpRouting()

	addr := getAddress()
	log.Printf("start listening and serving. port: %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}

func setUpGomniauthProviders() {
	gomniauth.SetSecurityKey(gomniauthSecret)
	gomniauth.WithProviders(
		google.New(googleClientKey, googleClientSecret, googleCallbackURL),
	)
}

func setUpRouting() {
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/upload", &templateHandler{fileName: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	r := newRoom()
	http.Handle("/room", r)
	go r.run()
}

func getAddress() string {
	addr := flag.String("addr", ":8080", "the application address")
	flag.Parse()

	return *addr
}
