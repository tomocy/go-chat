package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan *message
	room   *room
	user   map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		}

		msg.Name = c.user["name"].(string)
		if avatarURL, ok := c.user["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		msg.When = time.Now()
		c.room.forward <- msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
