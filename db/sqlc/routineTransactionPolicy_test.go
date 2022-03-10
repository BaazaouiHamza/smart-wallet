package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"git.digitus.me/pfe/smart-wallet/util"
	"github.com/stretchr/testify/require"
)

func createRandomRoutineTransactionPolicy(t *testing.T) RoutineTransactionPolicy {
	arg := CreateRTPParams{
		Name:              util.RandomName(),
		Description:       util.RandomPostDescriptionOrText(),
		NymID:             util.RandomNym(11),
		ScheduleStartDate: time.Now(),
		ScheduleEndDate:   time.Now().AddDate(1, 1, 1),
		Frequency:         "Monthly",
		Amount: json.RawMessage{
			12: 122,
		},
	}

	routineTransactionPolicy, err := testQueries.CreateRTP(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, routineTransactionPolicy)

	return routineTransactionPolicy
}
func TestCreateRoutineTransactionPolicy(t *testing.T) {
	createRandomRoutineTransactionPolicy(t)
}
func TestUpdateAccount(t *testing.T) {
	routineTransactionPolicy := createRandomRoutineTransactionPolicy(t)

	arg := UpdateRTPParams{
		ID:          routineTransactionPolicy.ID,
		Name:        "Modified",
		Description: "Hello im modified",
	}

	routineTransactionPolicy2, err := testQueries.UpdateRTP(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, routineTransactionPolicy2)

}
