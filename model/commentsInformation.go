package model

type CommentsInformation struct {
	PostID   int       `json:"post_id"`
	Comments []Comment `json:"comments"`
}
