package model

type User struct {
	ID        int    `json:"id,omitempty"`
	Username  string `json:"username"`
	UserGroup string `json:"user_group"`
}
