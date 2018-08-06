package main

import gomniauthcommon "github.com/stretchr/gomniauth/common"

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (c chatUser) UniqueID() string {
	return c.uniqueID
}
