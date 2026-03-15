package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	structs "main/internal/structs"
	"net/http"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type SESSIONS struct {
	rdb *redis.Client
}

func (s *SESSIONS) Checkout(r *http.Request) (structs.Session, error) {
	var Session structs.Session
	cookie_session, err := r.Cookie("session")
	if err != nil {
		return Session, err
	}
	MainSessionId := cookie_session.Value
	MainSession, err := s.FindMainSession(MainSessionId)
	if err != nil {
		return Session, err
	}
	user_session, err := s.GetUserSessionId(r.Header.Get("user_id"), MainSession)
	if err != nil {
		return Session, err
	}
	session_value, err := s.CheckSession(5, user_session)
	if err != nil {
		return Session, err
	}
	Session.Id = user_session
	Session.Value = session_value
	return Session, nil
}
func (s *SESSIONS) FindMainSession(session_id string) (structs.MainSession, error) {
	var MainSessionValue structs.MainSessionValue
	var MainSession structs.MainSession
	MainSession.Id = session_id
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
	fmt.Println("main session value:", MainSession.Value)
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
	MainSession.Id = MainSessionIdStr
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
	if s.rdb == nil {
		panic("rdb is nil")
	}

	if err != nil {
		return MainSession, err
	}
	return MainSession, nil
}
func (s *SESSIONS) CreateSession(ctx context.Context, user_id string) (string, error) {
	var session structs.Session
	var newctx context.Context = context.Background()
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
func (s *SESSIONS) GetUserSessionId(user_id string, MainSession structs.MainSession) (string, error) {
	user_session, ok := MainSession.Value.UserSessions[user_id]
	if !ok {
		err := errors.New("session is absent")
		return "", err
	}
	return user_session, nil

}

func NewStruct(rdb *redis.Client) *SESSIONS {
	return &SESSIONS{rdb: rdb}
}
func (s *SESSIONS) CheckSession(timeout int, session_id string) (structs.SessionValue, error) {
	var ctx context.Context = context.Background()
	var c context.CancelFunc
	var SessionValue structs.SessionValue

	if timeout > 0 {
		ctx, c = context.WithTimeout(context.Background(), time.Duration(timeout*int(time.Second)))
		defer c()
	}
	stored, err := s.rdb.Get(ctx, session_id).Result()
	if err != nil {
		return SessionValue, err
	}
	err = json.Unmarshal([]byte(stored), &SessionValue)
	if err != nil {
		return SessionValue, err
	}
	if time.Now().After(SessionValue.Exp) {
		return SessionValue, errors.New("session expired")
	}
	return SessionValue, nil
}

func (s *SESSIONS) DeleteSession(session_id string) {

}
