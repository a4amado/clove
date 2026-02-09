package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	ERROR_LIST_APP_KEYS_INVALID_ID            = "ERROR_LIST_APP_KEYS_INVALID_ID"
	ERROR_LIST_APP_KEYS_PAGE_IDX_NOT_A_NUMBER = "ERROR_LIST_APP_KEYS_PAGE_IDX_NOT_A_NUMBER"
)

func ListAppApiKeys(w http.ResponseWriter, r *http.Request) {
	session, err := auth.ParseSessionFromRequest(r)
	if err != nil {
		auth.UnAuthResponse(w)
		return
	}
	if !session.Permessions.Can(auth.KEY, auth.READ) {
		auth.UnAuthResponse(w)
		return
	}
	appId, err := uuid.Parse(r.PathValue("app_id"))

	if err != nil || appId == uuid.Nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_LIST_APP_KEYS_INVALID_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	pageIdxStr := r.URL.Query().Get("page_idx")
	if pageIdxStr == "" {
		pageIdxStr = "0"
	}
	page_idx, err := strconv.ParseInt(pageIdxStr, 10, 64)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_LIST_APP_KEYS_PAGE_IDX_NOT_A_NUMBER,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	srvs := services.New(r.Context()).WithCache()

	keys, err := srvs.Apps.Keys.List(appservice.ListKeysParams{
		AppId: appId,
		Page:  int32(page_idx),
	})

	if err != nil {
		http.Error(w, "Insternal server error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
