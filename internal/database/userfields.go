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

// func (u UpgradeToRedByIDRow) GetID() uuid.UUID {
// 	return u.ID
// }

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u CreateUserRow) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u UpdateUserByIDRow) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// func (u UpgradeToRedByIDRow) GetCreatedAt() time.Time {
// 	return u.CreatedAt
// }

func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u CreateUserRow) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u UpdateUserByIDRow) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// func (u UpgradeToRedByIDRow) GetUpdatedAt() time.Time {
// 	return u.UpdatedAt
// }

func (u User) GetEmail() string {
	return u.Email
}

func (u CreateUserRow) GetEmail() string {
	return u.Email
}

func (u UpdateUserByIDRow) GetEmail() string {
	return u.Email
}

// func (u UpgradeToRedByIDRow) GetEmail() string {
// 	return u.Email
// }

func (u User) GetChirpyRed() bool {
	return u.IsChirpyRed
}

func (u CreateUserRow) GetChirpyRed() bool {
	return u.IsChirpyRed
}

func (u UpdateUserByIDRow) GetChirpyRed() bool {
	return u.IsChirpyRed
}

// func (u UpgradeToRedByIDRow) GetChirpyRed() bool {
// 	return u.IsChirpyRed
// }
