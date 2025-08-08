package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/google/uuid"
	"github.com/molkobahn/Chirpy/internal/database"
)

type cleanedParameters struct {
	CleanedBody string `json:"cleaned_body"`
	UserID	uuid.UUID	`json:"user_id"`
}

func replaceBadWord(body string) string {
	bodySlice := strings.Split(body, " ")
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	for _, badWord := range profaneWords {
		for i, word := range bodySlice {
			if strings.ToLower(word) == badWord {
				bodySlice[i] = strings.Replace(word, word, "****", -1)
			}
		}
	}
	return strings.Join(bodySlice, " ")
}

func validateChirp(w http.ResponseWriter, r *http.Request) database.CreateChirpParams {
	type parameters struct {
		Body string `json:"body"`
		UserID	uuid.UUID	`json:"user_id"`
	}
	type returnVals struct {
		Valid bool `json"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return database.CreateChirpParams{}
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return database.CreateChirpParams{}
	}
	return database.CreateChirpParams{
		Body: replaceBadWord(params.Body),
		UserID:	params.UserID,
	}
}

func (cfg *apiConfig)chirpHandler(w http.ResponseWriter, r *http.Request) {
	arg := validateChirp(w, r)
	chirp, err := cfg.db.CreateChirp(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, mapChirp(chirp))
}

func (cfg *apiConfig)getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context()) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}
	respondWithArray(w, http.StatusOK, chirps)
}

func (cfg *apiConfig)getChirpHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, 404, "Failed to get chirp", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Failed to get chirp", err)
		return
	}
	mappedChirp := mapChirp(chirp)
	respondWithJSON(w, http.StatusOK, mappedChirp)
}