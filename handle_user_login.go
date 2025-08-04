package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
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

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Incorrect email or password", err)
			return
		}
		respondWithError(w, 500, "Error getting user", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	accessTokenExpiresIn := time.Hour
	jwt, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, accessTokenExpiresIn)
	if err != nil {
		respondWithError(w, 500, "Couldn't generate jwt authorization", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Couldn't create refresh token", err)
		return
	}

	refreshTokenExpiresIn := time.Hour * 24 * 60

	dbRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiresIn),
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't insert refresh token in database", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        jwt,
		RefreshToken: dbRefreshToken.Token,
	})
}
