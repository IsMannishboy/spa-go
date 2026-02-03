package structs

type User struct {
	Username string
	Password string
	Email    string
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
