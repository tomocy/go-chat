package main

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}
