package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	requestDto "github.com/scheduler-prototype/dto/request"
)

func (h *Handler) MGraphCreateCalendarViewSubscription(w http.ResponseWriter, r *http.Request) {
	req := &requestDto.MGraphCreateCalendarViewSubscriptionDto{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// create request to Microsoft Graph to create the event
	subscriptions, err := h.client.CreateCalendarViewSubscription(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("subscriptions: %s", subscriptions)

	// if subscription, add subscription id to user (under the assumption that this user will only have one subscription to calendar)
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	userDto, err := h.repo.GetUserByUserId(&userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	userDto.SubscriptionId = subscriptions.GetId()
	userDto.SubscriptionExpiresAt = subscriptions.GetExpirationDateTime()
	log.Printf("userDto.SubscriptionId: %s", *userDto.SubscriptionId)
	err = h.repo.UpdateSubscriptionIdByUser(&userDto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subscriptions)
}
