package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

type CSRF struct {
}

func RandomKey() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
func (c *CSRF) MakeTokenAndKey() (string, string, error) {
	key, err := RandomKey()
	if err != nil {

		return "", "", err
	}
	rawkey, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {

		return "", "", err
	}
	salt := make([]byte, 10)
	if _, err := rand.Read(salt); err != nil {

		return "", "", err
	}
	mac := hmac.New(sha256.New, rawkey)
	mac.Write(salt)
	token := mac.Sum(nil)
	resp := append(token, salt...)
	return base64.RawURLEncoding.EncodeToString(resp), key, nil

}
func (c *CSRF) CheckToken(token string, key string) error {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return err
	}
	raw_salt := raw[len(raw)-10:]
	raw_token := raw[:len(raw)-10]
	rawkey, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, rawkey)
	mac.Write(raw_salt)
	new := mac.Sum(nil)
	if !hmac.Equal(new, raw_token) {
		return errors.New("invalid CSRF token")
	}
	return nil
}
func (c *CSRF) PostRequest(w http.ResponseWriter, r *http.Request) int {
	csrf := r.Header.Get("CSRF")
	fmt.Println("token check:", csrf)
	cookie, err := r.Cookie("KEY")
	if err != nil {
		fmt.Println("cookie error:", err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return 1
	}
	key := cookie.Value
	CSRFError := c.CheckToken(csrf, key)
	if CSRFError != nil {
		fmt.Println("csrf error:", CSRFError)
		http.Error(w, CSRFError.Error(), http.StatusForbidden)
		return 1
	}
	return 0
}
