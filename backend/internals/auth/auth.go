package auth

import (
	"errors"
	"net/http"
)

type SessionType string

const (
	OOT            SessionType = "OOT"
	SDKToken       SessionType = "SDK_TOKEN"
	RegualrSession SessionType = "REGULAR_SESSION"
)

type Session struct {
	Permissions PermissionsBuilder
	SessionType SessionType
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ParseSession(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func UnAuthResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
func ParseSession(r *http.Request) (*Session, error) {

	OOTSession, _ := ParseOneTimeTokenFromRequest(r)

	if OOTSession != nil {
		return &Session{
			Permissions: OOTSession.Permessions,
			SessionType: OOT,
		}, nil
	}

	SDKSession, _ := ParseSDKTokenFromRequest(r)
	if SDKSession != nil {
		return &Session{
			Permissions: SDKSession.Permessions,
			SessionType: SDKToken,
		}, nil
	}

	s, _ := ParseSessionFromRequest(r)
	if s != nil {
		return &Session{
			Permissions: s.Permessions,
			SessionType: RegualrSession,
		}, nil
	}

	return nil, errors.New("unauthrized")
}
