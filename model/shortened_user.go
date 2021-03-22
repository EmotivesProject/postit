package model

//User struct declaration
type ShortenedUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	EncodedID string `json:"encoded_id"`
}
