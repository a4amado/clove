package apiguard

import (
	"net/http"
	"strings"
)

func GetHeaderApi(r *http.Request) string {
	token := r.Header.Get("Authorization")
	_, after, _ := strings.Cut(token, " ")
	return after
}
