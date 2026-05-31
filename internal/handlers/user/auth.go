package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, domainErrors.InvalidContentTypeMsg, http.StatusBadRequest)
		return
	}
	var body models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, domainErrors.DecodeBodyErrMsg, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(body); err != nil {
		http.Error(w, domainErrors.ValidateBodyErrMsg, http.StatusBadRequest)
		return
	}
	jwtString, err := h.service.Auth(context.Background(), body.Login, body.Password)
	if errors.Is(err, domainErrors.ErrInvalidAuthData) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("error while auth user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", jwtString)
	w.WriteHeader(http.StatusOK)
}
