package main

import (
	"net/http"

	"github.com/mjossany/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	authorization, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 404, "Couldn't find token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), authorization)
	if err != nil {
		respondWithError(w, 500, "Couldn't revoke token", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
