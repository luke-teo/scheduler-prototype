package handler

import (
	"encoding/json"
	"net/http"

	requestDto "github.com/scheduler-prototype/dto/request"
)

func (h *Handler) MGraphCreateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this handler
	// read the request body and create a MGraphCreateEventDto
	req := &requestDto.MGraphCreateEventDto{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// create request to Microsoft Graph to create the event
	event, err := h.client.PostCreateEvent(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}
