package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) MGraphHandleCalendarViewNotification(w http.ResponseWriter, r *http.Request) {
	// check if it is for verification
	validationToken := r.URL.Query().Get("validationToken")
	if validationToken != "" {
		log.Printf("validationToken exists")
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
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the raw JSON data
	fmt.Println("Received Raw JSON Request Body:")
	fmt.Println(string(body))

	// otherwise, check notification and update accordingly

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
