package models

import (
	"database/sql"
	"time"
)

type RoutineTransactionPolicy struct {
	ID                sql.NullInt64 `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Sender            string        `json:"sender"`
	Receiver          string        `json:"receiver"`
	CreatedAt         time.Time     `json:"created_at"`
	ScheduleStartDate time.Time     `json:"schedule_start_date"`
	ScheduleEndDate   time.Time     `json:"schedule_end_date"`
	Frequency         string        `json:"frequency"`
	Amount            int32         `json:"amount"`
}
