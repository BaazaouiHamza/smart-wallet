package service

import (
	"encoding/json"

	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

func marshalBoth(a, b ptclTypes.Balance) ([]byte, []byte, error) {
	x, err := a.MarshalJSON()
	if err != nil {
		return nil, nil, err
	}

	y, err := b.MarshalJSON()

	return x, y, err
}

func unmarshalBoth(a, b []byte) (ptclTypes.Balance, ptclTypes.Balance, error) {
	var x, y ptclTypes.Balance

	if err := json.Unmarshal(a, &x); err != nil {
		return nil, nil, err
	}

	return x, y, json.Unmarshal(b, &y)
}
