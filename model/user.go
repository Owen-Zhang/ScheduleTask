package model

type User struct {
	Id        int
	UserName  string
	Password  string
	Salt      string
	Email     string
	LastLogin int64
	LastIp    string
	Status    int
}
