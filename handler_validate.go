package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json"valid"`
	}
	type profaneParams struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		responWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, profaneParams{
		CleanedBody: replaceBadWord(params.Body),
	})
}