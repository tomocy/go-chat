package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		log.Fatalf("could not get file and the header: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("could not read all of the file: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	if err := ioutil.WriteFile(fileName, data, 0777); err != nil {
		log.Fatalf("could not write the data to file: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusOK)
}
