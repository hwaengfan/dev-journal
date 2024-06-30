package authenticationServices

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hwaengfan/dev-journal-backend/configs"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(configs.GlobalEnvironmentVariables.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, error := token.SignedString(secret)
	if error != nil {
		return "", error
	}

	return tokenString, nil
}
