package main

type TwitchMessage struct {
	UserID       string
	Username     string
	Text         string
	Channel      string
	PlatformName string
}

type MessageRequest struct {
	PlatformName  string `json:"platform_name"`
	PlatformID    string `json:"platform_id"`
	MessageLength int    `json:"message_length"`
}
