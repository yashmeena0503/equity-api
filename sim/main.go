package main

import (
	"log/slog"
	"os"
)

func main() {
	config, err := NewConfig("config.json")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	api := NewAPIClient(config.BaseURL)

	messages := make(chan TwitchMessage, 1024)

	go StartTwitch(messages, config.Platforms)
	slog.Info("message processing started")

	for msg := range messages {
		slog.Info("message received",
			"channel", msg.Channel,
			"platform", msg.PlatformName,
			"user", msg.Username,
			"length", len(msg.Text),
		)

		req := MessageRequest{
			PlatformName:  msg.PlatformName,
			PlatformID:    msg.UserID,
			MessageLength: len(msg.Text),
		}

		if err := api.SendMessage(req); err != nil {
			slog.Error("failed to send message",
				"channel", msg.Channel,
				"platform", msg.PlatformName,
				"user", msg.Username,
				"err", err,
			)
			continue
		}
	}
}
