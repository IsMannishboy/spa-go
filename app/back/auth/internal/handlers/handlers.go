package handlers

import (
	a "auth/internal/auth"
	d "auth/internal/db"
	m "auth/internal/models"
	r "auth/internal/redis"
	s "auth/internal/structs"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"encoding/json"
	"io"
	"net/http"
)

func LoginHandler(CSRF *a.CSRF, db *d.DB, cash *r.Cash, SESSIONS *a.SESSIONS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CSRF.PostRequest(w, r) == 1 {
			http.Error(w, "forbidden", 403)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		r.Body.Close()
		var LoginStruct s.Login
		err = json.Unmarshal(data, &LoginStruct)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		UserModel := db.Models["users"].(*m.UserModel)
		user, err := UserModel.FindOne(context.Background(), map[string]string{"username": LoginStruct.Username}, db.DB, 5)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, err.Error(), 404)
				return
			} else if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("context closed")
				http.Error(w, err.Error(), 500)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		UserOk, ok := user.(s.User)
		if !ok {
			panic(ok)
		}
		err = UserModel.CheckPassword(LoginStruct.Password, UserOk.Password)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		session_id, err := SESSIONS.CreateSession(context.Background(), UserOk.Id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session_id,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			MaxAge:   3600 * 24,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("msg", "login sucessfull")
		w.Write([]byte("session set"))

	}

}
func RegisterHandler(CSRF *a.CSRF, db *d.DB, cash *r.Cash, SESSIONS *a.SESSIONS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CSRF.PostRequest(w, r) == 1 {
			http.Error(w, "forbidden", 403)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		var Register s.Register
		err = json.Unmarshal(data, &Register)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		UserModel := db.Models["users"].(*m.UserModel)
		user_id, err := UserModel.Create(context.Background(), int(db.Timeout), Register, db.DB)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		session_id, err := SESSIONS.CreateSession(context.Background(), user_id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session_id,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   3600 * 24,
		}
		http.SetCookie(w, cookie)
		w.Header().Set("msg", "registration sucessfull")
		w.Write([]byte("session set"))

	}
}
func GetCSRF(CSRF *a.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, key, err := CSRF.MakeTokenAndKey()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println("token:", token)
		fmt.Println("key:", key)
		cookie := &http.Cookie{
			Name:     "KEY",
			Value:    key,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		w.Header().Set("CSRF", token)
		w.Write([]byte("take your token"))
	}
}
