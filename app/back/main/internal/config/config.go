package config

import (
	"github.com/joho/godotenv"
)

type Server struct {
}

func GetConf() {
	path := "/home/savage21/hardWork/go/go-spa/app/back/main/env.env"
	godotenv.Load(path)
}
