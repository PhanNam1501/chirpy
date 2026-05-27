package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/util"
)

type ValidateChirpRequest struct {
	Body string `json:"body"`
}

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (cfg *Config) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, ValidateChirpResponse{
			Error: "Something went wrong",
		})
		return
	}

	var req ValidateChirpRequest
	if err := json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, ValidateChirpResponse{
			Error: "Something went wrong",
		})
		return
	}

	if len(req.Body) > 140 {
		util.RespondWithJSON(w, http.StatusBadRequest, ValidateChirpResponse{
			Error: "Chirp is too long",
		})
		return
	}

	cleaned := util.CleanProfanity(req.Body)
	util.RespondWithJSON(w, http.StatusOK, ValidateChirpResponse{
		CleanedBody: cleaned,
	})
}
