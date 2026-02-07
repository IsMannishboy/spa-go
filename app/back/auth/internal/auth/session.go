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

func (s *SESSIONS) FindMainSession(session_id string) (structs.MainSession, error) {
	var MainSessionValue structs.MainSessionValue
	var MainSession structs.MainSession
	stored, err := s.rdb.Get(context.Background(), session_id).Result()
	if err != nil {
		return MainSession, err
	}
	bytes, err := base64.RawURLEncoding.DecodeString(stored)
	if err != nil {
		return MainSession, err
	}
	err = json.Unmarshal(bytes, &MainSessionValue)
	if err != nil {
		return MainSession, err
	}
	MainSession.Value = MainSessionValue
	return MainSession, nil
}
func (s *SESSIONS) UpdateMainSession(MainSession *structs.MainSession, user_id string, user_session_id string) error {
	MainSession.Value.UserSessions[user_id] = user_session_id
	bytes, err := json.Marshal(MainSession.Value)
	if err != nil {
		return err
	}
	str_value := base64.RawURLEncoding.EncodeToString(bytes)
	_, err = s.rdb.Set(context.Background(), MainSession.Id, str_value, time.Hour*24).Result()
	if err != nil {
		return err
	}
	return nil
}
func (s *SESSIONS) CreateMainSession(user_id string) (structs.MainSession, error) {
	//id
	var MainSession structs.MainSession
	MainSessionId := make([]byte, 20)
	_, err := rand.Read(MainSessionId)
	if err != nil {
		return MainSession, err
	}
	MainSessionIdStr := base64.RawURLEncoding.EncodeToString(MainSessionId)

	var MainSessionValue structs.MainSessionValue
	MainSession.Value = MainSessionValue
	var UserSessions = make(map[string]string)
	MainSession.Value.UserSessions = UserSessions
	MainSession.Value.Exp = time.Now().Add(time.Hour * 24)
	MainSessionByte, err := json.Marshal(MainSession.Id)
	str_value := base64.RawURLEncoding.EncodeToString(MainSessionByte)
	if err != nil {
		return MainSession, err
	}
	_, err = s.rdb.Set(context.Background(), MainSessionIdStr, str_value, 24*time.Hour).Result()
	if err != nil {
		return MainSession, err
	}
	return MainSession, nil
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
