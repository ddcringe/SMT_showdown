package models

import "time"

type User struct {
	ID           int       `json:"-" db:"id"`
	Username     string    `json:"-" db:"username"`
	Email        string    `json:"-" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastLogin    time.Time `json:"-" db:"last_login"`
	Bio          string    `json:"bio" db:"bio"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
}

type UserProfile struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Bio       string    `json:"bio" db:"bio"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
}
