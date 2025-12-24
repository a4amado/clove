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

// UpdateUser handles HTTP requests to update a user's email identified by the `user_id` path parameter.
// 
// It validates that `user_id` is a UUID, rejects request bodies larger than 1MB with 413, and decodes a JSON
// payload into UpdateUserBody. The request must include `Email` and must not include both `Email` and `Password`
// at the same time; otherwise it responds 400. On success it updates the user's email via the repository and
// returns 200 with a JSON confirmation message. If the user is not found it returns 404; other repository errors
// result in 500.
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