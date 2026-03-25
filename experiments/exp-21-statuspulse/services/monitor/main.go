package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"statuspulse/monitor/store"
)

var s *store.Store

func init() {
	s = store.NewStore()
}

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/api/checks", handleChecks)
	http.HandleFunc("/api/checks/", handleCheckDetail)

	log.Println("Starting monitor server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "monitor",
	})
}

func handleChecks(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		handleCreateCheck(w, r)
	case "GET":
		handleListChecks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCreateCheck(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		IntervalSeconds int    `json:"interval_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	check, err := s.CreateCheck(req.Name, req.URL, req.IntervalSeconds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(check)
}

func handleListChecks(w http.ResponseWriter, r *http.Request) {
	checks := s.ListChecks()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"checks": checks,
	})
}

func handleCheckDetail(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/checks/")

	// Handle ping endpoint
	if strings.HasSuffix(id, "/ping") {
		id = strings.TrimSuffix(id, "/ping")
		if r.Method == "POST" {
			handlePing(w, r, id)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		handleGetCheck(w, r, id)
	case "DELETE":
		handleDeleteCheck(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGetCheck(w http.ResponseWriter, r *http.Request, id string) {
	check, err := s.GetCheck(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "check not found"})
		return
	}

	results := s.GetResults(id, 10)

	response := struct {
		*store.Check
		Results []*store.Result `json:"results"`
	}{
		Check:   check,
		Results: results,
	}

	json.NewEncoder(w).Encode(response)
}

func handleDeleteCheck(w http.ResponseWriter, r *http.Request, id string) {
	err := s.DeleteCheck(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "check not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func handlePing(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	check, err := s.GetCheck(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "check not found"})
		return
	}

	start := time.Now()
	resp, err := http.Get(check.URL)
	latency := time.Since(start).Milliseconds()

	var result *store.Result
	if err != nil {
		result = &store.Result{
			Status:     store.StatusDown,
			LatencyMs:  int(latency),
			StatusCode: 0,
			CheckedAt:  time.Now().Format(time.RFC3339),
		}
	} else {
		defer resp.Body.Close()
		io.ReadAll(resp.Body)

		status := store.StatusUp
		if resp.StatusCode >= 400 {
			status = store.StatusDown
		}

		result = &store.Result{
			Status:     status,
			LatencyMs:  int(latency),
			StatusCode: resp.StatusCode,
			CheckedAt:  time.Now().Format(time.RFC3339),
		}
	}

	if err := s.RecordResult(id, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
