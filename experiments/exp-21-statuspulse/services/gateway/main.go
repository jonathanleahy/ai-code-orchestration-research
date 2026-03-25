package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Response structures
type Check struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ChecksResponse struct {
	Checks []Check `json:"checks"`
}

type Incident struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
}

type IncidentsResponse struct {
	Incidents []Incident `json:"incidents"`
}

type StatusResponse struct {
	OverallStatus  string      `json:"overall_status"`
	Checks         []Check     `json:"checks"`
	OpenIncidents  []Incident  `json:"open_incidents"`
	LastUpdated    time.Time   `json:"last_updated"`
}

// Middleware to add CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Handler for /api/status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch checks from service 1
	checksResp := ChecksResponse{Checks: []Check{}}
	checkReq, _ := http.NewRequest("GET", "http://localhost:8081/api/checks", nil)
	if resp, err := http.DefaultClient.Do(checkReq); err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&checksResp)
	}

	// Fetch incidents from service 2
	incidentsResp := IncidentsResponse{Incidents: []Incident{}}
	incidentReq, _ := http.NewRequest("GET", "http://localhost:8082/api/incidents?open=true", nil)
	if resp, err := http.DefaultClient.Do(incidentReq); err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&incidentsResp)
	}

	// Determine overall status
	overallStatus := "operational"
	if len(incidentsResp.Incidents) > 0 {
		// Check for critical severity
		hasCritical := false
		for _, incident := range incidentsResp.Incidents {
			if incident.Severity == "critical" {
				hasCritical = true
				break
			}
		}
		if hasCritical {
			overallStatus = "major_outage"
		} else {
			overallStatus = "degraded"
		}
	}

	// Build response
	response := StatusResponse{
		OverallStatus: overallStatus,
		Checks:        checksResp.Checks,
		OpenIncidents: incidentsResp.Incidents,
		LastUpdated:   time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// Handler for /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "gateway",
	})
}

// Handler for / - serve HTML status page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlPage))
}

