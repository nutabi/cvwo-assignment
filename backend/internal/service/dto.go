package service

import "time"

type UserProfile struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	AvatarUrl *string    `json:"avatar_url"`
	Bio       *string    `json:"bio"`
}
