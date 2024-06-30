package authenticationServices

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if error != nil {
		return "", error
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, plain []byte) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plain)
	return error == nil
}
