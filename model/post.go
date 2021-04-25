package model

import (
	"time"
)

type Post struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username"`
	Content   JSONB     `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}
