package api

import (
	"encoding/json"
	"net/http"

	"github.com/mjossany/Chirpy/internal/auth"
)

type loginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *Config) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, err := cfg.HandleGetUserByEmail(r, loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	userResponse := userResponse{
		ID:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}
