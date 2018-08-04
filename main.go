package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "the application address")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	log.Printf("start listening and serving. port: %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}
