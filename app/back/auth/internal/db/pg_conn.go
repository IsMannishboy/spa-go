package db

import (
	c "auth/internal/config"
	"time"

	m "auth/internal/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	db      *sql.DB
	Timeout time.Duration
	Models  map[string]*m.Model
}

func GetPostgresConn(conf *c.Server) *DB {

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Db,
		conf.Postgres.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	var PingError error
	for i := 0; i < 10; i++ {
		PingError = db.Ping()
		if PingError == nil {
			break
		}
	}
	if PingError != nil {
		panic(PingError)
	}
	Models := make(map[string]m.Model)
	var NewUsersModel = &m.UserModel{}
	Models["users"] = NewUsersModel
	return &DB{db, time.Duration(conf.Postgres.DialTimeout), Models}
}
