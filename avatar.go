package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatarURL = errors.New("could not get avatar url")

type Avatar interface {
	GetAvatarURL(c ChatUser) (string, error)
}

type TryAvatars []Avatar

func (as TryAvatars) GetAvatarURL(c ChatUser) (string, error) {
	for _, a := range as {
		if url, err := a.GetAvatarURL(c); err == nil {
			return url, nil
		}
	}

	return "", ErrNoAvatarURL
}

type AuthAvatar struct {
}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if url == "" {
		return "", ErrNoAvatarURL
	}

	return url, nil
}

type GravatarAvatar struct {
}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + c.UniqueID(), nil
}

type FileSystemAvatar struct {
}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if match, _ := filepath.Match(c.UniqueID()+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}
	}

	return "", ErrNoAvatarURL
}
