package src

import (
	"context"
	"sync"
	"time"
)

// NOTE: What are bundles ?
// bundles are the data we store to later evaluate and reward users with coins according to the data.
// examples: messages, playtime, achievements.

// TODO: make BundleType type

// platformName > bundleType > platformID > bundles
type PayoutQueue map[string]map[string]map[string][]MessageBundle

type PayoutBuffer struct {
	mu sync.Mutex
	queue PayoutQueue
}
func (b *PayoutBuffer) Enqueue(platformName, bundleType, platformID string, bundle MessageBundle) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.queue[platformName] == nil {
		b.queue[platformName] = make(map[string]map[string][]MessageBundle)
	}
	if b.queue[platformName][bundleType] == nil {
		b.queue[platformName][bundleType] = make(map[string][]MessageBundle)
	}
	slice := b.queue[platformName][bundleType][platformID]
	b.queue[platformName][bundleType][platformID] = append(slice, bundle)
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
    // deep copy so the handler doesn't hold the lock or expose internal state
    out := make(PayoutQueue)
    for platform, byType := range b.queue {
        out[platform] = make(map[string]map[string][]MessageBundle)
        for bundleType, byUser := range byType {
            out[platform][bundleType] = make(map[string][]MessageBundle)
            for platformID, bundles := range byUser {
                cp := make([]MessageBundle, len(bundles))
                copy(cp, bundles)
                out[platform][bundleType][platformID] = cp
            }
        }
    }
    return out
}

// Payout cycle stuff.

type PayoutOptions struct {
	// Minimum Pending Bundles.
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
			for _, bundles := range byUser {
				n += len(bundles)
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