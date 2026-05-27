package handler

import (
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/util"
	"github.com/google/uuid"
)

func (cfg *Config) DeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid chirp ID",
		})
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, ChirpResponse{
			Error: "Missing or invalid authorization",
		})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, ChirpResponse{
			Error: "Invalid token",
		})
		return
	}

	realUser, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		util.RespondWithJSON(w, http.StatusNotFound, map[string]string{
			"error": "Chirp not found",
		})
		return
	}

	if realUser.UserID != userID {
		util.RespondWithJSON(w, http.StatusForbidden, map[string]string{
			"error": "You can only delete your own chirps",
		})
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete chirp",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
