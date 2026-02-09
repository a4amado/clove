package auth

import (
	envConsts "clove/internals/consts/env"
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const authorization = "Authorization"

func extractBearerTokenFrom(token string) string {
	_, after, ok := strings.Cut(token, "Bearer ")
	if !ok {
		return ""
	}
	return after
}

type OneTimeTokenClaim struct {
	App         repository.App `json:"app"`
	ChannelID   string         `json:"channel_id"`
	Permessions PermissionsBuilder
	jwt.RegisteredClaims
}

func GenerateOneTimeToken(app repository.App, channelID string) (string, error) {
	perms := NewPermissionsBuilder()
	perms.Allow(DELIVERY, READ)
	claims := OneTimeTokenClaim{
		App:       app,
		ChannelID: channelID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "Clove One-Time Token",
			Issuer:    fmt.Sprintf("%v:%s", envConsts.Region()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		Permessions: perms,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(envConsts.JWTSecret())) // Fixed: convert to []byte
}

func ValidateOneTimeToken(tokenString string) (*OneTimeTokenClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&OneTimeTokenClaim{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(envConsts.JWTSecret()), nil // Fixed: convert to []byte
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

func ParseOneTimeTokenFromRequest(r *http.Request) (*OneTimeTokenClaim, error) {
	authHeader := r.Header.Get(authorization)
	tokenString := extractBearerTokenFrom(authHeader)
	return ValidateOneTimeToken(tokenString)
}
