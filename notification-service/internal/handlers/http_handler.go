// handlers/http_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"notification-service/internal/services"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) HandleErrorNotification(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Error    string `json:"error"`
		Metadata string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SendErrorNotification(r.Context(), body.Error, body.Metadata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "notification sent"})
}
