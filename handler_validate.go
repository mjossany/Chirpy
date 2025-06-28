package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handleChirpValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := bodyCleaning(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
}

func bodyCleaning(chirpy string) string {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	splittedChirpy := strings.Split(chirpy, " ")
	for index, word := range splittedChirpy {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			splittedChirpy[index] = "****"
		}
	}
	return strings.Join(splittedChirpy, " ")
}
