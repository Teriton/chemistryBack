// Package models contains models for package
package models

import "time"

type AddUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type User struct {
	ID           int       `json:"id"`
	Xp           int       `json:"xp"`
	Streak       int       `json:"streak"`
	CreationDate time.Time `json:"creation_date"`
	AddUser
}
