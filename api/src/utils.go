package src

import (
	"context"
	"fmt"
)

func (a *App) DoesUserExists(ctx context.Context, platformName, platformID string) (bool, error) {
	switch platformName {
		case "discord": return a.Queries.UserExistsByDiscordID(ctx, platformID)
		case "minecraft": return a.Queries.UserExistsByMinecraftID(ctx, platformID)
		default: return false, fmt.Errorf("unknown platform: %q", platformName)
	}
}

func (a *App) DoesUserExistsByUnid(ctx context.Context, platformName, platformID string) (bool, error) {
	return a.Queries.UserExistsByUnid(ctx, platformID)
}