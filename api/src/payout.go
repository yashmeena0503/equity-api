package src

import (
	"context"
	"sync"
	"time"
)

// NOTE: What are batches ?
// batches are the data we store to later evaluate and reward users with coins according to the data.
// examples: messages, playtime, achievements.

// TODO: make BatchType type

// platformName > batchType > platformID > batches
type PayoutQueue map[string]map[string]map[string][]MessageBatch

type PayoutBuffer struct {
	mu sync.Mutex
	queue PayoutQueue
}

func (b *PayoutBuffer) Enqueue(platformName, batchType, platformID string, batch MessageBatch) {

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.queue[platformName] == nil {
		b.queue[platformName] = make(map[string]map[string][]MessageBatch)
	}
	if b.queue[platformName][batchType] == nil {
		b.queue[platformName][batchType] = make(map[string][]MessageBatch)
	}
	slice := b.queue[platformName][batchType][platformID]
	b.queue[platformName][batchType][platformID] = append(slice, batch)

}

func (b *PayoutBuffer) Drain() PayoutQueue {
	b.mu.Lock()
	defer b.mu.Unlock()
	snapshot := b.queue
	b.queue = make(PayoutQueue)
	return snapshot
}

// Payout cycle stuff.

type PayoutOptions struct {
	// Minimum Pending Batches.
	MinPending int // 0 = disable threshold
}

func RunPayoutLoop(ctx context.Context, app *App, every time.Duration, opt PayoutOptions) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		select {
			case <- ctx.Done(): return
			case <- ticker.C: runPayoutIfNeeded(app, opt)
		}
	}
}

func (b *PayoutBuffer) PendingCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for _, byType := range b.queue {
		for _, byUser := range byType {
			for _, batches := range byUser {
				n += len(batches)
			}
		}
	}
	return n
}

func runPayoutIfNeeded(app *App, opt PayoutOptions) {
	if opt.MinPending > 0 && app.Payout.PendingCount() < opt.MinPending { return }
	snapshot := app.Payout.Drain()
	if len(snapshot) == 0 { return }
	// TODO: evaluateAll(snapshot)
}