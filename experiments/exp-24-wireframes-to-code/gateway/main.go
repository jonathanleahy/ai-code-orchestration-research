package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Check struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	LastCheck time.Time `json:"last_check"`
}

type Incident struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
}

var checks []Check
var incidents []Incident

func init() {
	checks = []Check{
		{ID: "1", Name: "Service A", URL: "https://example.com/api", Type: "HTTP Status", Status: "operational", LastCheck: time.Now().Add(-2 * time.Minute)},
		{ID: "2", Name: "Service B", URL: "https://example.com/db", Type: "Ping", Status: "down", LastCheck: time.Now().Add(-15 * time.Minute)},
		{ID: "3", Name: "Service C", URL: "https://example.com/cache", Type: "HTTP Status", Status: "degraded", LastCheck: time.Now().Add(-5 * time.Minute)},
	}
	incidents = []Incident{
		{ID: "12345", Title: "Database Connection Issue", Status: "resolved", Severity: "high", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: "12346", Title: "High Latency on API", Status: "open", Severity: "medium", CreatedAt: time.Now().Add(-1 * time.Hour)},
	}
}

func main() {
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/admin/checks/new", addCheckHandler)
	http.HandleFunc("/incidents", incidentsHandler)
	http.HandleFunc("/incidents/new", newIncidentHandler)
	http.HandleFunc("/api/status", apiStatusHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("StatusPulse gateway listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, dashboardHTML())
}

func addCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, addCheckHTML())
}

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, incidentsHTML())
}

func newIncidentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, newIncidentHTML())
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"overall_status": "operational",
		"checks":         checks,
		"incidents":      incidents,
	}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "gateway",
	})
}

func dashboardHTML() string {
	return "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>StatusPulse - Dashboard</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;background:#1a1a1a;color:#ffffff;line-height:1.6}header{background:#0d0d0d;border-bottom:1px solid #333;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.header-left{display:flex;align-items:center;gap:2rem}.logo{font-size:1.5rem;font-weight:bold;color:#4ade80}nav a{color:#999;text-decoration:none;margin-right:1.5rem;transition:color 0.2s}nav a:hover{color:#fff}.header-right{display:flex;gap:1rem}button{background:#4ade80;color:#000;border:none;padding:0.6rem 1.2rem;border-radius:0.4rem;cursor:pointer;font-weight:600;transition:background 0.2s}button:hover{background:#22c55e}main{max-width:1200px;margin:0 auto;padding:2rem}h1{margin-bottom:2rem}.controls{margin-bottom:2rem;display:flex;gap:1rem;justify-content:space-between;align-items:center}.cards-grid{display:grid;grid-template-columns:repeat(auto-fill, minmax(280px, 1fr));gap:1.5rem;margin-bottom:3rem}.card{background:#2a2a2a;border:1px solid #333;border-radius:0.6rem;padding:1.5rem;cursor:pointer;transition:transform 0.2s, border-color 0.2s}.card:hover{transform:translateY(-2px);border-color:#4ade80}.status-indicator{display:inline-block;width:12px;height:12px;border-radius:50%;margin-right:0.5rem}.status-operational{background:#4ade80}.status-degraded{background:#facc15}.status-down{background:#ef4444}.card-header{display:flex;align-items:center;margin-bottom:1rem}.card-name{font-size:1.1rem;font-weight:600}.card-url{color:#999;font-size:0.9rem;margin-bottom:0.5rem}.card-time{color:#666;font-size:0.85rem}.empty-state{text-align:center;padding:2rem;color:#666}.incidents-section{margin-top:3rem}.incidents-section h2{margin-bottom:1rem}.incident-card{background:#2a2a2a;border-left:4px solid #ef4444;padding:1rem;margin-bottom:1rem;border-radius:0.4rem}.incident-title{font-weight:600;margin-bottom:0.5rem}.incident-meta{color:#999;font-size:0.9rem}@media(max-width:768px){header{flex-direction:column;gap:1rem;text-align:center}.header-left{flex-direction:column;width:100%}nav{width:100%}nav a{display:block;margin:0.5rem 0}.header-right{width:100%;flex-direction:column}button{width:100%}.cards-grid{grid-template-columns:1fr}.controls{flex-direction:column}}</style></head><body><header><div class=\"header-left\"><div class=\"logo\">StatusPulse</div><nav><a href=\"/\">Dashboard</a> <a href=\"/incidents\">Incidents</a></nav></div><div class=\"header-right\"><button onclick=\"openAddCheck()\">+ Add Check</button><button onclick=\"openNotifications()\">Notifications</button><button onclick=\"openUser()\">User</button></div></header><main><h1>Dashboard</h1><div class=\"controls\"><span></span></div><div id=\"checksContainer\" class=\"cards-grid\"></div><div id=\"incidentsSection\" class=\"incidents-section\" style=\"display:none;\"><h2>Active Incidents</h2><div id=\"incidentsContainer\"></div></div></main><script>function loadChecks(){fetch('/api/status').then(function(response){return response.json();}).then(function(data){var container=document.getElementById('checksContainer');if(!data.checks||data.checks.length===0){container.innerHTML='<div class=\"empty-state\">No checks configured yet</div>';return;}var html='';data.checks.forEach(function(check){var statusClass='status-operational';if(check.status==='degraded')statusClass='status-degraded';if(check.status==='down')statusClass='status-down';var lastCheck=new Date(check.last_check);var now=new Date();var minutes=Math.floor((now-lastCheck)/60000);var timeStr=minutes+' min ago';html=html+'<div class=\"card\" onclick=\"viewCheck(' + \"'\" + check.id + \"'\" + ')\"><div class=\"card-header\"><span class=\"status-indicator '+statusClass+'\"></span><span class=\"card-name\">'+check.name+'</span></div><div class=\"card-url\">'+check.url+'</div><div class=\"card-time\">'+timeStr+'</div></div>';});container.innerHTML=html;}).catch(function(error){console.error('Error loading checks:',error);});}function loadIncidents(){fetch('/api/status').then(function(response){return response.json();}).then(function(data){if(!data.incidents||data.incidents.length===0){document.getElementById('incidentsSection').style.display='none';return;}document.getElementById('incidentsSection').style.display='block';var container=document.getElementById('incidentsContainer');var html='';data.incidents.forEach(function(incident){var createdAt=new Date(incident.created_at);var now=new Date();var hours=Math.floor((now-createdAt)/3600000);var timeStr=hours>0?hours+' hour ago':'just now';html=html+'<div class=\"incident-card\"><div class=\"incident-title\">'+incident.title+'</div><div class=\"incident-meta\">Status: '+incident.status+' | Severity: '+incident.severity+' | '+timeStr+'</div></div>';});container.innerHTML=html;});}function openAddCheck(){window.location.href='/admin/checks/new';}function openNotifications(){alert('Notifications coming soon');}function openUser(){alert('User menu coming soon');}function viewCheck(id){alert('Check detail view coming soon for check '+id);}loadChecks();loadIncidents();</script></body></html>"
}

