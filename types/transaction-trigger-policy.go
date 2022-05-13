package types

import (
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type TransactionTriggerPolicy struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	NymID           identity.PublicKey `json:"nymID"`
	TargetedBalance ptclTypes.Balance  `json:"targetedBalance"`
	Recipient       identity.PublicKey `json:"recipient"`
	Amount          ptclTypes.Balance  `json:"amount"`
}
