package main

import (
	"log/slog"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func StartTwitch(out chan<- TwitchMessage, platforms map[string][]string) {
	client := twitch.NewAnonymousClient()

	channelPlatforms := make(map[string]string)

	for platform, channels := range platforms {
		for _, channel := range channels {
			channelPlatforms[channel] = platform

			client.Join(channel)

			slog.Debug("joined twitch channel",
				"channel", channel,
				"platform", platform,
			)
		}
	}

	slog.Info("connecting to twitch",
		"channels", len(channelPlatforms),
	)

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		platformName, ok := channelPlatforms[msg.Channel]
		if !ok {
			slog.Warn("received message from unknown channel",
				"channel", msg.Channel,
				"user", msg.User.DisplayName,
			)
			return
		}

		out <- TwitchMessage{
			UserID:       msg.User.ID,
			Username:     msg.User.DisplayName,
			Text:         msg.Message,
			Channel:      msg.Channel,
			PlatformName: platformName,
		}
	})

	if err := client.Connect(); err != nil {
		slog.Error("twitch connection failed",
			"err", err,
		)
		return
	}
}
