package main

import(
	"net/http"
	"github.com/molkobahn/Chirpy/internal/auth"
	"time"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}	
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No authorization token found", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), authToken)	
	if err != nil {
		respondWithError(w, 401, "Couldn't find user to match token", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Failed to create token: %v", err)
		return
	}
	respondWithJSON(w, 200, response{
		Token: token,
	})
} 