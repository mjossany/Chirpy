package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mjossany/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, 401, "Invalid authorization", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID", err)
		return
	}
	_, err = cfg.db.UpdateUserChirpyRed(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "User can't be found", err)
			return
		}
		respondWithError(w, 500, "Couldn't update user chirpy red", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
