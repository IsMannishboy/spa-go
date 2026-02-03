package handlers

import (
	a "auth/internal/auth"
	d "auth/internal/db"
	s "auth/internal/structs"
	"encoding/json"
	"io"
	"net/http"

	redis "github.com/redis/go-redis/v9"
)

func LoginHandler(CSRF *a.CSRF, db *d.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CSRF.PostRequest(w, r) == 1 {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		var LoginStruct s.Login
		err = json.Unmarshal(data, &LoginStruct)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	}

}
