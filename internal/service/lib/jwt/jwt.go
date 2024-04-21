package jwt

import (
	"auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	clims := token.Claims.(jwt.MapClaims)
	clims["uid"] = user.Id
	clims["email"] = user.Email
	clims["ext"] = time.Now().Add(duration).Unix()

	tokenSecret, err := token.SignedString([]byte(app))
	if err != nil {
		return "", err
	}
	return tokenSecret, nil
}
