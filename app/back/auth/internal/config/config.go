package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Server struct {
	Host     string
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
	Host        string
	Port        string
	DialTimeout int
	Password    string
	Db          string
}

func GetConf() *Server {
	path := "../env.env"
	err := godotenv.Load(os.Getenv(path))
	fmt.Println("ENVPATH =", os.Getenv(path))
	if err != nil {
		panic(err)
	}
	var Server Server
	Server.Host = os.Getenv("host")
	Server.Port = os.Getenv("port")
	Server.Postgres.Host = os.Getenv("pg_host")
	Server.Postgres.Port = os.Getenv("pg_port")
	Server.Postgres.Password = os.Getenv("pg_password")
	Server.Postgres.SSLMode = os.Getenv("pd_ssl")
	raw_pd_timeout := os.Getenv("pd_timeout")
	Server.Postgres.DialTimeout, err = strconv.Atoi(raw_pd_timeout)
	if err != nil {
		panic(err)
	}
	Server.Redis.Db = os.Getenv("redis_db")
	Server.Redis.Host = os.Getenv("redis_host")
	Server.Redis.Port = os.Getenv("redis_port")
	Server.Redis.Password = os.Getenv("redis_password")
	redis_raw_timeout := os.Getenv("redis_timeout")
	Server.Redis.DialTimeout, err = strconv.Atoi(redis_raw_timeout)
	if err != nil {
		panic(err)
	}
	return &Server
}
