package main

import(
	"net/http"
	"encoding/json"
	"github.com/molkobahn/Chirpy/internal/auth"
	"log"
	"time"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email", err)
		return
	}
	log.Printf("The user: %v", user)
	err = auth.CheckPasswordHash(params.Password, user.HashedPasswords)
	if err != nil {
		respondWithError(w, 401, "Incorrect password", err)
		return
	}
	var expiresIn time.Duration
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		expiresIn = time.Duration(3600) * time.Second
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token: %v", err)
	}
	responseUser := mapUser(user)
	responseUser.Token = token
	respondWithJSON(w, http.StatusOK, responseUser)
}