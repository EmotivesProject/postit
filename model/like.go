package model

import "time"

type Like struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey"`
	PostID    int       `json:"post_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Active    bool      `json:"active"`
}
