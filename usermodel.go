package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type DatabaseUser interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetEmail() string
}

func mapUser[U DatabaseUser](u U, token string) User {
	if token == "" {
		return User{
			ID:        u.GetID(),
			CreatedAt: u.GetCreatedAt(),
			UpdatedAt: u.GetUpdatedAt(),
			Email:     u.GetEmail(),
		}
	}
	return User{
		ID:        u.GetID(),
		CreatedAt: u.GetCreatedAt(),
		UpdatedAt: u.GetUpdatedAt(),
		Email:     u.GetEmail(),
		Token:     token,
	}
}
