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
	addr := flag.String("addr", ":8080", "the application address")
	flag.Parse()

	setUpGomniauthProviders()

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/upload", &templateHandler{fileName: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/room", r)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	go r.run()

	log.Printf("start listening and serving. port: %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}

func setUpGomniauthProviders() {
	gomniauth.SetSecurityKey(gomniauthSecret)
	gomniauth.WithProviders(
		google.New(googleClientKey, googleClientSecret, googleCallbackURL),
	)
}
