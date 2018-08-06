package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	if err := setAvatarFile(r); err != nil {
		log.Fatalf("could not set avatar file: %s\n", err)
		return
	}
	redirect(w, "/chat")
}

func setAvatarFile(r *http.Request) error {
	userID := r.FormValue("userID")
	if err := deleteAvatarFile(userID); err != nil {
		return err
	}

	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fileName := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	if err := ioutil.WriteFile(fileName, data, 0777); err != nil {
		return err
	}

	return nil
}

func deleteAvatarFile(userID string) error {
	dirName := "avatars"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if match, _ := filepath.Match(userID+"*", file.Name()); match {
			fileName := filepath.Join(dirName, file.Name())
			os.Remove(fileName)
		}
	}

	return nil
}
