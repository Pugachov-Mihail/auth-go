package tests

import (
	auth "auth/protos/gen/dota_traker.auth.v1"
	"auth/tests/suite"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRegisterLogin(t *testing.T) {
	ctx, st := suite.New(t)

	login := gofakeit.Email()
	pass := randomPassword()
	steamId := gofakeit.Int32()

	respReq, err := st.AuthClient.AuthRegistration(ctx, &auth.AuthRegistrationRequest{
		Password:  pass,
		Password2: pass,
		Login:     login,
		Email:     login,
		SteamId:   int64(steamId),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	respLogin, err := st.AuthClient.AuthLogin(ctx, &auth.AuthLoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	tokenString := respLogin.GetToken()

	require.NotEmpty(t, tokenString)

	ken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
			return nil, fmt.Errorf("cyka")
		}
		return []byte(st.Cfg.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := ken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReq.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, login, claims["email"].(string))

	const deltaTime = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["ext"].(float64), deltaTime)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 10)
}
