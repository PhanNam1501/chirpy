package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/database"
	"github.com/PhanNam1501/chirpy/internal/util"
)

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	// ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoginResponse struct {
	ID           interface{} `json:"id"`
	CreatedAt    interface{} `json:"created_at"`
	UpdatedAt    interface{} `json:"updated_at"`
	Email        string      `json:"email"`
	IsChirpyRed  bool        `json:"is_chirpy_red"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	Error        string      `json:"error,omitempty"`
}

func (cfg *Config) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, LoginResponse{
			Error: "Failed to create user",
		})
		return
	}

	var req LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, LoginResponse{
			Error: "Failed to create user",
		})
		return
	}
	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, LoginResponse{
			Error: "Failed to get user",
		})
		return
	}
	check, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, LoginResponse{
			Error: "Failed to checkPasswordHash",
		})
		return
	}
	if !check {
		util.RespondWithJSON(w, http.StatusUnauthorized, LoginResponse{
			Error: "Incorrect email or password",
		})
		return
	}
	expiresIn := time.Hour // default 1 hour
	// if req.ExpiresInSeconds > 0 {
	// 	expiresIn = time.Duration(req.ExpiresInSeconds) * time.Second
	// 	if expiresIn > time.Hour {
	// 		expiresIn = time.Hour
	// 	}
	// }

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, LoginResponse{
			Error: "Failed to create token",
		})
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, LoginResponse{
			Error: "Failed to create refresh token",
		})
		return
	}

	expiresAt := time.Now().UTC().Add(60 * 24 * time.Hour) // 60 days
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		// revoked_at should be NULL (don't set it)
	})
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, LoginResponse{
			Error: "Failed to save refresh token",
		})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
	})

}
