package main

import "github.com/gorilla/websocket"

type client struct {
	socket websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		c.room.forward <- msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