func addCheckHTML() string {
	return "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Add Check - StatusPulse</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;background:#1a1a1a;color:#ffffff}header{background:#0d0d0d;border-bottom:1px solid #333;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.logo{font-size:1.5rem;font-weight:bold;color:#4ade80}main{max-width:600px;margin:2rem auto;padding:2rem}.form-container{background:#2a2a2a;border:1px solid #333;border-radius:0.6rem;padding:2rem}h1{margin-bottom:1.5rem}.form-group{margin-bottom:1.5rem}label{display:block;margin-bottom:0.5rem;font-weight:600;color:#ccc}input[type=\"text\"],select{width:100%;padding:0.8rem;background:#1a1a1a;border:1px solid #333;border-radius:0.4rem;color:#ffffff;font-size:1rem}input[type=\"text\"]:focus,select:focus{outline:none;border-color:#4ade80;box-shadow:0 0 0 2px rgba(74,222,128,0.2)}.button-group{display:flex;gap:1rem;margin-top:2rem}button{flex:1;padding:0.8rem;border:none;border-radius:0.4rem;cursor:pointer;font-weight:600;font-size:1rem;transition:background 0.2s}.btn-primary{background:#4ade80;color:#000}.btn-primary:hover{background:#22c55e}.btn-secondary{background:#333;color:#fff}.btn-secondary:hover{background:#444}.error-message{color:#ef4444;font-size:0.9rem;margin-top:0.5rem;display:none}@media(max-width:768px){main{padding:1rem}.form-container{padding:1.5rem}.button-group{flex-direction:column}}</style></head><body><header><div class=\"logo\">StatusPulse</div></header><main><h1>Add New Check</h1><div class=\"form-container\"><form id=\"checkForm\"><div class=\"form-group\"><label for=\"name\">Name</label><input type=\"text\" id=\"name\" name=\"name\" placeholder=\"e.g., API Gateway\" required> <div class=\"error-message\" id=\"nameError\"></div></div><div class=\"form-group\"><label for=\"url\">URL</label><input type=\"text\" id=\"url\" name=\"url\" placeholder=\"https://example.com\" required> <div class=\"error-message\" id=\"urlError\"></div></div><div class=\"form-group\"><label for=\"type\">Check Type</label><select id=\"type\" name=\"type\" required><option value=\"\">Select a check type</option><option value=\"http\">HTTP Status Code</option><option value=\"dns\">DNS Resolution</option><option value=\"ssl\">SSL Certificate</option><option value=\"ping\">Ping</option></select></div><div class=\"button-group\"><button type=\"button\" class=\"btn-secondary\" onclick=\"goBack()\">Cancel</button><button type=\"submit\" class=\"btn-primary\">Create Check</button></div></form></div></main><script>var form=document.getElementById('checkForm');form.addEventListener('submit',function(e){e.preventDefault();var name=document.getElementById('name').value.trim();var url=document.getElementById('url').value.trim();var type=document.getElementById('type').value;document.getElementById('nameError').style.display='none';document.getElementById('urlError').style.display='none';var valid=true;if(!name){document.getElementById('nameError').textContent='Please enter a name';document.getElementById('nameError').style.display='block';valid=false;}if(!url||!isValidUrl(url)){document.getElementById('urlError').textContent='Please enter a valid URL';document.getElementById('urlError').style.display='block';valid=false;}if(!type){valid=false;}if(valid){alert('Check created: '+name);setTimeout(function(){window.location.href='/';},500);}});function isValidUrl(string){try{new URL(string);return true;}catch(_){return false;}}function goBack(){window.location.href='/';}</script></body></html>"
}

