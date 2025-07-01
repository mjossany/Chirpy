package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mjossany/Chirpy/internal/auth"
	"github.com/mjossany/Chirpy/internal/database"
)

type userRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userItem struct {
	Id             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword sql.NullString
}

func (cfg *Config) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req userRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	now := time.Now().UTC()
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          req.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	resp := userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func (cfg *Config) HandleGetUserByEmail(r *http.Request, email string) (userItem, error) {
	user, err := cfg.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		return userItem{}, err
	}

	return userItem{
		Id:             user.ID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}, nil
}
