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
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

func mapUser(user database.User) User {
	return User{
		ID:	user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt:	user.UpdatedAt,
		Email:	user.Email,
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