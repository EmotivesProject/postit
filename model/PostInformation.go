package model

type PostInformation struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
	Likes    []Like    `json:"likes"`
}
