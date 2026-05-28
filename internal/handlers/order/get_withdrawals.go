package user

import (
	"encoding/json"
	"net/http"
)

// GetWithdrawals получение списка списаний баллов пользователя
func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserID(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("error while getting userID")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	list, err := h.service.GetWithdrawals(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("error while getList")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(list)
	if err != nil {
		h.logger.Error().Err(err).Msg("error while serialized request")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("error while write body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
