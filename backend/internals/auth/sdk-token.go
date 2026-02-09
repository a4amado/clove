package auth

import (
	envConsts "clove/internals/consts/env"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SdkTokenClaim struct {
	AppID       uuid.UUID `json:"app_id"`
	Permessions PermissionsBuilder
	jwt.RegisteredClaims
}

func GenerateSDKToken(appId uuid.UUID) (string, error) {
	perms := NewPermissionsBuilder()
	perms.Allow(DELIVERY, CREATE)
	perms.Allow(OneTimeToken, CREATE)

	claims := SdkTokenClaim{
		AppID:       appId,
		Permessions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Server NAME" + "REGION",
			Subject: fmt.Sprintf("App:[%s] SDK LONG LIVED TOKEN", appId),
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add((((time.Hour * 24) * 30) * 12) * 100),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(envConsts.JWTSecret())
}

func ValidateSDKToken(tokenString string) (*SdkTokenClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&SdkTokenClaim{},
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

	claims, ok := parsedToken.Claims.(*SdkTokenClaim)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
func ParseSDKTokenFromRequest(r *http.Request) (*SdkTokenClaim, error) {
	headerToken := r.Header.Get(authorization)
	token := extractBearerTokenFrom(headerToken)
	return ValidateSDKToken(token)
}
