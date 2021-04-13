package api

import (
	"encoding/json"
	"net/http"
	"postit/internal/logger"
	"postit/model"
)

func resultResponseJSON(w http.ResponseWriter, status int, result interface{}) {
	response := model.Response{
		Result: result,
	}

	responseJSON(w, status, response)
}

func messageResponseJSON(w http.ResponseWriter, status int, message model.Message) {
	response := model.Response{
		Message: []model.Message{message},
	}

	responseJSON(w, status, response)
}

func responseJSON(w http.ResponseWriter, status int, response interface{}) {
	payload, err := json.Marshal(response)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Sending response %v", response)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
	_, err = w.Write(payload)

	if err != nil {
		logger.Error(err)
	}
}
