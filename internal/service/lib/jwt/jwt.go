package jwt

import (
	configapp "auth/internal/config"
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

func ValidateToken(token string, st configapp.Config) bool {
	ken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
			return nil, fmt.Errorf("cyka")
		}
		return []byte(st.Secret), nil
	})
	if err != nil {
		return false
	}

	claim := ken.Claims.(jwt.MapClaims)
	tokenTime := claim["ext"]

	return deltaTime(tokenTime.(float64))
}

func deltaTime(tt float64) bool {
	ct := time.Now().Unix()

	if float64(ct) > tt {
		return true
	}
	return false

}
