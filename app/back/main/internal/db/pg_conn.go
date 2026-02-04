package db

import (
	"database/sql"
)

type Storage struct {
	Db      *sql.DB
	Timeout int
}

func GetStorage() {

}
