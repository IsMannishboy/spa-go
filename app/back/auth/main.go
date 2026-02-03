package auth

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
	rdb := R.GetRedisConn(Server)
	// auth structs
	var CSRF = new(a.CSRF)
	http.HandleFunc("/auth/login", h.LoginHandler(CSRF, db, rdb))
}
