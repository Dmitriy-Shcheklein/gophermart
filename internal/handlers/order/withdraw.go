package user

import (
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, domainErrors.InvalidContentTypeMsg, http.StatusBadRequest)
		return
	}
	userID, err := h.authService.GetUserID(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("error while getting userID")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var body models.RequestWithdrawn
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, domainErrors.DecodeBodyErrMsg, http.StatusBadRequest)
		return
	}
	if err = h.validate.Struct(body); err != nil {
		http.Error(w, domainErrors.ValidateBodyErrMsg, http.StatusBadRequest)
		return
	}
	if err = h.service.Withdraw(r.Context(), userID, body.Sum, body.Order); err != nil {
		if errors.Is(err, domainErrors.ErrOrderNotEnoughBalance) {
			http.Error(w, "not enough balance", http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, domainErrors.ErrOrderInvalidNumber) {
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
			return
		}
		h.logger.Error().Err(err).Msg("error while withdraw")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
