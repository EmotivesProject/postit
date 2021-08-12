package model

type PostInformation struct {
	Post       Post         `json:"post"`
	Comments   []Comment    `json:"comments"`
	EmojiCount []EmojiCount `json:"emoji_count"`
	Likes      []Like       `json:"likes"`
}
