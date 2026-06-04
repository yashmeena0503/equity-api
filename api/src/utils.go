package src

import (
	"context"
	"fmt"
)

func (a *App) DoesUserExist(ctx context.Context, platformName, platformID string) (bool, error) {
	switch platformName {
		// TODO: create platform enum
		case "discord": return a.Queries.UserExistsByDiscordID(ctx, platformID)
		case "minecraft": return a.Queries.UserExistsByMinecraftID(ctx, platformID)
		default: return false, fmt.Errorf("unknown platform: %q", platformName)
	}
}

func (a *App) DoesUserExistByUnid(ctx context.Context, unid string) (bool, error) {
	return a.Queries.UserExistsByUnid(ctx, unid)
}
