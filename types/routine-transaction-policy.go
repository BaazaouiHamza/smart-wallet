package types

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type RoutineTransactionPolicy struct {
	ID                int
	Name              string
	Description       string
	NymID             identity.PublicKey
	Recipient         identity.PublicKey
	ScheduleStartDate time.Time
	ScheduleEndDate   time.Time
	Frequency         string
	Amount            ptclTypes.Balance
}
