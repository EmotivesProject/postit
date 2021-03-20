package model

type Response struct {
	Result  interface{} `json:"result"`
	Message []Message   `json:"message"`
}
