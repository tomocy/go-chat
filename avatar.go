package main

import (
	"errors"
)

var ErrNoAvatarURL = errors.New("could not get avatar url")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct {
}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.user["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil
}

type GravatarAvatar struct {
}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userID, ok := c.user["userID"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return "//www.gravatar.com/avatar/" + userIDStr, nil
}
