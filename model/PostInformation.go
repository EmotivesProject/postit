package model

type PostInformation struct {
	Post           Post      `json:"post"`
	Comments       []Comment `json:"comments"`
	EmojiCount     string    `json:"emoji_count"`
	SelfEmojiCount string    `json:"self_emoji_count"`
	Likes          []Like    `json:"likes"`
}
