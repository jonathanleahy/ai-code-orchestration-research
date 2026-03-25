package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"statuspulse/incidents/store"
)

var s *store.Store

func init() {
	s = store.NewStore()
}

func main() {
	http.HandleFunc("/api/incidents", handleIncidents)
	http.HandleFunc("/api/incidents/", handleIncidentDetail)
	http.HandleFunc("/health", handleHealth)

	log.Println("Starting incidents service on :8082")
	log.Fatal(http.ListenAndServe(":8082", http.HandlerFunc(corsMiddleware)))
}

func corsMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.DefaultServeMux.ServeHTTP(w, r)
}

func handleIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleCreateIncident(w, r)
	case http.MethodGet:
		handleListIncidents(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

func handleIncidentDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/incidents/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		return
	}

	id := parts[0]

	// GET /api/incidents/{id}
	if len(parts) == 1 && r.Method == http.MethodGet {
		handleGetIncident(w, r, id)
		return
	}

	// PATCH /api/incidents/{id}
	if len(parts) == 1 && r.Method == http.MethodPatch {
		handleUpdateIncident(w, r, id)
		return
	}

	// DELETE /api/incidents/{id}
	if len(parts) == 1 && r.Method == http.MethodDelete {
		handleDeleteIncident(w, r, id)
		return
	}

	// POST /api/incidents/{id}/resolve
	if len(parts) == 2 && parts[1] == "resolve" && r.Method == http.MethodPost {
		handleResolveIncident(w, r, id)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
}

func handleCreateIncident(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Severity    string `json:"severity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	status := store.IncidentStatus(req.Status)
	severity := store.Severity(req.Severity)

	incident, err := s.Create(req.Title, req.Description, status, severity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(incident)
}

func handleListIncidents(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")
	openParam := r.URL.Query().Get("open")

	var status *store.IncidentStatus
	if statusParam != "" {
		s := store.IncidentStatus(statusParam)
		status = &s
	}

	openOnly := openParam == "true" || openParam == "1"

	incidents := s.List(status, openOnly)
	if incidents == nil {
		incidents = make([]*store.Incident, 0)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"incidents": incidents})
}

func handleGetIncident(w http.ResponseWriter, r *http.Request, id string) {
	incident, err := s.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "incident not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incident)
}

func handleUpdateIncident(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Message     string `json:"message"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	var status *store.IncidentStatus
	var description *string
	var message *string

	if req.Status != "" {
		s := store.IncidentStatus(req.Status)
		status = &s
	}
	if req.Description != "" {
		description = &req.Description
	}
	if req.Message != "" {
		message = &req.Message
	}

	incident, err := s.Update(id, status, description, message)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "incident not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incident)
}

func handleResolveIncident(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	incident, err := s.Resolve(id, req.Message)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incident)
}

func handleDeleteIncident(w http.ResponseWriter, r *http.Request, id string) {
	err := s.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "incident not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"deleted": id})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "incidents",
	})
}
