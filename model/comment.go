package model

import (
	"time"
)

type Comment struct {
	ID        int       `json:"id,omitempty"`
	PostID    int       `json:"post_id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}
