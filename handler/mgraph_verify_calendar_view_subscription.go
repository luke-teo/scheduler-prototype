package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) MGraphVerifyCalendarViewSubscription(w http.ResponseWriter, r *http.Request) {
	validationToken := r.URL.Query().Get("validationToken")
	urlDecodedValidationToken, err := url.QueryUnescape(validationToken)
	log.Printf("validationToken: %s", validationToken)
	log.Printf("urlDecodedValidationToken: %s", urlDecodedValidationToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(urlDecodedValidationToken))
}
