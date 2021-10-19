package model

type PostInformation struct {
	Post         Post      `json:"post"`
	Comments     []Comment `json:"comments"`
	EmojiSummary string    `json:"emoji_summary"`
	Likes        []Like    `json:"likes"`
}
