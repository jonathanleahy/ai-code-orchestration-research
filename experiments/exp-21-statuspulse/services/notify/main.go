package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"statuspulse/notify/store"
)

var s *store.Store

func main() {
	s = store.NewStore()

	http.HandleFunc("/api/subscribers", handleSubscribers)
	http.HandleFunc("/api/notify", handleNotify)
	http.HandleFunc("/health", handleHealth)

	log.Println("Notify service listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", http.HandlerFunc(withCORS)))
}

func withCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.DefaultServeMux.ServeHTTP(w, r)
}

func handleSubscribers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handleAddSubscriber(w, r)
	} else if r.Method == "GET" {
		handleListSubscribers(w, r)
	} else if r.Method == "DELETE" {
		handleDeleteSubscriber(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddSubscriber(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string   `json:"name"`
		WebhookURL string   `json:"webhook_url"`
		Events     []string `json:"events"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscriber, err := s.AddSubscriber(req.Name, req.WebhookURL, req.Events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscriber)
}

func handleListSubscribers(w http.ResponseWriter, r *http.Request) {
	subscribers := s.ListSubscribers()
	response := map[string]interface{}{
		"subscribers": subscribers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/subscribers/")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	err := s.DeleteSubscriber(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Event   string      `json:"event"`
		Payload interface{} `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscribers := s.GetMatchingSubscribers(req.Event)

	var results []store.DispatchResult
	dispatched := 0
	failed := 0

	for _, subscriber := range subscribers {
		result := dispatchWebhook(subscriber, req.Event, req.Payload)
		results = append(results, result)

		if result.Success {
			dispatched++
		} else {
			failed++
		}
	}

	response := map[string]interface{}{
		"dispatched": dispatched,
		"failed":     failed,
		"results":    results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func dispatchWebhook(subscriber *store.Subscriber, event string, payload interface{}) store.DispatchResult {
	result := store.DispatchResult{
		SubscriberID: subscriber.ID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		result.Success = false
		result.Error = "Failed to marshal payload"
		return result
	}

	req, err := http.NewRequest("POST", subscriber.WebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event", event)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	io.ReadAll(resp.Body)

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "notify",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
