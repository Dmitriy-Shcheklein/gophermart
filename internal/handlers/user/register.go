package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, domainErrors.InvalidContentTypeMsg, http.StatusBadRequest)
		return
	}
	var body models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, domainErrors.DecodeBodyErrMsg, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(body); err != nil {
		http.Error(w, domainErrors.ValidateBodyErrMsg, http.StatusBadRequest)
		return
	}
	err := h.service.Register(context.Background(), body.Login, body.Password)
	if errors.Is(err, domainErrors.ErrLoginDuplicate) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("error while register user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
