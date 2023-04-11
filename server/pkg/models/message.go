package models

import "time"

type Message struct {
	Sender    *User
	Type      string
	Content   string
	Timestamp time.Time
}
