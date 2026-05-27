package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/util"
	"github.com/google/uuid"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type PolkaWebhookResponse struct {
	Error string
}

func (cfg *Config) PolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.ApiKey {
		util.RespondWithJSON(w, http.StatusUnauthorized, PolkaWebhookResponse{
			Error: "Api key is not matched",
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, PolkaWebhookResponse{
			Error: "Something went wrong",
		})
		return
	}

	var req PolkaWebhookRequest
	if err = json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, PolkaWebhookResponse{
			Error: "Cannot decode body",
		})
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = cfg.DB.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
