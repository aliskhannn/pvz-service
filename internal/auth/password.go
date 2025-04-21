package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	CheckPassword(password, hash string) error
}

type bcryptHasher struct{}

func NewBcryptHasher() PasswordHasher {
	return &bcryptHasher{}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (b *bcryptHasher) CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err
}
