package authenticationServices

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hwaengfan/dev-journal-backend/configs"
	userModel "github.com/hwaengfan/dev-journal-backend/internal/models/user"
	"github.com/hwaengfan/dev-journal-backend/internal/utils"
)

type contextKey string

const UserKey contextKey = "userID"

// CreateJWT creates a new JWT token
func CreateJWT(secret []byte, userID uuid.UUID) (string, error) {
	expiration := time.Second * time.Duration(configs.GlobalEnvironmentVariables.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, error := token.SignedString(secret)
	if error != nil {
		return "", error
	}

	return tokenString, nil
}

// GetUserIDFromContext retrieves the userID from the context
func GetUserIDFromContext(ctx context.Context) uuid.NullUUID {
	userID, exists := ctx.Value(UserKey).(uuid.UUID)
	if !exists {
		return uuid.NullUUID{UUID: userID, Valid: false}
	}

	return uuid.NullUUID{UUID: userID, Valid: true}
}

// JWTAuthentication check for logged in users
func JWTAuthentication(handlerFunction http.HandlerFunc, userStore userModel.UserStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// get JWT token from request header
		tokenString := getTokenFromRequest(request)

		// parse the JWT token
		token, error := parseToken(tokenString)
		if error != nil {
			log.Printf("failed to parse token: %v", error)
			utils.WritePermissionDenied(writer)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			utils.WritePermissionDenied(writer)
			return
		}

		// extract userID if the JWT token is valid
		claims := token.Claims.(jwt.MapClaims)
		userIDString, ok := claims["userID"].(string)
		if !ok {
			log.Printf("userID in claims is not a UUID: %v", claims["userID"])
			utils.WritePermissionDenied(writer)
			return
		}

		userID, err := uuid.Parse(userIDString)
		if err != nil {
			log.Printf("failed to parse userID: %v", err)
			utils.WritePermissionDenied(writer)
			return
		}

		// get user by ID
		user, err := userStore.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by ID: %v", err)
			utils.WritePermissionDenied(writer)
			return
		}

		// add userID to the context
		ctx := request.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		request = request.WithContext(ctx)

		// call the handler function
		handlerFunction(writer, request)
	}
}

// getTokenFromRequest retrieves the JWT token from the request header
func getTokenFromRequest(request *http.Request) string {
	tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
	if tokenString != "" {
		return tokenString
	}
	return ""
}

// parseToken parses the JWT token
func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(configs.GlobalEnvironmentVariables.JWTSecret), nil
	})
}
