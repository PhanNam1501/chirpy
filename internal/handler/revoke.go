package handler

import (
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
)

func (cfg *Config) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
