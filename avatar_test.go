package main

import "testing"

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
