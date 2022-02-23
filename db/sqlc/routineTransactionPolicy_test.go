package db

import (
	"context"
	"testing"
	"time"

	"git.digitus.me/pfe/smart-wallet/util"
	"github.com/stretchr/testify/require"
)

func createRandomRoutineTransactionPolicy(t *testing.T) RoutineTransactionPolicy {
	arg := CreateRoutineTransactionPolicyParams{
		Name:              util.RandomName(),
		Description:       util.RandomPostDescriptionOrText(),
		NymID:             util.RandomNym(11),
		ScheduleStartDate: time.Now(),
		ScheduleEndDate:   time.Now().AddDate(1, 1, 1),
		Frequency:         "Monthly",
		Amount:            10,
	}

	routineTransactionPolicy, err := testQueries.CreateRoutineTransactionPolicy(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, routineTransactionPolicy)

	return routineTransactionPolicy
}
func TestCreateRoutineTransactionPolicy(t *testing.T) {
	createRandomRoutineTransactionPolicy(t)
}
func TestUpdateAccount(t *testing.T) {
	routineTransactionPolicy := createRandomRoutineTransactionPolicy(t)

	arg := UpdateRoutineTransactionPolicyParams{
		ID:          routineTransactionPolicy.ID,
		Name:        "Modified",
		Description: "Hello im modified",
	}

	routineTransactionPolicy2, err := testQueries.UpdateRoutineTransactionPolicy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, routineTransactionPolicy2)

}
