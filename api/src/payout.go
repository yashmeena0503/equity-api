package src

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// NOTE: What are activities ?
// activities are ways to earn equity or other stuff.
// there are ActivityType, Examples: Messages, Playtime, Achievements.
// then there are ActivityObjects, which are the instances of the activity held onto the PayoutQueue
// to be evaluated and rewarded to the activity initiator later.

type ActivityType string
const (
	// NOTE: NEW_ACTIVITY
	ActivityTypeMessage ActivityType = "message"
)
type ActivityObject interface {
	Type() ActivityType
	OccuredAt() time.Time
}

// NOTE: NEW_ACTIVITY

type ActivityObjectMessage struct {
	Timestamp time.Time `json:"timestamp"`
	MessageLength int `json:"message_length"`
}
func (a ActivityObjectMessage) Type() ActivityType { return ActivityTypeMessage }
func (a ActivityObjectMessage) OccuredAt() time.Time { return a.Timestamp }

// END OF NEW_ACTIVITY

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
			ProcessPayout(ctx, app)
		}()
	}}
}

func ProcessPayout(ctx context.Context, app *App) {
	queue := app.Payout.Drain()
	if len(queue) == 0 { return }

	for platform, activityByType := range queue {
		for activityType, activityByUser := range activityByType {
			for platformID, activityObject := range activityByUser {

				_ = platform
				_ = activityType
				_ = platformID
				_ = activityObject

			}
		}
	}

}
