// Package models contains models for package
package models

type AddUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type User struct {
	ID int `json:"id"`
	Xp int `json:"xp"`
	AddUser
}
