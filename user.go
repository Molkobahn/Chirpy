package main

import (
	"github.com/google/uuid"
	"time"
	"github.com/molkobahn/Chirpy/internal/database"
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