func incidentsHTML() string {
	return "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Incidents - StatusPulse</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;background:#1a1a1a;color:#ffffff}header{background:#0d0d0d;border-bottom:1px solid #333;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.header-left{display:flex;align-items:center;gap:2rem}.logo{font-size:1.5rem;font-weight:bold;color:#4ade80}nav a{color:#999;text-decoration:none;margin-right:1.5rem;transition:color 0.2s}nav a:hover{color:#fff}.header-right{display:flex;gap:1rem}button{background:#4ade80;color:#000;border:none;padding:0.6rem 1.2rem;border-radius:0.4rem;cursor:pointer;font-weight:600;transition:background 0.2s}button:hover{background:#22c55e}main{max-width:1000px;margin:0 auto;padding:2rem}h1{margin-bottom:1.5rem}.filters{display:flex;gap:1rem;margin-bottom:2rem}.filter-btn{background:#333;color:#fff;padding:0.6rem 1.2rem}.filter-btn.active{background:#4ade80;color:#000}.incident-card{background:#2a2a2a;border:1px solid #333;border-left:4px solid #ef4444;padding:1.5rem;margin-bottom:1rem;border-radius:0.4rem;cursor:pointer;transition:transform 0.2s, border-color 0.2s}.incident-card:hover{transform:translateX(4px);border-left-color:#facc15}.incident-header{display:flex;justify-content:space-between;align-items:start;margin-bottom:0.5rem}.incident-title{font-weight:600;font-size:1.1rem}.badge{padding:0.3rem 0.8rem;border-radius:0.3rem;font-size:0.85rem;font-weight:600}.badge-high{background:#ef4444;color:#fff}.badge-medium{background:#facc15;color:#000}.badge-low{background:#4ade80;color:#000}.incident-meta{color:#999;font-size:0.9rem;margin-bottom:0.5rem}.incident-status{display:inline-block;padding:0.2rem 0.6rem;background:#333;border-radius:0.2rem;font-size:0.85rem}.empty-state{text-align:center;padding:3rem;color:#666}@media(max-width:768px){header{flex-direction:column;gap:1rem}.header-left{flex-direction:column;width:100%}nav{width:100%}nav a{display:block;margin:0.5rem 0}.filters{flex-wrap:wrap}.incident-header{flex-direction:column}.badge{margin-top:0.5rem}}</style></head><body><header><div class=\"header-left\"><div class=\"logo\">StatusPulse</div><nav><a href=\"/\">Dashboard</a> <a href=\"/incidents\">Incidents</a></nav></div><div class=\"header-right\"><button onclick=\"newIncident()\">+ New Incident</button></div></header><main><h1>Incidents</h1><div class=\"filters\"><button class=\"filter-btn active\" onclick=\"filterIncidents('all')\">All</button><button class=\"filter-btn\" onclick=\"filterIncidents('open')\">Open</button><button class=\"filter-btn\" onclick=\"filterIncidents('investigating')\">Investigating</button><button class=\"filter-btn\" onclick=\"filterIncidents('resolved')\">Resolved</button></div><div id=\"incidentsContainer\"></div></main><script>var currentFilter='all';function loadIncidents(){fetch('/api/status').then(function(response){return response.json();}).then(function(data){var container=document.getElementById('incidentsContainer');if(!data.incidents||data.incidents.length===0){container.innerHTML='<div class=\"empty-state\">No incidents found</div>';return;}var filtered=data.incidents;if(currentFilter!=='all'){filtered=data.incidents.filter(function(incident){return incident.status===currentFilter;});}if(filtered.length===0){container.innerHTML='<div class=\"empty-state\">No incidents with this status</div>';return;}var html='';filtered.forEach(function(incident){var createdAt=new Date(incident.created_at);var now=new Date();var hours=Math.floor((now-createdAt)/3600000);var timeStr=hours>0?hours+' hour ago':'just now';var badgeClass='badge-medium';if(incident.severity==='high')badgeClass='badge-high';if(incident.severity==='low')badgeClass='badge-low';html=html+'<div class=\"incident-card\"><div class=\"incident-header\"><span class=\"incident-title\">'+incident.title+'</span><span class=\"badge '+badgeClass+'\">'+incident.severity.toUpperCase()+'</span></div><div class=\"incident-meta\">'+timeStr+'</div><div><span class=\"incident-status\">'+incident.status+'</span></div></div>';});container.innerHTML=html;}).catch(function(error){console.error('Error loading incidents:',error);});}function filterIncidents(status){currentFilter=status;var buttons=document.querySelectorAll('.filter-btn');buttons.forEach(function(btn){btn.classList.remove('active');});event.target.classList.add('active');loadIncidents();}function newIncident(){window.location.href='/incidents/new';}loadIncidents();</script></body></html>"
}

