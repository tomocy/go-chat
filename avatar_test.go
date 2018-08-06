package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvadar when client has no data\n")
		t.Errorf("but return %s", err)
	}

	want := "http://to-avatar-url"
	client.user = map[string]interface{}{
		"avatar_url": want,
	}
	have, err := authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should not return error when client has data\n")
		t.Errorf("but return %s", err)
	}
	if have != want {
		t.Errorf("have %s, but want %s\n", have, want)
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.user = map[string]interface{}{
		"userID": "3a34d1d201c7b30d0143c3d6fbe2b4e5",
	}
	want := "//www.gravatar.com/avatar/3a34d1d201c7b30d0143c3d6fbe2b4e5"
	have, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Errorf("GravatarAvatar.GetAvatarURL should not return any error when client.email is valid\n: %s\n", err)
	}
	if have != want {
		t.Error("GravatarAvatar does not return expected url\n")
		t.Errorf("have %s, but want %s\n", have, want)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	userID := "abc"
	fileName := filepath.Join("avatars", userID+".jpg")
	ioutil.WriteFile(fileName, []byte{}, 0777)
	defer os.Remove(fileName)

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.user = map[string]interface{}{
		"userID": userID,
	}
	want := "/" + fileName
	have, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Errorf("FileSystemAvatar.GetAvatarURL should not return any error when client.userID is valid: %s\n", err)
	}
	if have != want {
		t.Errorf("have %s, but want %s", have, want)
	}
}
