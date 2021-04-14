package model

type Event struct {
	Username      string      `json:"username"`
	CustomerEvent string      `json:"customer_event"`
	Data          interface{} `json:"event_data"`
}
