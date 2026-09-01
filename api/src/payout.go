package src

import (
	"api/sql/db"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// platformName > activityType > platformID > activityObject
type PayoutQueue map[PlatformType]map[ActivityType]map[string][]ActivityObject

type PayoutBuffer struct {
	mu sync.Mutex
	queue PayoutQueue
}
func (b *PayoutBuffer) Enqueue(platform PlatformType, platformID string, activityObject ActivityObject) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.queue[platform] == nil { b.queue[platform] = make(map[ActivityType]map[string][]ActivityObject) }
	t := activityObject.Type()
	if b.queue[platform][t] == nil { b.queue[platform][t] = make(map[string][]ActivityObject) }
	b.queue[platform][t][platformID] = append(b.queue[platform][t][platformID], activityObject)
}
func (b *PayoutBuffer) Drain() PayoutQueue {
	b.mu.Lock()
	defer b.mu.Unlock()
	snapshot := b.queue
	b.queue = make(PayoutQueue)
	return snapshot
}
func (b *PayoutBuffer) Snapshot() PayoutQueue {
    b.mu.Lock()
    defer b.mu.Unlock()
    out := make(PayoutQueue)
    for platform, byType := range b.queue {
        out[platform] = make(map[ActivityType]map[string][]ActivityObject)
        for activityType, byUser := range byType {
            out[platform][activityType] = make(map[string][]ActivityObject)
            for platformID, activityObjects := range byUser {
                cp := make([]ActivityObject, len(activityObjects))
                copy(cp, activityObjects)
                out[platform][activityType][platformID] = cp
            }
        }
    }
    return out
}

// Payout cycle stuff.

type PayoutOptions struct {
	// Minimum Pending Acitivty objects.
	MinPending int // 0 = disable threshold
}

// Pending Acitivty Objects in PayoutBuffer.
func (b *PayoutBuffer) GetPendingActivityObjectsInBuffer() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for _, byType := range b.queue {
		for _, byUser := range byType {
			for _, activityObject := range byUser {
				n += len(activityObject)
			}
		}
	}
	return n
}

// Payout Scheduler (processes payout if it's been x amount since last payout or if max threshold of Activity Objects reached)
// Threshold = 0 means only time based payouts.
func RunPayoutScheduler(ctx context.Context, app *App, every time.Duration, threshold int) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	var lastPayout time.Time
	var running atomic.Bool
	for { select {
	case <- ctx.Done():
		return
	case <- ticker.C:
		// TODO: time since + threshold reached (so low-amount payouts aren't done), also something wrong here.
		shouldRun := time.Since(lastPayout) >= every || (threshold > 0 && app.Payout.GetPendingActivityObjectsInBuffer() >= threshold)
		if !shouldRun || running.Load() { continue }
		lastPayout = time.Now()
		go func() {
			running.Store(true)
			defer running.Store(false)
			err := ProcessPayout(ctx, app)
			if err != nil { slog.Error("Failed to process payout", "err", err) }
		}()
	}}
}

func ProcessPayout(ctx context.Context, app *App) error {
	queue := app.Payout.Drain()
	if len(queue) == 0 { return ErrEmptyPayout() }

	// begin db transaction
	dbtx, err := app.DB.BeginTx(ctx, nil)
	if err != nil { return ErrBeginTx(err) }
	qtx := app.Queries.WithTx(dbtx)
	for platform, activityByType := range queue {
		for activityType, activityByUser := range activityByType {

			// fetch stuff
			reason, err := qtx.GetTransactionReasonID(ctx, string(ActivityTypeToGenesisTransactionReasons[activityType]))
			// TODO: handle error
			if err != nil { dbtx.Rollback(); return fmt.Errorf("GetTransactionReasonID") }
			// TODO: handle error
			platformTypeID, err := qtx.GetPlatformTypeID(ctx, string(platform))
			if err != nil { dbtx.Rollback(); return fmt.Errorf("GetPlatformTypeID") }
			// create new transaction, get txid
			txidResult, err := qtx.NewTransaction(ctx, db.NewTransactionParams{
				Reason: reason,
				Platform: platformTypeID,
				InitiatedBy: ADMIN_UNID,
			}); if err != nil { dbtx.Rollback(); return ErrNewTransaction(err) }
			txid, err := txidResult.LastInsertId(); if err != nil { dbtx.Rollback(); return ErrLastInsertId(err) }

			// evaluate payout amount per user
			for platformID, activityObjects := range activityByUser {
				var coinValue int64
				var messageCount int32

				// evaluate amount from activityObject
				// NOTE: NEW_ACTIVITY
				switch activityType {
				case ActivityTypeMessage:
					// TODO: implement using timestamp aka time-gap
					for _, object := range activityObjects {
						msg := object.(ActivityObjectMessage)
						messageCount += int32(msg.MessageLength)
					}
					coinValue = int64(messageCount)
				}; if coinValue <= 0 { continue }

				// get unid
				unid, err := app.GetUnidByPlatformTx(ctx, qtx, platform, platformID)
				if err != nil { dbtx.Rollback(); return ErrGetUnidByPlatform(err) }
				// create coin
				coinResult, err := qtx.NewCoin(ctx, db.NewCoinParams{ OwnerUnid: unid, Value: coinValue })
				if err != nil { dbtx.Rollback(); return ErrNewCoin(err) }
				coinId, err := coinResult.LastInsertId()
				if err != nil { dbtx.Rollback(); return ErrLastInsertId(err) }
				// create transaction_link
				txlinkResult, err := qtx.NewTransactionLink(ctx, db.NewTransactionLinkParams{
					Txid: txid,
					CoinID: coinId,
					LinkType: TransactionLinkTypeOutput,
				})
				if err != nil { dbtx.Rollback(); return ErrNewTransactionLink(err) }
				txlink, err := txlinkResult.LastInsertId()
				if err != nil { dbtx.Rollback(); return ErrLastInsertId(err) }
				// create payout row
				err = qtx.NewPayout(ctx, db.NewPayoutParams{ Txlink: txlink, MessageCount: messageCount })
				if err != nil { dbtx.Rollback(); return ErrNewPayout(err) }
			}
		}
	}
	// end db transaction
	if err := dbtx.Commit(); err != nil { dbtx.Rollback(); return ErrCommitTx(err) }
	return nil
}
