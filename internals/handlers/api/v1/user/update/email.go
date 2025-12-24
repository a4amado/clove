package UserHandlersV1

import (
	dbPool "clove/internals/data/database/pool"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	httpheaders "github.com/utain/httpheaders"
	mediatypes "github.com/utain/httpheaders/mediatypes"
)

var db = repository.New(dbPool.Client())

type UpdateUserBody struct {
	Email    *string
	Password *string
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// this endpoint will always return json
	w.Header().Add(httpheaders.ContentType, mediatypes.ApplicationJson)

	user_id := r.PathValue("user_id")
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1000) // 1MB limit
	defer r.Body.Close()

	var body UpdateUserBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		maxBytesErr := &http.MaxBytesError{}
		if errors.As(err, &maxBytesErr) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now validate the body
	if body.Email == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.Email != nil && body.Password != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.UpdateUserEmail(context.Background(), repository.UpdateUserEmailParams{
		Email:  *body.Email,
		UserID: userUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"message": "user updated successfully"}`)

}
