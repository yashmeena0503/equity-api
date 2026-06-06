package src

import (
	"api/sql/db"
	"context"
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewUlid() (string, error) {
	t := time.Now()
	entropy := rand.Reader
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil { return "", ErrNewUnid(err) }
	return id.String(), nil
}

func (a *App) DoesUserExistByPlatform(ctx context.Context, platformName PlatformType, platformID string) (bool, error) {
	switch platformName {
		// NOTE: NEW_PLATFORM
	case PlatformDiscord:
		result, err := a.Queries.DoesUserExistByDiscordID(ctx, platformID)
		if err != nil { return false, ErrCheckUserExistence(err) }
		return result, nil
	case PlatformMinecraft:
		result, err := a.Queries.DoesUserExistByMinecraftID(ctx, platformID)
		if err != nil { return false, ErrCheckUserExistence(err) }
		return result, nil
	}
	return false, ErrUnknownPlatform(string(platformName))
}

func (a *App) NewUserByPlatform(ctx context.Context, platformName PlatformType, platformID string) (string, error) {
	unid, err := NewUlid()
	if err != nil { return unid, err }

	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil { return "", ErrBeginTx(err) }

	qtx := a.Queries.WithTx(tx)
	if err := qtx.NewUser(ctx, unid); err != nil {
		_ = tx.Rollback()
		return "", ErrNewUser(err)
	}

	switch platformName {
		// NOTE: NEW_PLATFORM
		case PlatformDiscord: err = qtx.NewDiscordUser(ctx, db.NewDiscordUserParams { Unid: unid, DiscordID: platformID })
		case PlatformMinecraft: err = qtx.NewMinecraftUser(ctx, db.NewMinecraftUserParams { Unid: unid, MinecraftID: platformID })
		default: _ = tx.Rollback(); return "", ErrUnknownPlatform(string(platformName))
	}
	if err != nil {
		_ = tx.Rollback()
		return "", ErrNewUserByPlatform(err)
	}
	if err := tx.Commit(); err != nil { return "", ErrCommitTx(err) }
	return unid, nil
}
