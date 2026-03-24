package models

type AddUser struct {
	Email    string
	Password string
	Username string
}

type User struct {
	ID int
	Xp int
	AddUser
}
