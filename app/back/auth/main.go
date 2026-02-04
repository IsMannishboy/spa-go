package main

import (
	a "auth/internal/auth"
	c "auth/internal/config"
	DB "auth/internal/db"
	h "auth/internal/handlers"
	R "auth/internal/redis"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// config
	Server := c.GetConf()
	// db + redis
	db := DB.GetPostgresConn(Server)
	cash := R.GetRedisConn(Server)
	// auth structs
	var CSRF = new(a.CSRF)
	var SESSIONS = a.NewStruct(cash.Rdb)

	http.HandleFunc("/auth/login", h.LoginHandler(CSRF, db, cash, SESSIONS))
	http.HandleFunc("/auth/register", h.RegisterHandler(CSRF, db, cash, SESSIONS))
	http.ListenAndServe("8080", nil)
}
