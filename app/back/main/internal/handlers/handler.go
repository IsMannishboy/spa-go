package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	auth "main/internal/auth"
	r "main/internal/redis"
	"math/big"
	"net/http"
	"time"
)

func GetCodeVerivication(SessionValidator *auth.SESSIONS, cash *r.Cash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_value, err := SessionValidator.Checkout(r)
		if err != nil {
			http.Error(w, err.Error(), 403)
			return
		}
		if time.Now().Before(session_value.Email) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("email is already validated"))
			return
		}
		code, err := MakeCode()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, err = cash.Rdb.Set(context.Background(), r.Header.Get("user_id"), code, time.Duration(time.Minute*1)).Result()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// send code to email
		w.WriteHeader(http.StatusCreated)
		w.Write(nil)
	}
}
func PostCodeVerification(SessionValidator *auth.SESSIONS, cash *r.Cash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserSession, err := SessionValidator.Checkout(r)
		if err != nil {
			http.Error(w, err.Error(), 403)
			return
		}
		if time.Now().Before(UserSession.Value.Email) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("email is already validated"))
			return
		}
		code := r.Header.Get("code")

		code_from_redis, err := cash.Rdb.Get(context.Background(), r.Header.Get("user_id")).Result()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if code != code_from_redis {
			w.Header().Set("validated", "false")
			http.Error(w, "code doesnt match", 403)
			return
		}
		_, err = cash.Rdb.Del(context.Background(), r.Header.Get("user_id")).Result()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		UserSession.Value.Email = time.Now().Add(time.Duration(time.Minute * 30))
		_, err = cash.Rdb.Set(context.Background(), UserSession.Id, UserSession.Value, time.Duration(time.Minute*30)).Result()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("validated", "true")
		w.Write(nil)
	}
}
func GetEmailVerivication(SessionValidator *auth.SESSIONS, cash *r.Cash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_value, CheckoutError := SessionValidator.Checkout(r)
		if CheckoutError != nil {
			http.Error(w, CheckoutError.Error(), 403)
			return
		}
		if time.Now().Before(session_value.Email) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("email is already validated"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(nil)

	}
}
func MakeCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
