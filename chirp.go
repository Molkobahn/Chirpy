package main

import(
	"github.com/google/uuid"
	"time"
	"github.com/molkobahn/Chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body	string	`json:"body"`
	UserID	uuid.UUID	`json:"user_id"`
}

func mapChirp(chirp database.Chirp) Chirp {
	return Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:	chirp.Body,
		UserID:	chirp.UserID,
	}
}