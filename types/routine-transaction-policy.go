package types

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type RoutineTransactionPolicy struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	NymID             identity.PublicKey `json:"nymID"`
	Recipient         identity.PublicKey `json:"recipient"`
	ScheduleStartDate time.Time          `json:"scheduleStartDate"`
	ScheduleEndDate   time.Time          `json:"scheduleEndDate"`
	Frequency         string             `json:"frequency"`
	Amount            ptclTypes.Balance  `json:"amount"`
	RequestType       string             `json:"requestType"`
}
