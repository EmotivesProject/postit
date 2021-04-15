package model

const (
	EventPost    = "post"
	EventLike    = "like"
	EventComment = "comment"
)

type Event struct {
	Username      string      `json:"username"`
	CustomerEvent string      `json:"customer_event"`
	Data          interface{} `json:"event_data"`
}
