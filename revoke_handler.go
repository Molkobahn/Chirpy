package main

import(
	"net/http"
	"github.com/molkobahn/Chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No authorization token found", err)
		return
	}
	token, err := cfg.db.GetRefreshToken(r.Context(), authToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token not found", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
	}
	w.WriteHeader(204)
}

