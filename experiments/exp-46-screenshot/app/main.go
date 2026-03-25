package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"statuspage/store"
	"strings"
	"time"
)

var s *store.Store

func init() {
	s = store.New()
}

func main() {
	http.HandleFunc("/", handleStatus)
	http.HandleFunc("/api/status", handleAPIStatus)
	http.HandleFunc("/health", handleHealth)

	log.Println("Starting status page server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	overallStatus := s.GetOverallStatus()
	components := s.GetComponents()
	incidents := s.GetIncidents()

	var allComponents []*store.Component
	for _, categoryComps := range components {
		allComponents = append(allComponents, categoryComps...)
	}

	response := map[string]interface{}{
		"status":     overallStatus,
		"components": allComponents,
		"incidents":  incidents,
		"timestamp":  time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	overallStatus := s.GetOverallStatus()
	components := s.GetComponents()
	incidents := s.GetIncidents()

	// Sort categories
	var categories []string
	for cat := range components {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// Determine banner color and message
	var bannerColor, bannerIcon, statusMsg string
	if overallStatus == "operational" {
		bannerColor = "#10b981"
		bannerIcon = "✓"
		statusMsg = "All Systems Operational"
	} else if overallStatus == "degraded" {
		bannerColor = "#f59e0b"
		bannerIcon = "!"
		statusMsg = "Partial System Outage"
	} else {
		bannerColor = "#ef4444"
		bannerIcon = "!"
		statusMsg = "Major System Outage"
	}

	html := "<!DOCTYPE html>" +
		"<html lang=\"en\">" +
		"<head>" +
		"<meta charset=\"UTF-8\">" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">" +
		"<title>Status Page</title>" +
		"<style>" +
		"* { margin: 0; padding: 0; box-sizing: border-box; }" +
		"body { font-family: system-ui, -apple-system, sans-serif; background: #ffffff; color: #1f2937; }" +
		"header { position: fixed; top: 0; left: 0; right: 0; background: #ffffff; border-bottom: 1px solid #e5e7eb; padding: 16px 20px; z-index: 100; }" +
		"header nav { display: flex; justify-content: space-between; align-items: center; max-width: 800px; margin: 0 auto; }" +
		".logo { font-weight: 700; font-size: 20px; color: #1f2937; }" +
		".subscribe-btn { background: #3b82f6; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-size: 14px; }" +
		".subscribe-btn:hover { background: #2563eb; }" +
		"main { margin-top: 60px; padding: 32px 20px 40px; }" +
		".container { max-width: 800px; margin: 0 auto; }" +
		".status-banner { background: " + bannerColor + "; color: white; padding: 16px 20px; border-radius: 8px; margin-bottom: 32px; display: flex; align-items: center; gap: 12px; font-size: 18px; font-weight: 600; }" +
		".status-icon { font-size: 24px; }" +
		".section { margin-bottom: 32px; }" +
		".section-title { font-size: 16px; font-weight: 700; color: #374151; margin-bottom: 16px; text-transform: uppercase; letter-spacing: 0.05em; }" +
		".component-row { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid #e5e7eb; }" +
		".component-row:last-child { border-bottom: none; }" +
		".component-name { font-size: 14px; color: #1f2937; font-weight: 500; }" +
		".status-badge { display: inline-block; padding: 4px 12px; border-radius: 16px; font-size: 12px; font-weight: 600; }" +
		".status-operational { background: #dbeafe; color: #1e40af; }" +
		".status-degraded { background: #fef3c7; color: #92400e; }" +
		".status-outage { background: #fee2e2; color: #991b1b; }" +
		".uptime-chart { display: flex; gap: 2px; margin-top: 8px; height: 24px; }" +
		".uptime-bar { flex: 1; height: 100%; border-radius: 2px; background: #e5e7eb; cursor: pointer; }" +
		".uptime-bar.up { background: #10b981; }" +
		".uptime-bar.down { background: #ef4444; }" +
		".uptime-percent { font-size: 12px; color: #6b7280; margin-top: 4px; }" +
		".incidents { background: #f9fafb; padding: 20px; border-radius: 8px; }" +
		".incident { margin-bottom: 20px; }" +
		".incident:last-child { margin-bottom: 0; }" +
		".incident-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px; }" +
		".incident-title { font-size: 14px; font-weight: 600; color: #1f2937; }" +
		".incident-date { font-size: 12px; color: #6b7280; }" +
		".incident-status { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; text-transform: uppercase; }" +
		".incident-status.investigating { background: #dbeafe; color: #1e40af; }" +
		".incident-status.identified { background: #fef3c7; color: #92400e; }" +
		".incident-status.resolved { background: #dbeafe; color: #1e40af; }" +
		".incident-update { font-size: 13px; color: #4b5563; margin-left: 16px; padding-left: 12px; border-left: 2px solid #e5e7eb; padding-top: 8px; padding-bottom: 8px; }" +
		".incident-update-time { font-size: 11px; color: #9ca3af; display: block; margin-top: 2px; }" +
		"footer { background: #f3f4f6; padding: 20px; text-align: center; color: #6b7280; font-size: 12px; margin-top: 32px; }" +
		"footer a { color: #3b82f6; text-decoration: none; }" +
		"footer a:hover { text-decoration: underline; }" +
		"</style>" +
		"</head>" +
		"<body>" +
		"<header>" +
		"<nav>" +
		"<div class=\"logo\">Status</div>" +
		"<button class=\"subscribe-btn\">Subscribe to updates</button>" +
		"</nav>" +
		"</header>" +
		"<main>" +
		"<div class=\"container\">" +
		"<div class=\"status-banner\">" +
		"<span class=\"status-icon\">" + bannerIcon + "</span>" +
		"<span>" + statusMsg + "</span>" +
		"</div>"

	// Render components grouped by category
	for _, category := range categories {
		html += "<div class=\"section\">" +
			"<div class=\"section-title\">" + escapeHTML(category) + "</div>"

		for _, comp := range components[category] {
			badgeClass := "status-operational"
			if comp.Status == "degraded" {
				badgeClass = "status-degraded"
			} else if comp.Status == "outage" {
				badgeClass = "status-outage"
			}

			statusText := strings.Title(comp.Status)

			html += "<div class=\"component-row\">" +
				"<div>" +
				"<div class=\"component-name\">" + escapeHTML(comp.Name) + "</div>" +
				"<div class=\"uptime-chart\" title=\"90-day uptime\">"

			// Render uptime bars
			for _, day := range comp.UptimeDays {
				barClass := "uptime-bar "
				if day.Status == "up" {
					barClass += "up"
				} else {
					barClass += "down"
				}
				html += "<div class=\"" + barClass + "\" title=\"" + escapeHTML(day.Date) + "\"></div>"
			}

			html += "</div>" +
				"<div class=\"uptime-percent\">" + fmt.Sprintf("%.2f%% uptime", comp.Uptime) + "</div>" +
				"</div>" +
				"<span class=\"status-badge " + badgeClass + "\">" + statusText + "</span>" +
				"</div>"
		}

		html += "</div>"
	}

	// Render incidents
	if len(incidents) > 0 {
		html += "<div class=\"section\">" +
			"<div class=\"section-title\">Past Incidents</div>" +
			"<div class=\"incidents\">"

		for _, incident := range incidents {
			statusClass := "investigating"
			if incident.Status == "identified" {
				statusClass = "identified"
			} else if incident.Status == "resolved" {
				statusClass = "resolved"
			}

			html += "<div class=\"incident\">" +
				"<div class=\"incident-header\">" +
				"<div>" +
				"<div class=\"incident-title\">" + escapeHTML(incident.Title) + "</div>"

			if len(incident.Updates) > 0 {
				lastUpdate := incident.Updates[len(incident.Updates)-1]
				html += "<div class=\"incident-date\">" + formatDate(lastUpdate.Timestamp) + "</div>"
			}

			html += "<span class=\"incident-status " + statusClass + "\">" + strings.Title(incident.Status) + "</span>" +
				"</div>" +
				"</div>"

			// Render timeline updates
			for _, update := range incident.Updates {
				html += "<div class=\"incident-update\">" +
					escapeHTML(update.Message) +
					"<span class=\"incident-update-time\">" + formatTime(update.Timestamp) + "</span>" +
					"</div>"
			}

			html += "</div>"
		}

		html += "</div>" +
			"</div>"
	}

	html += "</div>" +
		"</main>" +
		"<footer>" +
		"Powered by <a href=\"#\">Status Page</a>" +
		"</footer>" +
		"</body>" +
		"</html>"

	fmt.Fprint(w, html)
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func formatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

func formatTime(t time.Time) string {
	return t.Format("3:04 PM MST")
}
