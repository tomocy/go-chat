package main

import gomniauthcommon "github.com/stretchr/gomniauth/common"

type ChatUser interface {
	UniqueID() string
	Name() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

func (u chatUser) Name() string {
	return u.User.Name()
}
