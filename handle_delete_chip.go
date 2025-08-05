package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, 404, "chirpID must not be blank", nil)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Couldn't find chirp", err)
			return
		}
		respondWithError(w, 500, "Couldn't find chirp", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, 403, "Unauthorized action", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't delete chirp", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
