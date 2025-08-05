package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	authorization, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid authorization", err)
		return
	}

	userID, err := auth.ValidateJWT(authorization, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid authorization", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash password", err)
		return
	}

	dbUser, err := cfg.db.UpdateUserLoginInfo(r.Context(), database.UpdateUserLoginInfoParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 404, "Couldn't find user", err)
		return
	}

	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, 200, response{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	})
}
