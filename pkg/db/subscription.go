package db

import (
	"database/sql"
	"time"
)

type Subscription struct {
	ID          string       `db:"id"`
	ServiceName string       `db:"service_name"`
	Price       int64        `db:"price"`
	UserID      string       `db:"user_id"`
	StartDate   time.Time    `db:"start_date"`
	EndDate     sql.NullTime `db:"end_date"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
