package models

import "time"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Pin      string `json:"pin"`
}

type Message struct {
	ID        string    `json:"id"`
	FromPin   string    `json:"from_pin"`
	FromUser  string    `json:"from_user"`
	ToPin     string    `json:"to_pin"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type WSMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
