package main

import "errors"

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
