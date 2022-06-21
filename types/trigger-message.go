package types

import (
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type TriggerMessage struct {
	PolicyID int                                      `json:"policyID"`
	Amounts  map[identity.PublicKey]ptclTypes.Balance `json:"amounts"`
}
