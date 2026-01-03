package tokenguard

import (
	envConsts "clove/internals/consts/env"
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type OneTimeTokenClaim struct {
	App       repository.App `json:"app"`
	ChannelID string         `json:"channel_id"`
	jwt.RegisteredClaims
}

func GenerateOneTimeToken(app repository.App, channelID string, keyId uuid.UUID) (string, error) {
	claims := OneTimeTokenClaim{
		App:       app,
		ChannelID: channelID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "Clove One-Time Token",
			Issuer:    fmt.Sprintf("%v:%s", envConsts.Region(), keyId.String()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(envConsts.JWTSecret()) // Convert to []byte
}

func ValidateOneTimeToken(tokenString string) (*OneTimeTokenClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&OneTimeTokenClaim{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return envConsts.JWTSecret(), nil // Return []byte, not string
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*OneTimeTokenClaim)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
