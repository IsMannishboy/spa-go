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
	"time"

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

		session_cookie, err := r.Cookie("session")
		session_id, err := MainSessionFunc(err, SESSIONS, UserOk.Id, session_cookie)
		if err != nil {
			fmt.Println("MainSessionFunc error:", err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println("New Main session :", session_id)
		cookie := &http.Cookie{
			Name:     "session",
			Value:    session_id,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			MaxAge:   3600 * 24,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("user_id", UserOk.Id)
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
		fmt.Println("register:", Register)
		UserModel := db.Models["users"].(*m.UserModel)
		_, err = UserModel.FindOne(context.Background(), map[string]string{"username": Register.Username}, db.DB, time.Duration(db.Timeout))
		if err != sql.ErrNoRows && err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err == nil {
			http.Error(w, "this username already is used", 403)
			return
		}
		user_id, err := UserModel.Create(context.Background(), int(db.Timeout), Register, db.DB)
		if err != nil {
			fmt.Println("sreate user error:", err)
			http.Error(w, err.Error(), 500)
			return
		}
		session_cookie, err := r.Cookie("session")
		session_id, err := MainSessionFunc(err, SESSIONS, user_id, session_cookie)
		if err != nil {
			fmt.Println("MainSessionFunc error:", err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println("New Main session :", session_id)

		cookie := &http.Cookie{
			Name:     "session",
			Value:    session_id,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			MaxAge:   3600 * 24,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("user_id", user_id)
		w.Write([]byte("session set"))
	}
}
func MainSessionFunc(err error, SESSIONS *a.SESSIONS, user_id string, session_cookie *http.Cookie) (string, error) {
	var session_id string
	switch err {
	case nil:

		session_id = session_cookie.Value
		fmt.Println("session from cookie:", session_id)
		MainSession, err := SESSIONS.FindMainSession(session_id)
		if err != nil {
			fmt.Println("find session error", err)

			return session_id, err
		}
		fmt.Println("FindMainSession MainSession.ID:", MainSession.Id)
		user_session_id, err := SESSIONS.CreateSession(nil, user_id)
		if err != nil {
			fmt.Println("create user session error", err)

			return session_id, err
		}
		err = SESSIONS.UpdateMainSession(&MainSession, user_id, user_session_id)
		fmt.Println("Main SESSION VALUE:", MainSession.Value)
		return session_id, nil

	default:
		fmt.Println("session is absent(")
		MainSession, err := SESSIONS.CreateMainSession(user_id)

		if err != nil {
			fmt.Println("create CreateMainSession session error", err)

			return session_id, err
		}

		fmt.Println(" CreateMainSession MainSession.ID:", MainSession.Id)
		user_session, err := SESSIONS.CreateSession(nil, user_id)
		if err != nil {
			fmt.Println("create  session error", err)

			return session_id, err
		}
		err = SESSIONS.UpdateMainSession(&MainSession, user_id, user_session)
		if err != nil {
			fmt.Println("create CreateMainSession session error", err)

			return session_id, err
		}
		fmt.Println("Main SESSION VALUE:", MainSession.Value)
		return MainSession.Id, nil

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
		token_cookie := &http.Cookie{
			Name:     "CSRF",
			Value:    key,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		http.SetCookie(w, token_cookie)
		w.Header().Set("CSRF", token)
		w.Write([]byte("take your token"))
	}
}
