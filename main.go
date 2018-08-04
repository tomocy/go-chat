package main

import (
	"log"
	"net/http"
)

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}
