package models

import (
	"context"
	"database/sql"
	"time"
)

type Model interface {
	FindOne(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration) (any, error)
	Find(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration, max int64) (any, error)
}