func newIncidentHTML() string {
	return "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>New Incident - StatusPulse</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;background:#1a1a1a;color:#ffffff}header{background:#0d0d0d;border-bottom:1px solid #333;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.logo{font-size:1.5rem;font-weight:bold;color:#4ade80}main{max-width:600px;margin:2rem auto;padding:2rem}.form-container{background:#2a2a2a;border:1px solid #333;border-radius:0.6rem;padding:2rem}h1{margin-bottom:1.5rem}.form-group{margin-bottom:1.5rem}label{display:block;margin-bottom:0.5rem;font-weight:600;color:#ccc}input[type=\"text\"],textarea,select{width:100%;padding:0.8rem;background:#1a1a1a;border:1px solid #333;border-radius:0.4rem;color:#ffffff;font-size:1rem;font-family:inherit}input[type=\"text\"]:focus,textarea:focus,select:focus{outline:none;border-color:#4ade80;box-shadow:0 0 0 2px rgba(74,222,128,0.2)}textarea{resize:vertical;min-height:150px}.button-group{display:flex;gap:1rem;margin-top:2rem}button{flex:1;padding:0.8rem;border:none;border-radius:0.4rem;cursor:pointer;font-weight:600;font-size:1rem;transition:background 0.2s}.btn-primary{background:#4ade80;color:#000}.btn-primary:hover{background:#22c55e}.btn-secondary{background:#333;color:#fff}.btn-secondary:hover{background:#444}@media(max-width:768px){main{padding:1rem}.form-container{padding:1.5rem}.button-group{flex-direction:column}}</style></head><body><header><div class=\"logo\">StatusPulse</div></header><main><h1>Create New Incident</h1><div class=\"form-container\"><form id=\"incidentForm\"><div class=\"form-group\"><label for=\"title\">Title</label><input type=\"text\" id=\"title\" name=\"title\" placeholder=\"e.g., Database Connection Issue\" required></div><div class=\"form-group\"><label for=\"description\">Description</label><textarea id=\"description\" name=\"description\" placeholder=\"Describe the incident...\" required></textarea></div><div class=\"form-group\"><label for=\"severity\">Severity</label><select id=\"severity\" name=\"severity\" required><option value=\"\">Select severity</option><option value=\"low\">Low</option><option value=\"medium\">Medium</option><option value=\"high\">High</option></select></div><div class=\"button-group\"><button type=\"button\" class=\"btn-secondary\" onclick=\"goBack()\">Cancel</button><button type=\"submit\" class=\"btn-primary\">Create Incident</button></div></form></div></main><script>var form=document.getElementById('incidentForm');form.addEventListener('submit',function(e){e.preventDefault();var title=document.getElementById('title').value.trim();var description=document.getElementById('description').value.trim();var severity=document.getElementById('severity').value;if(title&&description&&severity){alert('Incident created: '+title);setTimeout(function(){window.location.href='/incidents';},500);}});function goBack(){window.location.href='/';}</script></body></html>"
}
