```go
// 1. Exact Go types

type Check struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Status      string `json:"status"` // "operational", "degraded", "major-outage"
    Category    string `json:"category"`
    Uptime      string `json:"uptime"` // "99.98%"
    UptimeData  []UptimeDay `json:"uptime_data"`
}

type UptimeDay struct {
    Date   string `json:"date"`
    Status string `json:"status"` // "up", "down", "no-data"
}

type Incident struct {
    ID          string         `json:"id"`
    Title       string         `json:"title"`
    Status      string         `json:"status"` // "investigating", "identified", "resolved"
    Timestamp   string         `json:"timestamp"`
    Timeline    []TimelineEntry `json:"timeline"`
}

type TimelineEntry struct {
    ID        string `json:"id"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
    Type      string `json:"type"` // "update", "resolution"
}

// 2. API endpoints

GET /api/status
GET /api/components
GET /api/incidents
GET /api/component/{id}/uptime
GET /api/component/{id}/details

// 3. HTML layout description

/*
<div class="status-page">
  <nav class="top-nav">
    <div class="logo">[logo]</div>
    <button class="subscribe-btn">Subscribe to updates</button>
  </nav>

  <div class="status-banner" data-status="operational">
    <span class="status-icon">✓</span>
    <span class="status-text">All Systems Operational</span>
  </div>

  <div class="main-content">
    <div class="components-section">
      <div class="category-group" data-category="API">
        <h3>API</h3>
        <div class="component-row">
          <span class="component-name">API Gateway</span>
          <span class="status-badge operational">Operational</span>
        </div>
        <!-- more component rows -->
      </div>
      <!-- more category groups -->
    </div>

    <div class="incidents-section">
      <h2>Past Incidents</h2>
      <div class="incident">
        <div class="incident-header">
          <h3>Database Connection Issues</h3>
          <span class="incident-status">Resolved</span>
          <span class="incident-timestamp">2023-05-15 14:30 UTC</span>
        </div>
        <div class="timeline">
          <div class="timeline-entry">
            <span class="timeline-message">Initial report received</span>
            <span class="timeline-timestamp">2023-05-15 14:30 UTC</span>
          </div>
          <!-- more timeline entries -->
        </div>
      </div>
      <!-- more incidents -->
    </div>
  </div>

  <footer class="footer">
    <a href="https://instatus.com">Powered by Instatus</a>
  </footer>
</div>
*/

// 4. Color values
// green = #10b981
// yellow = #f59e0b
// red = #ef4444

// 5. CSS details

/*
:root {
  --primary-green: #10b981;
  --primary-yellow: #f59e0b;
  --primary-red: #ef4444;
  --border-color: #e5e7eb;
  --text-color: #1f2937;
  --bg-white: #ffffff;
  --bg-grey: #f9fafb;
  --footer-bg: #f3f4f6;
  --border-radius: 8px;
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  color: var(--text-color);
  background-color: var(--bg-white);
  margin: 0;
  padding: 0;
}

.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.subscribe-btn {
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 9999px;
  padding: var(--spacing-sm) var(--spacing-md);
  font-weight: 500;
  cursor: pointer;
}

.status-banner {
  padding: 16px;
  border-radius: var(--border-radius);
  margin: var(--spacing-lg) auto;
  max-width: 800px;
  text-align: center;
  font-weight: 600;
}

.status-banner[data-status="operational"] {
  background-color: rgba(16, 185, 129, 0.1);
  color: var(--primary-green);
}

.status-banner[data-status="degraded"] {
  background-color: rgba(245, 158, 11, 0.1);
  color: var(--primary-yellow);
}

.status-banner[data-status="outage"] {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--primary-red);
}

.main-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
}

.components-section {
  margin-bottom: var(--spacing-xl);
}

.category-group {
  margin-bottom: var(--spacing-lg);
}

.category-group h3 {
  margin-top: 0;
  margin-bottom: var(--spacing-md);
  font-size: 1.25rem;
  font-weight: 600;
}

.component-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--border-color);
}

.component-row:last-child {
  border-bottom: none;
}

.component-name {
  font-weight: 500;
}

.status-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius);
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge.operational {
  background-color: rgba(16, 185, 129, 0.1);
  color: var(--primary-green);
}

.status-badge.degraded {
  background-color: rgba(245, 158, 11, 0.1);
  color: var(--primary-yellow);
}

.status-badge.major-outage {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--primary-red);
}

.incidents-section h2 {
  margin-top: 0;
  margin-bottom: var(--spacing-lg);
}

.incident {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
}

.incident-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.incident-status {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius);
  font-size: 0.875rem;
  font-weight: 500;
}

.incident-status.resolved {
  background-color: rgba(16, 185, 129, 0.1);
  color: var(--primary-green);
}

.incident-status.investigating {
  background-color: rgba(245, 158, 11, 0.1);
  color: var(--primary-yellow);
}

.incident-status.identified {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--primary-red);
}

.timeline {
  margin-left: var(--spacing-lg);
}

.timeline-entry {
  margin-bottom: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  border-left: 2px solid var(--border-color);
  padding-left: var(--spacing-md);
}

.footer {
  background-color: var(--footer-bg);
  padding: var(--spacing-md) var(--spacing-lg);
  text-align: center;
  border-top: 1px solid var(--border-color);
}
*/
```