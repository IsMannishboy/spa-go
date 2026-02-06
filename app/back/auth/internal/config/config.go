package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Server struct {
	Addr     string
	Port     string
	Postgres Postgres
	Redis    Redis
}
type Postgres struct {
	Host        string
	Port        string
	DialTimeout int
	Db          string
	User        string
	Password    string
	SSLMode     string
}
type Redis struct {
	Addr        string
	DialTimeout int
	Password    string
}

func GetConf() *Server {
	path := "/home/savage21/hardWork/go/go-spa/app/back/auth/env.env"
	err := godotenv.Load(path)
	fmt.Println("ENVPATH =", path)
	if err != nil {
		panic(err)
	}
	var Server Server
	Server.Addr = os.Getenv("addr")
	Server.Postgres.Host = os.Getenv("pg_host")
	Server.Postgres.Port = os.Getenv("pg_port")
	Server.Postgres.User = os.Getenv("pg_user")
	Server.Postgres.Db = os.Getenv("pg_db")
	Server.Postgres.Password = os.Getenv("pg_password")
	Server.Postgres.SSLMode = os.Getenv("pg_ssl")
	raw_pd_timeout := os.Getenv("pg_timeout")
	Server.Postgres.DialTimeout, err = strconv.Atoi(raw_pd_timeout)
	if err != nil {
		panic(err)
	}

	Server.Redis.Addr = os.Getenv("redis_addr")

	redis_raw_timeout := os.Getenv("redis_timeout")
	Server.Redis.DialTimeout, err = strconv.Atoi(redis_raw_timeout)
	if err != nil {
		panic(err)
	}
	return &Server
}
