package models

import (
	s "auth/internal/structs"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserModel struct {
}

func (u *UserModel) FindOne(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration) (interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var conditions []string
	for key, value := range params {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, value))
	}
	var args = strings.Join(conditions, "AND")
	query := "select * from users where " + args
	var User s.User
	newctx, c := context.WithTimeout(ctx, timeout)
	defer c()
	err := db.QueryRowContext(newctx, query).Scan(&User)
	if err != nil {
		if err == sql.ErrNoRows {
			return User, err
		}
		return User, err
	}
	return User, nil
}
