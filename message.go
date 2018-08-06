package main

import "time"

type message struct {
	Name      string
	AvatorURL string
	Message   string
	When      time.Time
}
