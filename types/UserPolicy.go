package types

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
)

type UserPolicy struct {
	ID          int
	Name        string
	Description string
	NymID       identity.PublicKey
	CreatedAt   time.Time
}
