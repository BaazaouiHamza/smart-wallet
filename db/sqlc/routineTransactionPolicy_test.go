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
		Sender:            util.RandomNym(11),
		Receiver:          util.RandomNym(11),
		CreatedAt:         time.Now(),
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
