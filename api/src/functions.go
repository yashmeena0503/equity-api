package src

import (
	"api/sql/db"
	"context"
)

// FIXME: get library for ts, unfinished.
func NewUnid() string {
	return "123"
}

func (a *App) DoesUserExistByPlatform(ctx context.Context, platformName PlatformType, platformID string) (bool, error) {
	switch platformName {
		// NOTE: NEW_PLATFORM
	case PlatformDiscord:
		result, err := a.Queries.DoesUserExistByDiscordID(ctx, platformID)
		if err != nil { return false, Err.CheckUserExistence(err) }
		return result, nil
	case PlatformMinecraft:
		result, err := a.Queries.DoesUserExistByMinecraftID(ctx, platformID)
		if err != nil { return false, Err.CheckUserExistence(err) }
		return result, nil
	}
	return false, Err.UnknownPlatform(string(platformName))
}

func (a *App) NewUserByPlatform(ctx context.Context, platformName PlatformType, platformID string) (string, error) {
	unid := NewUnid()

	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil { return "", Err.BeginTx(err) }

	qtx := a.Queries.WithTx(tx)
	if err := qtx.NewUser(ctx, unid); err != nil {
		_ = tx.Rollback()
		return "", Err.NewUser(err)
	}

	switch platformName {
		// NOTE: NEW_PLATFORM
		case PlatformDiscord: err = qtx.NewDiscordUser(ctx, db.NewDiscordUserParams { Unid: unid, DiscordID: platformID })
		case PlatformMinecraft: err = qtx.NewMinecraftUser(ctx, db.NewMinecraftUserParams { Unid: unid, MinecraftID: platformID })
		default: _ = tx.Rollback(); return "", Err.UnknownPlatform(string(platformName))
	}
	if err != nil {
		_ = tx.Rollback()
		return "", Err.NewUserByPlatform(err)
	}
	if err := tx.Commit(); err != nil { return "", Err.CommitTx(err) }
	return unid, nil
}
