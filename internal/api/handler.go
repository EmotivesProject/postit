package api

import (
	"net/http"
	"postit/model"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: "Health ok"})
}
