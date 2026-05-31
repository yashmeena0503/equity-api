package fn

import (
	"api/src"
	"context"
)

type res struct {
	Result bool `json:"result"`
}

func Message(ctx context.Context, app *src.App, req src.MessageRequest) (res, error) {
	return res { Result: true }, nil
}