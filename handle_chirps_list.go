package main

import "net/http"

func (cfg *apiConfig) handleChirpList(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't get chirps", err)
		return
	}

	responseChirps := make([]Chirp, len(chirps))
	for i, dbChirp := range chirps {
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
