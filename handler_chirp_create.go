package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpCreation(w http.ResponseWriter, r *http.Request) {
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

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedChirp := cleanChirp(params.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func cleanChirp(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(body, " ")

	for i, word := range chirpWords {
		for _, pw := range profaneWords {
			if pw == strings.ToLower(word) {
				chirpWords[i] = "****"
				break
			}
		}
	}

	return strings.Join(chirpWords, " ")
}