// Embedded HTML status page
const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>System Status</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			background: #0f1419;
			color: #e0e0e0;
			line-height: 1.6;
		}

		.container {
			max-width: 1200px;
			margin: 0 auto;
			padding: 20px;
		}

		header {
			text-align: center;
			padding: 30px 0;
			border-bottom: 1px solid #2a2f36;
			margin-bottom: 30px;
		}

		h1 {
			font-size: 32px;
			margin-bottom: 10px;
			transition: color 0.3s ease;
		}

		h1.operational {
			color: #4ade80;
		}

		h1.degraded {
			color: #facc15;
		}

		h1.major_outage {
			color: #ef4444;
		}

		.status-subtitle {
			color: #888;
			font-size: 14px;
		}

		.main-content {
			display: grid;
			grid-template-columns: 1fr 300px;
			gap: 30px;
		}

		.checks-section h2 {
			font-size: 20px;
			margin-bottom: 20px;
			color: #fff;
		}

		.checks-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
			gap: 15px;
		}

		.check-card {
			background: #1a1f26;
			border: 1px solid #2a2f36;
			border-radius: 8px;
			padding: 20px;
			transition: all 0.3s ease;
		}

		.check-card:hover {
			border-color: #3a3f46;
			transform: translateY(-2px);
		}

		.check-card-header {
			display: flex;
			align-items: center;
			gap: 10px;
			margin-bottom: 15px;
		}

		.status-dot {
			width: 12px;
			height: 12px;
			border-radius: 50%;
			flex-shrink: 0;
		}

		.status-dot.up {
			background: #4ade80;
			box-shadow: 0 0 8px rgba(74, 222, 128, 0.5);
		}

		.status-dot.down {
			background: #ef4444;
			box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
		}

		.status-dot.unknown {
			background: #eab308;
			box-shadow: 0 0 8px rgba(234, 179, 8, 0.5);
		}

		.check-name {
			font-weight: 600;
			color: #fff;
		}

		.check-status {
			font-size: 13px;
			color: #888;
			text-transform: capitalize;
		}

		.sidebar {
			display: flex;
			flex-direction: column;
		}

		.incidents-panel {
			background: #1a1f26;
			border: 1px solid #2a2f36;
			border-radius: 8px;
			padding: 20px;
		}

		.incidents-panel h3 {
			font-size: 16px;
			margin-bottom: 15px;
			color: #fff;
		}

		.incidents-list {
			display: flex;
			flex-direction: column;
			gap: 12px;
		}

		.incident-item {
			padding: 12px;
			background: #0f1419;
			border-left: 3px solid #888;
			border-radius: 4px;
			font-size: 13px;
		}

		.incident-item.critical {
			border-left-color: #ef4444;
			background: rgba(239, 68, 68, 0.1);
		}

		.incident-item.minor {
			border-left-color: #facc15;
			background: rgba(250, 204, 21, 0.1);
		}

		.incident-title {
			font-weight: 600;
			color: #fff;
			margin-bottom: 4px;
		}

		.incident-severity {
			display: inline-block;
			padding: 2px 6px;
			border-radius: 3px;
			font-size: 11px;
			font-weight: 600;
			text-transform: uppercase;
			margin-top: 6px;
		}

		.incident-severity.critical {
			background: rgba(239, 68, 68, 0.3);
			color: #fca5a5;
		}

		.incident-severity.minor {
			background: rgba(250, 204, 21, 0.3);
			color: #fef08a;
		}

		.empty-state {
			text-align: center;
			padding: 30px 20px;
			color: #666;
		}

		.empty-state svg {
			width: 48px;
			height: 48px;
			margin-bottom: 15px;
			opacity: 0.5;
		}

		.last-updated {
			text-align: center;
			color: #666;
			font-size: 12px;
			margin-top: 20px;
			padding-top: 20px;
			border-top: 1px solid #2a2f36;
		}

		@media (max-width: 768px) {
			.main-content {
				grid-template-columns: 1fr;
			}

			.checks-grid {
				grid-template-columns: 1fr;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1 id="statusTitle" class="operational">All Systems Operational</h1>
			<p class="status-subtitle">API Gateway Status</p>
		</header>

		<div class="main-content">
			<div class="checks-section">
				<h2>Service Checks</h2>
				<div class="checks-grid" id="checksContainer">
					<div class="empty-state">Loading...</div>
				</div>
			</div>

			<div class="sidebar">
				<div class="incidents-panel">
					<h3>Open Incidents</h3>
					<div class="incidents-list" id="incidentsContainer">
						<div class="empty-state" style="padding: 20px 0;">No incidents</div>
					</div>
				</div>
				<div class="last-updated" id="lastUpdated">Last updated: just now</div>
			</div>
		</div>
	</div>

	<script>
		var refreshInterval = 30000;
		var lastUpdateTime = new Date();

		function fetchStatus() {
			fetch('/api/status')
				.then(function(response) { return response.json(); })
				.then(function(data) {
					updatePage(data);
					lastUpdateTime = new Date();
				})
				.catch(function(error) {
					console.error('Error fetching status:', error);
				});
		}

		function updatePage(data) {
			// Update header
			var statusTitle = document.getElementById('statusTitle');
			var statusText = data.overall_status === 'operational'
				? 'All Systems Operational'
				: 'Degraded Performance';
			statusTitle.textContent = statusText;
			statusTitle.className = 'operational';
			if (data.overall_status === 'degraded') {
				statusTitle.className = 'degraded';
			} else if (data.overall_status === 'major_outage') {
				statusTitle.className = 'major_outage';
			}

			// Update checks
			var checksContainer = document.getElementById('checksContainer');
			if (data.checks && data.checks.length > 0) {
				checksContainer.innerHTML = '';
				data.checks.forEach(function(check) {
					var statusClass = check.status === 'up' ? 'up' : check.status === 'down' ? 'down' : 'unknown';
					var statusLabel = check.status.charAt(0).toUpperCase() + check.status.slice(1);
					var card = document.createElement('div');
					card.className = 'check-card';
					card.innerHTML = '<div class="check-card-header">' +
						'<div class="status-dot ' + statusClass + '"></div>' +
						'<div>' +
							'<div class="check-name">' + escapeHtml(check.name) + '</div>' +
							'<div class="check-status">' + statusLabel + '</div>' +
						'</div>' +
						'</div>';
					checksContainer.appendChild(card);
				});
			} else {
				checksContainer.innerHTML = '<div class="empty-state">No checks available</div>';
			}

			// Update incidents
			var incidentsContainer = document.getElementById('incidentsContainer');
			if (data.open_incidents && data.open_incidents.length > 0) {
				incidentsContainer.innerHTML = '';
				data.open_incidents.forEach(function(incident) {
					var item = document.createElement('div');
					item.className = 'incident-item ' + incident.severity;
					item.innerHTML = '<div class="incident-title">' + escapeHtml(incident.title) + '</div>' +
						'<div class="incident-severity ' + incident.severity + '">' +
							incident.severity.charAt(0).toUpperCase() + incident.severity.slice(1) +
						'</div>';
					incidentsContainer.appendChild(item);
				});
			} else {
				incidentsContainer.innerHTML = '<div class="empty-state" style="padding: 20px 0;">No incidents</div>';
			}

			// Update last updated time
			updateLastUpdated();
		}

		function updateLastUpdated() {
			var now = new Date();
			var seconds = Math.floor((now - lastUpdateTime) / 1000);
			var timeStr = 'just now';

			if (seconds < 60) {
				timeStr = seconds + 's ago';
			} else if (seconds < 3600) {
				timeStr = Math.floor(seconds / 60) + 'm ago';
			} else {
				timeStr = Math.floor(seconds / 3600) + 'h ago';
			}

			document.getElementById('lastUpdated').textContent = 'Last updated: ' + timeStr;
		}

		function escapeHtml(text) {
			var map = {
				'&': '&amp;',
				'<': '&lt;',
				'>': '&gt;',
				'"': '&quot;',
				"'": '&#039;'
			};
			return text.replace(/[&<>"']/g, function(m) { return map[m]; });
		}

		// Initial fetch
		fetchStatus();

		// Auto-refresh every 30 seconds
		setInterval(fetchStatus, refreshInterval);

		// Update "last updated" time every second
		setInterval(updateLastUpdated, 1000);
	</script>
</body>
</html>`

func main() {
	http.HandleFunc("/", corsMiddleware(indexHandler))
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/api/status", corsMiddleware(statusHandler))

	fmt.Println("Gateway server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
