package handler

import (
	"cooling-tower-diagnostics/internal/model"
	"errors"
	"net/http"
)

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func mapError(err error) (int, string) {
	if errors.Is(err, model.ErrInvalid) {
		return 400, "invalid_request"
	}
	if errors.Is(err, model.ErrNotFound) {
		return 404, "not_found"
	}
	if errors.Is(err, model.ErrConflict) {
		return 409, "conflict"
	}
	if errors.Is(err, model.ErrCanceled) {
		return 408, "canceled"
	}
	return 500, "internal_error"
}
func respondError(w http.ResponseWriter, r *http.Request, err error) {
	status, code := mapError(err)
	writeJSONIndent(w, status, ErrorBody{Code: code, Message: err.Error(), RequestID: requestID(r)})
}
func assertMethod(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	if methodAllowed(r, methods...) {
		return true
	}
	w.Header().Set("Allow", methods[0])
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	return false
}
