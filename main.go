package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}
