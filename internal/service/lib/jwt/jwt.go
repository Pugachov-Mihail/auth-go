package jwt

import (
	"auth/internal/domain/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["ext"] = time.Now().Add(duration).Unix()

	tokenSecret, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return tokenSecret, nil
}
