package authmngr

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	rawPassword := "654321"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	check(err, t)
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(rawPassword))
	check(err, t)
	fmt.Printf("Raw password: %v\n HashedPassword: %s\n", rawPassword, hashedPassword)
}

func TestPasswordHasher(t *testing.T) {
	pswHasher, err := NewPasswordHasher()
	check(err, t)
	password := "654321"
	hashedPassword, err := pswHasher.encode([]byte(password))
	check(err, t)
	if !pswHasher.check(hashedPassword, []byte(password)) {
		t.Error("Passwords checking is not passed")
	}
}
