package model

type Message struct {
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}
