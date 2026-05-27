package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/database"
	"github.com/PhanNam1501/chirpy/internal/util"
)

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	ID          interface{} `json:"id"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
	Email       string      `json:"email"`
	IsChirpyRed bool        `json:"is_chirpy_red"`
	Error       string      `json:"error,omitempty"`
}

func (cfg *Config) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, UpdateUserResponse{
			Error: "Missing or invalid authorization",
		})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, UpdateUserResponse{
			Error: "Invalid token",
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, UpdateUserResponse{
			Error: "Something went wrong",
		})
		return
	}
	var req UpdateUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, UpdateUserResponse{
			Error: "Something went wrong",
		})
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, UpdateUserResponse{
			Error: "Failed to hash password",
		})
		return
	}

	user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, UpdateUserResponse{
			Error: "Failed to update user",
		})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, UpdateUserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
