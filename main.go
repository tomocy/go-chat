package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

func main() {
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not listen and serve: %s", err)
	}
}

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	t.templ.Execute(w, nil)
}
