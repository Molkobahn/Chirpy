package main

import(
	"net/http"
	"encoding/json"
	"github.com/molkobahn/Chirpy/internal/auth"
	"log"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
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
	respondWithJSON(w, http.StatusOK, mapUser(user))
}