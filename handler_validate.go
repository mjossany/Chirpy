package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleChirpsValidation(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(params.Body, " ")

	for i, word := range chirpWords {
		for _, pw := range profaneWords {
			if pw == strings.ToLower(word) {
				chirpWords[i] = "****"
				break
			}
		}
	}

	validatedChirp := strings.Join(chirpWords, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: validatedChirp,
	})
}
