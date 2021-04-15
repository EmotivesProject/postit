package model

const (
	StatusCreated = "created"
	StatusDeleted = "deleted"
)

type EventData struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}
