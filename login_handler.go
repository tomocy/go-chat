package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("could not get proivder %s: %s", provider, err)
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("could not get login url of %s: %s", provider, err)
		}

		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("could not get provider %s: %s", provider, err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("could not get user credentials from %s: %s", provider, err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("could not get user from %s: %s", provider, err)
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		userID := fmt.Sprintf("%x", m.Sum(nil))
		authCookie := objx.New(map[string]interface{}{
			"userID":     userID,
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookie,
			Path:  "/",
		})

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not supported action: %s", action)
	}
}
