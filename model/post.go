package model

import (
	"time"
)

type Post struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Content   JSONB     `json:"content" sql:"type:jsonb"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Active    bool      `json:"active"`
}
