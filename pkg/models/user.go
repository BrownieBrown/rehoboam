package models

import "nhooyr.io/websocket"

type User struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}
