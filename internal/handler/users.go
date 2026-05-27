package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/database"
	"github.com/PhanNam1501/chirpy/internal/util"
)

type AddUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AddUserResponse struct {
	ID          interface{} `json:"id"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
	Email       string      `json:"email"`
	IsChirpyRed bool        `json:"is_chirpy_red"`
	Error       string      `json:"error,omitempty"`
}

func (cfg *Config) AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, AddUserResponse{
			Error: "Failed to create user",
		})
		return
	}

	var req AddUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, AddUserResponse{
			Error: "Failed to create user",
		})
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, AddUserResponse{
			Error: "Failed to create user",
		})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, AddUserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
