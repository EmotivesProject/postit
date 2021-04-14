package event

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"postit/internal/logger"
	"postit/model"
)

func SendPostEvent(username string, post *model.Post) {
	event := model.Event{
		Username:      username,
		CustomerEvent: "post",
		Data:          post,
	}
	sendEvent(event)
}

func sendEvent(event model.Event) {
	baseHost := os.Getenv("BASE_HOST")
	url := baseHost + "metrics/customer_event_token"

	requestBody, err := json.Marshal(event)
	if err != nil {
		logger.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		logger.Error(err)
	}
	req.Header.Add("Authorization", "qutSecret")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
	}
	defer resp.Body.Close()
	logger.Info("Sent event to metrics")
}
