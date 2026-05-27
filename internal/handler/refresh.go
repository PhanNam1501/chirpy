package handler

import (
	"net/http"
	"time"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/util"
)

type RefreshResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func (cfg *Config) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, RefreshResponse{
			Error: "Missing or invalid authorization header",
		})
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, RefreshResponse{
			Error: "Invalid refresh token",
		})
		return
	}

	if time.Now().UTC().After(user.ExpiresAt) {
		util.RespondWithJSON(w, http.StatusUnauthorized, RefreshResponse{
			Error: "Refresh token expired",
		})
		return
	}

	if user.RevokedAt.Valid {
		util.RespondWithJSON(w, http.StatusUnauthorized, RefreshResponse{
			Error: "Refresh token revoked",
		})
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, RefreshResponse{
			Error: "Failed to create access token",
		})
		return
	}

	// Return new access token
	util.RespondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: newAccessToken,
	})
}
