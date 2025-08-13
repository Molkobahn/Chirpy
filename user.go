package main

import (
	"github.com/google/uuid"
	"time"
	"github.com/molkobahn/Chirpy/internal/database"
	"net/http"
	"encoding/json"
	"github.com/molkobahn/Chirpy/internal/auth"
)

type User struct {
	ID				uuid.UUID	`json:"id"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	Email			string		`json:"email"`
	Token			string		`json:"token"`
	RefreshToken	string		`json:"refresh_token"`
	IsChirpyRed		bool		`json:"is_chirpy_red"`
}

func mapUser(user database.User) User {
	return User{
		ID:	user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt:	user.UpdatedAt,
		Email:	user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}
}

func (cfg *apiConfig)createUserHandler(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
	}
	arg := database.CreateUserParams{
		Email:	params.Email,
		HashedPasswords:	hashedPassword,
	}
	user, err := cfg.db.CreateUser(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
		}
	newUser := mapUser(user)
	respondWithJSON(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No authentication token found", err)
		return
	}
	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Failed to validate token", err)
		return
	}
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}
	arg := database.UpdateUserParams{
		Email:	params.Email,
		HashedPasswords: hashedPassword,
		ID: userID,
	}
	updatedUser, err := cfg.db.UpdateUser(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update User", err)
		return
	}
	respondWithJSON(w, 200, mapUser(updatedUser))
}

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	authKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "No authorization key found", err)
		return
	}
	if authKey != cfg.polka_key {
		respondWithError(w, 401, "Authorization key doesn't match", err)
		return 
	}
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse user ID", err)
		return
	}
	err = cfg.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "Couldn't upgrade user", err)
		return
	}
	w.WriteHeader(204)
}