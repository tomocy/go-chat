package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	log.Println(action, provider)
	switch action {
	case "login":
		log.Println("Todo: login with ", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not supported action: %s", action)
	}
}
