package src

import (
	"time"
)

// What is ActivityType ?
// activities are ways to earn equity or other stuff.
// there are ActivityType, Examples: Messages, Playtime, Achievements.
// then there are ActivityObjects, which are the instances of the activity held onto the PayoutQueue
// to be evaluated and rewarded to the activity initiator later.

type ActivityType string
const (
	// NOTE: NEW_ACTIVITY
	ActivityTypeMessage ActivityType = "message"
)
var ActivityTypes = []ActivityType {
	// NOTE: NEW_ACTIVITY
	ActivityTypeMessage,
}

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

// What are Transaction Reasons ?
// string definition of why a transaction happened
// genesis transactions are ones which generate new coins

type TransactionReason string
const (
	// NOTE: NEW_ACTIVITY
	TransactionReasonGenesisActivityMessage TransactionReason = "GenesisActivityMessage"
	// NOTE: NEW_TRANSACTION_REASON
)
var TransactionReasons = []TransactionReason {
	// NOTE: NEW_ACTIVITY
	TransactionReasonGenesisActivityMessage,
	// NOTE: NEW_TRANSACTION_REASON
}

// What are PlatformTypes ?
// platforms that are supports

type PlatformType string
const (
	// NOTE: NEW_PLATFORM
	PlatformDiscord PlatformType = "discord"
	PlatformMinecraft PlatformType = "minecraft"
)
var PlatformTypes = []PlatformType {
	// NOTE: NEW_PLATFORM
	PlatformDiscord,
	PlatformMinecraft,
}

// What are TransactionLinkType ?
// input or output, whether a coin is a input or an output to a transaction
// imagine x amount of people sitting around a table and either
// inputting (putting coins) or outputting (taking coins) onto and from the table respectively

const (
	TransactionLinkTypeInput string = "input"
	TransactionLinkTypeOutput string = "output"
)

// Link tables

var ActivityTypeToGenesisTransactionReasons = map[ActivityType]TransactionReason {
	// NOTE: NEW_ACTIVITY
	// NOTE: NEW_TRANSACTION_REASON
	ActivityTypeMessage: TransactionReasonGenesisActivityMessage,
}
