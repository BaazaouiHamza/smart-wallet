package types

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
)

type RoutineTransactionPolicy struct {
	Name              string
	Description       string
	NymID             identity.PublicKey
	ScheduleStartDate time.Time
	ScheduleEndDate   time.Time
	Frequency         string
	Amount            int64
}