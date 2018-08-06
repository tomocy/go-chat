package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	_, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvadar when client has no data\n")
		t.Errorf("but return %s", err)
	}

	want := "http://to-avatar-url"
	testUser = &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return(want, nil)
	testChatUser.User = testUser
	have, err := authAvatar.GetAvatarURL(testChatUser)
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
	uniqueID := "3a34d1d201c7b30d0143c3d6fbe2b4e5"
	testChatUser := &chatUser{
		uniqueID: uniqueID,
	}
	want := "//www.gravatar.com/avatar/" + uniqueID
	have, err := gravatarAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Errorf("GravatarAvatar.GetAvatarURL should not return any error when client.email is valid\n: %s\n", err)
	}
	if have != want {
		t.Error("GravatarAvatar does not return expected url\n")
		t.Errorf("have %s, but want %s\n", have, want)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	uniqueID := "abc"
	fileName := filepath.Join("avatars", uniqueID+".jpg")
	ioutil.WriteFile(fileName, []byte{}, 0777)
	defer os.Remove(fileName)

	var fileSystemAvatar FileSystemAvatar
	testChatUser := &chatUser{
		uniqueID: uniqueID,
	}
	want := "/" + fileName
	have, err := fileSystemAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Errorf("FileSystemAvatar.GetAvatarURL should not return any error when client.userID is valid: %s\n", err)
	}
	if have != want {
		t.Errorf("have %s, but want %s", have, want)
	}
}
