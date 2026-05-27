package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PhanNam1501/chirpy/internal/auth"
	"github.com/PhanNam1501/chirpy/internal/database"
	"github.com/PhanNam1501/chirpy/internal/util"
	"github.com/google/uuid"
)

type ChirpRequest struct {
	Body string `json:"body"`
	// UserID uuid.UUID `json:"user_id"`
}

type ChirpResponse struct {
	Id        uuid.UUID   `json:"id"`
	CreatedAt interface{} `json:"created_at"`
	UpdatedAt interface{} `json:"updated_at"`
	Body      string      `json:"body"`
	Error     string      `json:"error,omitempty"`
	UserId    uuid.UUID   `json:"user_id"`
}

func (cfg *Config) CreateChirp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, ChirpResponse{
			Error: "Something went wrong",
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

	var req ChirpRequest
	if err := json.Unmarshal(body, &req); err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, ChirpResponse{
			Error: "Something went wrong",
		})
		return
	}

	if len(req.Body) > 140 {
		util.RespondWithJSON(w, http.StatusBadRequest, ChirpResponse{
			Error: "Chirp is too long",
		})
		return
	}

	cleaned := util.CleanProfanity(req.Body)

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, ChirpResponse{
			Error: "Failed to create chirp",
		})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, ChirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *Config) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortStr := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			util.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid author_id",
			})
			return
		}

		chirps, err = cfg.DB.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		chirps, err = cfg.DB.GetChirps(r.Context())
	}

	if err != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get chirps",
		})
		return
	}

	if chirps == nil {
		chirps = []database.Chirp{}
	}

	chirpsResponse := make([]ChirpResponse, 0)
	if sortStr == "desc" {
		for i := len(chirps) - 1; i >= 0; i-- {
			chirpsResponse = append(chirpsResponse, ChirpResponse{
				Id:        chirps[i].ID,
				CreatedAt: chirps[i].CreatedAt,
				UpdatedAt: chirps[i].UpdatedAt,
				Body:      chirps[i].Body,
				UserId:    chirps[i].UserID,
			})
		}
	} else {
		for _, chirp := range chirps {
			chirpsResponse = append(chirpsResponse, ChirpResponse{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}
	}

	util.RespondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *Config) GetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid chirp ID",
		})
		return
	}

	chirp, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		util.RespondWithJSON(w, http.StatusNotFound, map[string]string{
			"error": "Chirp not found",
		})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, ChirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}
