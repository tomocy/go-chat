package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["User"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}
