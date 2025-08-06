package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpList(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorUUID)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}

	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []Chirp{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	responseChirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		responseChirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	respondWithJSON(w, 200, responseChirps)
}
