package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scheduler-prototype/dto"
	requestDto "github.com/scheduler-prototype/dto/request"
	"github.com/scheduler-prototype/utility"
)

func (h *Handler) MGraphCalendarViewFirstSync(w http.ResponseWriter, r *http.Request) {
	// read the request body and insert into users table
	req := &requestDto.MGraphCalendarViewFirstSyncDto{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// insert into database if user does not exists
	newUserUuid, err := uuid.Parse(req.UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = h.repo.GetUserByUserId(&newUserUuid)
	if err != nil {
		if err == utility.ErrNotFound {
			// create user if not exists
			newUser := &dto.UserDto{
				UserId:        newUserUuid,
				CurrentDelta:  nil,
				PreviousDelta: nil,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			err = h.repo.CreateUser(newUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := map[string]string{"error": err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// once user is confirmed to be in database, grab the user and make first delta queries to Microsoft Graph
	userDto, err := h.repo.GetUserByUserId(&newUserUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// make first delta queries to Microsoft Graph

	// Get the current time in the user's timezone
	now := time.Now().UTC().Add(time.Duration(time.Hour * -8))

	// Calculate the start of the current month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Calculate the start of the next month
	nextMonth := now.AddDate(0, 1, 0)
	startOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Calculate the end of the next month
	endOfNextMonth := startOfNextMonth.Add(-time.Second)

	events, err := h.client.GetDeltaCalendarView(startOfMonth.Format(time.RFC3339), endOfNextMonth.Format(time.RFC3339), userDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Events: %v", events)
	w.WriteHeader(http.StatusOK)
}
