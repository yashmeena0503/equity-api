package src

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Message.

type response struct {
	Result bool `json:"result"`
	Reason string `json:"reason"`
}

func HandleEndpointMessage(ctx context.Context, app *App, req MessageRequest) (response, error) {

	if req.MessageLength > 500 { req.MessageLength = 500 } // NOTE: hardcoded message length.

	// Handle new Users.
	userExists, err := app.DoesUserExistByPlatform(ctx, req.PlatformName, req.PlatformID)
	if err != nil { return response{}, err }
	if !userExists {
		_, err := app.NewUserByPlatform(ctx, req.PlatformName, req.PlatformID)
		if err != nil { return response{}, err }
		return response { Result: true, Reason: "New user created" }, nil
	}

	// Create a user evaluation object and Queue it.
	app.Payout.Enqueue(req.PlatformName, req.PlatformID, ActivityObjectMessage {
		Timestamp: time.Now(),
		MessageLength: req.MessageLength,
	})

	return response { Result: true, Reason: "Enqueued" }, nil
}

// Payout.

func HandleEndpointPayout(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any {
			"pending_activity_objects": app.Payout.GetPendingActivityObjectsInBuffer(),
			"queue": app.Payout.Snapshot(),
		})
	}
}
