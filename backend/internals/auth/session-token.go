package auth

import (
	envConsts "clove/internals/consts/env"
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SessionTokenClaim struct {
	User        repository.User `json:"app"`
	Permessions PermissionsBuilder

	jwt.RegisteredClaims
}

func GenerateSessionToken(user repository.User) (string, error) {
	perms := NewPermissionsBuilder()
	perms.Allow(RESOURCES_ALL, OPERATIONS_ALL)

	claims := SessionTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "Session Token",
			Issuer:    fmt.Sprintf("%v:%s", envConsts.Region(), "foo"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		User:        user,
		Permessions: perms,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(envConsts.JWTSecret()) // Convert to []byte
}

func ValidateSessionToken(tokenString string) (*SessionTokenClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&SessionTokenClaim{},
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

	claims, ok := parsedToken.Claims.(*SessionTokenClaim)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

const session_token_cookie_name = "euf_k15"

func ParseSessionFromRequest(r *http.Request) (*SessionTokenClaim, error) {
	isMobile := r.Header.Get("X-Client-Type") == "mobile"

	token := ""

	if isMobile {
		headerToken := r.Header.Get(authorization)
		token = extractBearerTokenFrom(headerToken)
	} else {
		session_token, err := r.Cookie(session_token_cookie_name)
		if err != nil {
			return nil, err
		}
		token = session_token.Value
	}

	return ValidateSessionToken(token)
}
