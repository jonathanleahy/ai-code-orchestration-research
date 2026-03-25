package store

import (
	"sync"
	"time"
)

// Component represents a service/component being monitored
type Component struct {
	ID       string
	Name     string
	Category string
	Status   string // "operational", "degraded", "outage"
	Uptime   float64
	UptimeDays []DayStatus
}

// DayStatus represents the status for a single day in the uptime chart
type DayStatus struct {
	Date   string // YYYY-MM-DD
	Status string // "up", "down", "no-data"
}

// Incident represents a past incident
type Incident struct {
	ID      string
	Title   string
	Status  string // "investigating", "identified", "resolved"
	Created time.Time
	Updates []IncidentUpdate
}

// IncidentUpdate represents a timeline update within an incident
type IncidentUpdate struct {
	Message   string
	Status    string
	Timestamp time.Time
}

// Store holds all status page data
type Store struct {
	mu         sync.RWMutex
	Components map[string]*Component
	Incidents  map[string]*Incident
}

// New creates a new Store with sample data
func New() *Store {
	s := &Store{
		Components: make(map[string]*Component),
		Incidents:  make(map[string]*Incident),
	}

	// Sample components
	s.Components["api"] = &Component{
		ID:       "api",
		Name:     "API",
		Category: "Core Services",
		Status:   "operational",
		Uptime:   99.98,
		UptimeDays: generateUptimeDays(90, 0.9985),
	}

	s.Components["web"] = &Component{
		ID:       "web",
		Name:     "Website",
		Category: "Core Services",
		Status:   "operational",
		Uptime:   99.95,
		UptimeDays: generateUptimeDays(90, 0.9995),
	}

	s.Components["db"] = &Component{
		ID:       "db",
		Name:     "Database",
		Category: "Infrastructure",
		Status:   "operational",
		Uptime:   99.99,
		UptimeDays: generateUptimeDays(90, 1.0),
	}

	s.Components["auth"] = &Component{
		ID:       "auth",
		Name:     "Authentication",
		Category: "Core Services",
		Status:   "operational",
		Uptime:   100.0,
		UptimeDays: generateUptimeDays(90, 1.0),
	}

	s.Components["cdn"] = &Component{
		ID:       "cdn",
		Name:     "CDN",
		Category: "Infrastructure",
		Status:   "degraded",
		Uptime:   98.5,
		UptimeDays: generateUptimeDays(90, 0.9850),
	}

	// Sample incidents
	now := time.Now()
	s.Incidents["inc1"] = &Incident{
		ID:      "inc1",
		Title:   "CDN Performance Degradation",
		Status:  "resolved",
		Created: now.Add(-24 * time.Hour),
		Updates: []IncidentUpdate{
			{
				Message:   "Investigating reports of slow CDN responses",
				Status:    "investigating",
				Timestamp: now.Add(-24 * time.Hour),
			},
			{
				Message:   "Root cause identified: traffic spike on edge location",
				Status:    "identified",
				Timestamp: now.Add(-20 * time.Hour),
			},
			{
				Message:   "Issue resolved by scaling edge nodes",
				Status:    "resolved",
				Timestamp: now.Add(-18 * time.Hour),
			},
		},
	}

	s.Incidents["inc2"] = &Incident{
		ID:      "inc2",
		Title:   "Brief Database Replication Delay",
		Status:  "resolved",
		Created: now.Add(-72 * time.Hour),
		Updates: []IncidentUpdate{
			{
				Message:   "Detected replication lag on secondary database",
				Status:    "investigating",
				Timestamp: now.Add(-72 * time.Hour),
			},
			{
				Message:   "Replication resumed, system returned to normal",
				Status:    "resolved",
				Timestamp: now.Add(-71 * time.Hour),
			},
		},
	}

	return s
}

// GetComponents returns all components grouped by category
func (s *Store) GetComponents() map[string][]*Component {
	s.mu.RLock()
	defer s.mu.RUnlock()

	grouped := make(map[string][]*Component)
	for _, comp := range s.Components {
		grouped[comp.Category] = append(grouped[comp.Category], comp)
	}
	return grouped
}

// GetIncidents returns all incidents sorted by creation time
func (s *Store) GetIncidents() []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var incidents []*Incident
	for _, inc := range s.Incidents {
		incidents = append(incidents, inc)
	}

	// Sort by created time descending
	for i := 0; i < len(incidents)-1; i++ {
		for j := i + 1; j < len(incidents); j++ {
			if incidents[j].Created.After(incidents[i].Created) {
				incidents[i], incidents[j] = incidents[j], incidents[i]
			}
		}
	}

	return incidents
}

// GetOverallStatus returns the overall system status
func (s *Store) GetOverallStatus() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hasOutage := false
	hasDegraded := false

	for _, comp := range s.Components {
		if comp.Status == "outage" {
			hasOutage = true
		} else if comp.Status == "degraded" {
			hasDegraded = true
		}
	}

	if hasOutage {
		return "outage"
	}
	if hasDegraded {
		return "degraded"
	}
	return "operational"
}

// generateUptimeDays creates 90 days of uptime data
func generateUptimeDays(days int, upRatio float64) []DayStatus {
	var result []DayStatus
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		// Simple hash-based status for demo
		hash := 0
		for _, c := range dateStr {
			hash = hash*31 + int(c)
		}

		var status string
		if float64(hash%100)/100.0 < upRatio {
			status = "up"
		} else {
			status = "down"
		}

		result = append(result, DayStatus{
			Date:   dateStr,
			Status: status,
		})
	}

	return result
}
