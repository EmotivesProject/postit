package model

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	Username string `bson:"username" json:"username"`
}
