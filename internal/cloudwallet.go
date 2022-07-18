package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

var _ interface {
	SendAmounts(
		ctx context.Context,
		sender identity.PublicKey,
		recipients map[identity.PublicKey]ptclTypes.Balance,
		context string,
	) (*ptclTypes.RawTransfer, error)
	GetUserState(
		context.Context, identity.PublicKey,
	) (*ptclTypes.UserState, error)
} = (*CloudwalletClient)(nil)

type (
	CloudwalletClient struct {
		*http.Client
	}

	sendRequest struct {
		Encoding   bool                                     `json:"std"`
		Context    string                                   `json:"context"`
		Recipients map[identity.PublicKey]ptclTypes.Balance `json:"recipients"`
	}
	cwTransfer struct {
		Tranfer *ptclTypes.RawTransfer `json:"transfer"`
	}
	coreUserState struct {
		Nym               identity.PublicKey          `json:"nym"`
		Balance           ptclTypes.Balance           `json:"balance"`
		TransferSequences ptclTypes.TransferSequences `json:"transferSequences"`
	}
	cwUserState struct {
		UserState *coreUserState
	}
)

func (c *CloudwalletClient) SendAmounts(
	ctx context.Context,
	sender identity.PublicKey,
	recipients map[identity.PublicKey]ptclTypes.Balance,
	context string,
) (*ptclTypes.RawTransfer, error) {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(sendRequest{
		Encoding:   true,
		Context:    context,
		Recipients: recipients,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("http://cloudwallet/wallet1.0/wallets/%s/transactions/send", sender),
		buf,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("could not send transaction")
	}
	defer res.Body.Close()
	var t cwTransfer
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, err
	}

	return t.Tranfer, nil
}

func (c *CloudwalletClient) GetUserState(
	ctx context.Context, pk identity.PublicKey,
) (*ptclTypes.UserState, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("http://cloudwallet/wallet1.0/wallets/%s", pk),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
	}

	defer res.Body.Close()
	var us cwUserState

	if err := json.NewDecoder(res.Body).Decode(&us); err != nil {
		return nil, err
	}

	return &ptclTypes.UserState{
		UserNymID:         us.UserState.Nym,
		Balance:           us.UserState.Balance,
		TransferSequences: us.UserState.TransferSequences,
		Signature:         nil,
	}, nil
}
