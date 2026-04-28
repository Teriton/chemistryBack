package authmngr

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	encode([]byte) ([]byte, error)
	check(hashedPassword []byte, rawPassword []byte) bool
}

type PasswordHasher struct {
}

func NewPasswordHasher() (*PasswordHasher, error) {
	return &PasswordHasher{}, nil
}

func (ph PasswordHasher) encode(rawPassword []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}

// check returns true if password is correct
func (ph PasswordHasher) check(hashedPassword []byte, rawPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(rawPassword))
	return err == nil
}
