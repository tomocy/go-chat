package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

type authRequest struct {
	action   string
	provider string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	authReq, err := parseAuthRequest(r)
	if err != nil {
		log.Fatalf("could not get auth request: %s\n", err)
	}
	switch authReq.action {
	case "login":
		loginURL, err := getLoginURL(authReq.provider)
		if err != nil {
			log.Fatalf("could not get login url: %s\n", err)
		}

		redirect(w, loginURL)
	case "callback":
		provider, err := gomniauth.Provider(authReq.provider)
		if err != nil {
			log.Fatalln("could not get provider %s: %s", authReq.provider, err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("could not get user credentials from %s: %s", authReq.provider, err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("could not get user from %s: %s", authReq.provider, err)
		}
		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalf("could not get avatar url: %s\n", err)
		}
		authCookie := objx.New(map[string]interface{}{
			"userID":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
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
		fmt.Fprintf(w, "not supported action: %s", authReq.action)
	}
}

func parseAuthRequest(r *http.Request) (authRequest, error) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		return authRequest{}, errors.New("the length of segs is not over 4")
	}

	return authRequest{
		action:   segs[2],
		provider: segs[3],
	}, nil
}

func getLoginURL(providerName string) (string, error) {
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		return "", err
	}
	loginURL, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		return "", err
	}

	return loginURL, nil
}

func redirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
