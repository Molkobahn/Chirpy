package main

import(
	"net/http"
	"encoding/json"
	"github.com/molkobahn/Chirpy/internal/auth"
	"log"
	"time"
	"github.com/molkobahn/Chirpy/internal/database"
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
	var expiresIn time.Duration
	expiresIn = time.Duration(3600) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token: %v", err)
		return
	}
	refTokenName, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token: %v", err)
		return
	}
	now := time.Now()
	expiresIn, err = time.ParseDuration("1440h")
	arg := database.CreateRefreshTokenParams{
		Token:	refTokenName,
		UserID:	user.ID,
		ExpiresAt: now.Add(expiresIn),
	}
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token: %v", err)
		return
	}
	responseUser := mapUser(user)
	responseUser.Token = token
	responseUser.RefreshToken = refreshToken.Token
	respondWithJSON(w, http.StatusOK, responseUser)
}