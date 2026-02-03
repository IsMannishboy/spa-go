package models

import (
	"context"
	"database/sql"
	"time"
)

type Model interface {
	FindOne(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration) (interface{}, error)
}
