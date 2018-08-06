package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
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
		user, err := getProviderUser(authReq.provider, r)
		if err != nil {
			log.Fatalf("could not get provider user: %s\n", err)
		}
		chatUser := getChatUser(user)
		setAuthCookie(w, chatUser)

		redirect(w, "/chat")
	default:
		fmt.Fprintf(w, "unsupported action: %s\n", authReq.action)
		w.WriteHeader(http.StatusNotFound)
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

func getProviderUser(providerName string, r *http.Request) (gomniauthcommon.User, error) {
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		return nil, err
	}
	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		return nil, err
	}
	user, err := provider.GetUser(creds)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getChatUser(providerUser gomniauthcommon.User) *chatUser {
	return &chatUser{
		User:     providerUser,
		uniqueID: getHasedUniqueID(providerUser.Email()),
	}
}

func getHasedUniqueID(s string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(s))

	return fmt.Sprintf("%x", m.Sum(nil))
}

func setAuthCookie(w http.ResponseWriter, u ChatUser) {
	authCookieValue, err := getEncodedAuthCookieValue(u)
	if err != nil {
		log.Fatalf("could not get encoded auth cookie value: %s\n", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})
}

func getEncodedAuthCookieValue(u ChatUser) (string, error) {
	avatarURL, err := avatars.GetAvatarURL(u)
	if err != nil {
		return "", err
	}
	authCookie := objx.New(map[string]interface{}{
		"userID":     u.UniqueID(),
		"name":       "user.Name()",
		"avatar_url": avatarURL,
	}).MustBase64()

	return authCookie, nil
}

func redirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
