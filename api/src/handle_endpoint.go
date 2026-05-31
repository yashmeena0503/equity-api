package src

import (
	"context"
	"fmt"
	"time"
)

type res struct {
	Result bool `json:"result"`
	Reason string `json:"reason"`
}

type MessageBatch struct {
	Timestamp int64 `json:"timestamp"`
	MessageLength int `json:"message_length"`
}

func HandleEndpointMessage(ctx context.Context, app *App, req MessageRequest) (res, error) {

	if req.MessageLength > 500 { req.MessageLength = 500 } // NOTE: hardcoded message length.

	// Handle new Users.
	userExists, err := app.DoesUserExists(ctx, req.PlatformName, req.PlatformID)
	if err != nil { return res {}, fmt.Errorf("Failed to check user existence: %w", err) }
	if !userExists {  } // TODO: actually handle new users.

	// Create a user evaluation object.
	batch := MessageBatch {
		Timestamp: time.Now().UTC().Unix(),
		MessageLength: req.MessageLength,
	}
	// Queue it.
	app.Payout.Enqueue(req.PlatformName, "message", req.PlatformID, batch)

	return res { Result: true, Reason: "reached end" }, nil
}
