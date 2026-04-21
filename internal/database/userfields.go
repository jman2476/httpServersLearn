package database

import (
	"time"

	"github.com/google/uuid"
)

func (u User) GetID() uuid.UUID {
	return u.ID
}

func (u CreateUserRow) GetID() uuid.UUID {
	return u.ID
}

func (u UpdateUserByIDRow) GetID() uuid.UUID {
	return u.ID
}

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u CreateUserRow) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u UpdateUserByIDRow) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u CreateUserRow) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u UpdateUserByIDRow) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u User) GetEmail() string {
	return u.Email
}

func (u CreateUserRow) GetEmail() string {
	return u.Email
}

func (u UpdateUserByIDRow) GetEmail() string {
	return u.Email
}
