package types

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type TransactionTriggerPolicy struct {
	ID              int
	Name            string
	Description     string
	NymID           identity.PublicKey
	CreatedAt       time.Time
	TargetedBalance ptclTypes.Balance
	Recipient       identity.PublicKey
	Amount          ptclTypes.Balance
}
