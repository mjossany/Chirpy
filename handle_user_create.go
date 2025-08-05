package main

import (
	"encoding/json"
	"net/http"

	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleUserCreation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash password", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't create user", err)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, 201, user)
}
