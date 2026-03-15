package structs

import "time"

type User struct {
	Id        string
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}
type SessionValue struct {
	UserId string
	Email  time.Time
	Exp    time.Time
}
type MainSession struct {
	Id    string
	Value MainSessionValue
}
type MainSessionValue struct {
	UserSessions map[string]string
	Exp          time.Time
}
type Session struct {
	Id    string
	Value SessionValue
}
type Login struct {
	Username string
	Password string
}
type Register struct {
	Username string
	Password string
	Email    string
}
