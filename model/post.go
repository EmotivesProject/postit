package model

import (
	"time"
)

type Post struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey"`
	Username  string    `json:"user"`
	Message   *string   `json:"message,omitempty"`
	ImagePath *string   `json:"image_path,omitempty"`
	Longitude *float64  `json:"longitude,omitempty"`
	Latitude  *float64  `json:"latitude,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Active    bool      `json:"active"`
}

func (p Post) Validate() bool {
	if p.Message == nil && p.ImagePath == nil && (p.Longitude == nil || p.Latitude == nil) {
		return false
	}
	return true
}
