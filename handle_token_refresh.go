package main

import (
	"net/http"
	"time"

	"github.com/mjossany/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 404, "Invalid authorization", err)
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Couldn't validate token", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Couldn't validate token", err)
		return
	}

	respondWithJSON(w, 200, response{
		Token: accessToken,
	})
}
