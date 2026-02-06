package models

import (
	s "auth/internal/structs"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
}

func (u *UserModel) FindOne(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration) (any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var conditions []string
	for key, value := range params {
		conditions = append(conditions, fmt.Sprintf("%s = '%s'", key, value))
	}
	var args = strings.Join(conditions, "AND")
	query := "select * from users where " + args
	fmt.Println("query:", query)
	var User s.User
	newctx, c := context.WithTimeout(ctx, timeout*time.Second)
	defer c()
	err := db.QueryRowContext(newctx, query).Scan(&User)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error() == sql.ErrNoRows.Error())
		return User, err
	}
	return User, nil
}
func (u *UserModel) Find(ctx context.Context, params map[string]string, db *sql.DB, timeout time.Duration, max int64) (any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var conditions []string
	for key, value := range params {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, value))
	}
	var args = strings.Join(conditions, "AND")
	query := "select * from users where " + args
	var Users []s.User
	newctx, c := context.WithTimeout(ctx, timeout)
	defer c()
	rows, err := db.QueryContext(newctx, query)
	if err != nil {

		return Users, err
	}
	var i int64 = 0
	for rows.Next() {
		if max != 0 {
			if i > max {
				break
			}
		}
		var user s.User
		err = rows.Scan(&user)
		if err != nil {
			fmt.Println("find users error:", err)
			return Users, err
		}
		Users = append(Users, user)
		i++
	}
	return Users, nil
}
func (u *UserModel) Create(ctx context.Context, timeout int, user s.Register, db *sql.DB) (string, error) {
	var newctx context.Context
	var c context.CancelFunc
	var id string
	if ctx != nil {
		newctx, c = context.WithTimeout(ctx, time.Duration(timeout))
		defer c()
	} else {
		newctx = nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return id, err
	}
	stored_hash := base64.RawURLEncoding.EncodeToString(hash)
	_, err = db.ExecContext(newctx, "insert into users (username,email,password,created_at) values ($1,$2,$3,$4);",
		user.Username, user.Email, stored_hash, time.Now())
	var int_id int
	err = db.QueryRowContext(newctx, "select id from users where username = $1", user.Username).Scan(&int_id)
	if err != nil {
		return id, err

	}
	id = strconv.Itoa(int(int_id))
	return id, nil
}
func (u *UserModel) CheckPassword(password string, hash string) error {
	h, err := base64.RawURLEncoding.DecodeString(hash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(h, []byte(password))
	return err
}
