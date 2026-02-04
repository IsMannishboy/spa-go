package auth

import (
	structs "auth/internal/structs"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type SESSIONS struct {
	rdb *redis.Client
}

func (s *SESSIONS) CreateSession(ctx context.Context, user_id string) (string, error) {
	var session structs.Session
	var newctx context.Context = nil
	var c context.CancelFunc
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return session.Id, err
	}
	session.Id = base64.RawURLEncoding.EncodeToString(b)
	session.Value.Exp = time.Now().Add(time.Hour * 24)
	session.Value.UserId = user_id
	if ctx != nil {
		newctx, c = context.WithTimeout(ctx, s.rdb.Options().DialTimeout*time.Second)
		defer c()
	}
	byte_value, err := json.Marshal(session.Value)
	if err != nil {
		return session.Id, err
	}
	str_value := base64.RawURLEncoding.EncodeToString(byte_value)
	_, err = s.rdb.Set(newctx, session.Id, str_value, time.Minute*30).Result()
	return session.Id, err
}
func NewStruct(rdb *redis.Client) *SESSIONS {
	return &SESSIONS{rdb: rdb}
}
func (s *SESSIONS) CheckSession(timeout int, session_id string) error {
	var ctx context.Context = nil
	var c context.CancelFunc
	if timeout > 0 {
		ctx, c = context.WithTimeout(context.Background(), time.Duration(timeout*int(time.Second)))
		defer c()
	}
	stored, err := s.rdb.Get(ctx, session_id).Result()
	if err != nil {
		return err
	}
	var SessionValue structs.SessionValue
	err = json.Unmarshal([]byte(stored), &SessionValue)
	if err != nil {
		return err
	}
	if time.Now().After(SessionValue.Exp) {
		return errors.New("session expired")
	}
	return nil
}
func (s *SESSIONS) DeleteSession(session_id string) {

}
