package main

import "errors"

var ErrNoAvatarURL = errors.New("could not get avatar url")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}
