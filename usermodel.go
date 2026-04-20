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
	Refresh   string    `json:"refresh_token"`
}

type DatabaseUser interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetEmail() string
}

func mapUser[U DatabaseUser](u U, access, refresh string) User {
	if access == "" {
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
		Token:     access,
		Refresh:   refresh,
	}
}
