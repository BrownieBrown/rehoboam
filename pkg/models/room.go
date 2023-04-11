package models

type Room struct {
	ID    string
	Name  string
	Users map[string]*User // Use a map for easier user management
}
