package order

import (
	"errors"
	"io"
	"net/http"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
)

// Upload загрузка заказов пользователя
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
		http.Error(w, domainErrors.InvalidContentTypeMsg, http.StatusBadRequest)
		return
	}
	num, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("error while read body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error().Err(err).Msg("error while close body")
		}
	}()
	if len(num) == 0 {
		http.Error(w, domainErrors.ValidateBodyErrMsg, http.StatusBadRequest)
		return
	}

	userID, err := h.authService.GetUserID(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("error while getting userID")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = h.service.Upload(r.Context(), userID, string(num)); err != nil {
		if errors.Is(err, domainErrors.ErrOrderAlreadyExists) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, domainErrors.ErrOrderBelongsAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, domainErrors.ErrOrderInvalidNumber) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		h.logger.Error().Err(err).Msg("error while upload orderNum")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
