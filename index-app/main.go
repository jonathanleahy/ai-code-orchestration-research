package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Data model, categories, and lookup maps are in data.go


var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"prevExp":       prevExp,
	"nextExp":       nextExp,
	"prevClone":     prevClone,
	"nextClone":     nextClone,
	"prevCloneName": prevCloneName,
	"nextCloneName": nextCloneName,
	"add":           func(a, b int) int { return a + b },
	"sub":           func(a, b int) int { return a - b },
	"progressPct":   func(idx, total int) int { if total == 0 { return 0 }; return idx * 100 / total },
	"costBarPct":    func(cost, max float64) int { if max <= 0 { return 0 }; return int(cost / max * 100) },
	"getExp":        getExp,
	"scoreColor":    scoreColor,
	"expSearchJSON": expSearchJSON,
	"sortedScores":  sortedScores,
	"mul":           func(a, b float64) float64 { return a * b },
}).Parse(allTemplates))

const allTemplates = `
{{define "style"}}
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --sidebar-width: 260px;
    --accent: #d4727a;
    --accent-light: #f5e1e3;
    --bg: #ffffff;
    --bg-sidebar: #f5f5f5;
    --bg-sidebar-hover: #ebebeb;
    --text: #1a1a1a;
    --text-muted: #555;
    --text-light: #777;
    --border: #e0e0e0;
    --table-stripe: #fafafa;
    --code-bg: #f0f0f0;
    --max-content: 1100px;
  }

  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    color: var(--text);
    background: var(--bg);
    line-height: 1.6;
  }

  a { color: var(--accent); text-decoration: none; }
  a:hover { text-decoration: underline; }
  code { background: var(--code-bg); padding: 2px 5px; border-radius: 3px; font-size: 0.9em; }

  /* Sidebar */
  .sidebar {
    position: fixed; top: 0; left: 0;
    width: var(--sidebar-width); height: 100vh;
    background: var(--bg-sidebar);
    border-right: 1px solid var(--border);
    overflow-y: auto; z-index: 100;
    padding: 24px 0 40px;
    transition: transform 0.3s ease;
  }
  .sidebar-header {
    padding: 0 20px 20px;
    border-bottom: 1px solid var(--border);
    margin-bottom: 8px;
  }
  .sidebar-header a { color: var(--text-muted); text-decoration: none; }
  .sidebar-header a:hover { color: var(--accent); text-decoration: none; }
  .sidebar-header h2 {
    font-size: 14px; font-weight: 700;
    letter-spacing: 0.05em; text-transform: uppercase;
  }
  .sidebar nav ul { list-style: none; }
  .sidebar nav > ul > li > a {
    display: block; padding: 10px 20px;
    font-size: 14px; font-weight: 600;
    color: var(--text); text-decoration: none;
    border-left: 3px solid transparent;
    transition: all 0.15s ease;
  }
  .sidebar nav > ul > li > a:hover {
    background: var(--bg-sidebar-hover);
    border-left-color: var(--border);
    text-decoration: none;
  }
  .sidebar nav > ul > li > a.active {
    background: var(--accent-light);
    border-left-color: var(--accent);
    color: var(--accent);
  }
  .sidebar nav > ul > li > ul { list-style: none; display: none; }
  .sidebar nav > ul > li.expanded > ul { display: block; }
  .sidebar nav > ul > li > ul > li > a {
    display: block; padding: 5px 20px 5px 36px;
    font-size: 12px; color: var(--text-muted);
    text-decoration: none; transition: color 0.15s ease;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }
  .sidebar nav > ul > li > ul > li > a:hover { color: var(--accent); text-decoration: none; }
  .sidebar nav > ul > li > ul > li > a.active { color: var(--accent); font-weight: 600; }

  /* Hamburger */
  .hamburger {
    display: none; position: fixed; top: 12px; left: 12px; z-index: 200;
    background: var(--bg); border: 1px solid var(--border); border-radius: 6px;
    padding: 8px 10px; cursor: pointer; font-size: 20px; line-height: 1;
  }
  .overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.3); z-index: 90; }

  /* Main */
  .main { margin-left: var(--sidebar-width); padding: 32px 48px 80px; }
  .main-inner { max-width: var(--max-content); margin: 0 auto; }

  /* Breadcrumbs */
  .breadcrumbs { font-size: 13px; color: var(--text-light); margin-bottom: 24px; }
  .breadcrumbs a { color: var(--text-muted); text-decoration: none; }
  .breadcrumbs a:hover { color: var(--accent); }
  .breadcrumbs .sep { margin: 0 6px; }

  /* Hero */
  .hero { margin-bottom: 48px; padding-bottom: 32px; border-bottom: 1px solid var(--border); }
  .hero h1 { font-size: 32px; font-weight: 800; line-height: 1.2; margin-bottom: 8px; }
  .hero .subtitle { font-size: 18px; color: var(--text-muted); font-style: italic; margin-bottom: 24px; }
  .stats { display: flex; gap: 32px; flex-wrap: wrap; }
  .stat { text-align: center; }
  .stat-value { font-size: 28px; font-weight: 800; color: var(--accent); display: block; }
  .stat-label { font-size: 13px; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.04em; }

  /* Section */
  .section { margin-bottom: 56px; }
  .section h2 {
    font-size: 24px; font-weight: 700; margin-bottom: 4px;
    padding-bottom: 8px; border-bottom: 2px solid var(--accent); display: inline-block;
  }
  .section .exp-range { font-size: 14px; color: var(--text-light); margin-bottom: 16px; }
  .section p { margin-bottom: 14px; color: var(--text); }
  .section .finding {
    background: var(--table-stripe); border-left: 3px solid var(--accent);
    padding: 12px 16px; margin: 16px 0; font-size: 14px;
  }
  .section .finding strong { color: var(--accent); }

  /* Tables */
  table { width: 100%; border-collapse: collapse; margin: 20px 0 24px; font-size: 14px; }
  thead th {
    background: var(--bg-sidebar); font-weight: 600; text-align: left;
    padding: 10px 12px; border-bottom: 2px solid var(--border); white-space: nowrap;
  }
  tbody td { padding: 8px 12px; border-bottom: 1px solid var(--border); }
  tbody tr:nth-child(even) { background: var(--table-stripe); }
  tbody tr:hover { background: var(--accent-light); }
  td.num, th.num { text-align: right; font-variant-numeric: tabular-nums; }
  td a { color: var(--accent); text-decoration: none; font-weight: 500; }
  td a:hover { text-decoration: underline; }

  /* Category cards */
  .cat-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 20px; margin: 24px 0; }
  .cat-card {
    border: 1px solid var(--border); border-radius: 8px; padding: 24px;
    transition: all 0.15s ease; text-decoration: none; color: var(--text); display: block;
  }
  .cat-card:hover { border-color: var(--accent); background: var(--accent-light); text-decoration: none; }
  .cat-card h3 { font-size: 18px; font-weight: 700; margin-bottom: 4px; color: var(--text); }
  .cat-card .range { font-size: 13px; color: var(--text-light); margin-bottom: 8px; }
  .cat-card .count { font-size: 14px; color: var(--text-muted); }

  /* Experiment grid on homepage */
  .exp-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 12px; margin: 16px 0 48px; }
  .exp-block {
    border: 1px solid var(--border); border-radius: 8px; padding: 16px;
    transition: all 0.15s ease; text-decoration: none; color: var(--text); display: block;
    background: var(--bg);
  }
  .exp-block:hover { border-color: var(--accent); background: var(--accent-light); transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.08); }
  .exp-block-num { font-size: 11px; font-weight: 700; color: var(--accent); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px; }
  .exp-block-title { font-size: 14px; font-weight: 600; color: var(--text); margin-bottom: 6px; line-height: 1.3; }
  .exp-block-result { font-size: 12px; color: var(--text-muted); margin-bottom: 4px; }
  .exp-block-finding { font-size: 11px; color: var(--text-light); line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }

  /* Score badges */
  .score-badge { display:inline-block; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600; color:#fff; margin-left:6px; }

  /* Search */
  .search-container { position:relative; margin-bottom:24px; }
  .search-input { width:100%; max-width:500px; padding:10px 16px; border:1px solid var(--border); border-radius:6px; font-size:15px; background:var(--bg); color:var(--text); font-family:inherit; }
  .search-input:focus { border-color:var(--accent); outline:none; box-shadow:0 0 0 3px rgba(212,114,122,0.15); }
  .search-results { position:absolute; z-index:50; background:var(--bg); border:1px solid var(--border); border-radius:8px; max-height:400px; overflow-y:auto; width:100%; max-width:500px; box-shadow:0 4px 16px rgba(0,0,0,0.1); margin-top:4px; }
  .search-result-item { display:block; padding:10px 16px; border-bottom:1px solid var(--border); color:var(--text); text-decoration:none; }
  .search-result-item:hover { background:var(--accent-light); }
  .search-result-item .sr-num { font-size:11px; font-weight:700; color:var(--accent); }
  .search-result-item .sr-title { font-size:14px; font-weight:600; }
  .search-result-item .sr-finding { font-size:12px; color:var(--text-light); }
  .search-result-item mark { background:#ffe066; border-radius:2px; padding:0 1px; }

  /* Grid sort */
  .grid-sort { padding:6px 12px; border:1px solid var(--border); border-radius:6px; font-size:13px; background:var(--bg); color:var(--text); cursor:pointer; font-family:inherit; }

  /* Related cards */
  .related-cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(220px, 1fr)); gap:12px; margin-top:12px; }
  .related-card { border:1px solid var(--border); border-radius:8px; padding:14px; text-decoration:none; color:var(--text); transition:all 0.15s; display:block; }
  .related-card:hover { border-color:var(--accent); background:var(--accent-light); transform:translateY(-1px); }
  .related-card-num { font-size:11px; font-weight:700; color:var(--accent); margin-bottom:4px; }
  .related-card-title { font-size:13px; font-weight:600; margin-bottom:4px; line-height:1.3; }
  .related-card-finding { font-size:11px; color:var(--text-light); line-height:1.4; display:-webkit-box; -webkit-line-clamp:2; -webkit-box-orient:vertical; overflow:hidden; }

  .exp-screenshots{margin:24px 0 32px;}
  .exp-screenshot-pair{display:grid;grid-template-columns:1fr 1fr;gap:16px;}
  .exp-screenshot{border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .exp-screenshot img{width:100%;display:block;cursor:zoom-in;}
  .exp-screenshot-label{font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:1px;padding:8px 12px;color:var(--text-muted);border-bottom:1px solid var(--border);background:var(--bg-sidebar);}
  @media(max-width:600px){.exp-screenshot-pair{grid-template-columns:1fr;}}

  /* Thumbnail strip */
  .exp-block-thumb { height:80px; border-radius:6px 6px 0 0; margin:-16px -16px 10px; background-size:cover; background-position:top center; }

  /* Completeness bar */
  .completeness-bar { margin-top:16px; max-width:500px; }
  .completeness-track { height:8px; background:var(--border); border-radius:4px; overflow:hidden; }
  .completeness-fill { height:100%; background:var(--accent); border-radius:4px; transition:width 0.3s; }

  /* Experiment detail */
  .exp-detail { margin-top: 16px; }
  .exp-detail .meta { display: flex; gap: 32px; flex-wrap: wrap; margin: 16px 0 24px; }
  .exp-detail .meta-label { font-size: 12px; text-transform: uppercase; letter-spacing: 0.04em; color: var(--text-light); }
  .exp-detail .meta-value { font-size: 18px; font-weight: 700; color: var(--text); }

  /* Nav arrows */
  .exp-nav { display: flex; justify-content: space-between; margin-top: 48px; padding-top: 24px; border-top: 1px solid var(--border); }
  .exp-nav a { font-size: 14px; }

  /* Experiment sections */
  .exp-section { margin-top: 28px; padding: 20px 24px; background: #fafafa; border-radius: 6px; border-left: 3px solid var(--border); }
  .exp-section h3 { font-size: 14px; text-transform: uppercase; letter-spacing: 1px; color: var(--accent); margin-bottom: 10px; font-weight: 600; }
  .exp-section p { line-height: 1.7; color: var(--text); }
  .exp-section code { background: #e8e8e8; }

  /* Source code viewer */
  .source-section { margin-top: 32px; }
  .source-section h3 { font-size: 16px; font-weight: 600; margin-bottom: 8px; }
  .source-file { font-size: 13px; color: var(--text-muted); margin-bottom: 12px; }
  .code-toggle { background: var(--accent); color: #fff; border: none; padding: 6px 16px; border-radius: 4px; cursor: pointer; font-size: 13px; margin-bottom: 8px; }
  .code-toggle:hover { opacity: 0.9; }
  .source-pre { border-radius: 6px; margin: 0; }
  .source-pre code { font-size: 13px !important; line-height: 1.5 !important; }

  /* Diagram lightbox */
  .mermaid { cursor: pointer; text-align: center; margin: 24px 0; }
  .mermaid svg { margin: 0 auto; display: block; }
  .diagram-overlay { display: none; position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.85); z-index: 9999; justify-content: center; align-items: center; flex-direction: column; padding: 40px; }
  .diagram-overlay.active { display: flex; }
  .diagram-overlay-close { position: absolute; top: 20px; right: 30px; color: #fff; font-size: 32px; cursor: pointer; z-index: 10000; background: none; border: none; }
  .diagram-overlay-close:hover { color: var(--accent); }
  .diagram-overlay-content { background: #fff; border-radius: 8px; padding: 40px; max-width: 95vw; max-height: 90vh; overflow: auto; }
  .diagram-overlay-content .mermaid { cursor: default; }

  /* Footer */
  .footer {
    margin-top: 64px; padding-top: 24px; border-top: 1px solid var(--border);
    font-size: 13px; color: var(--text-light); text-align: center;
  }

  /* Responsive */
  @media (max-width: 900px) {
    .sidebar { transform: translateX(-100%); }
    .sidebar.open { transform: translateX(0); }
    .overlay.open { display: block; }
    .hamburger { display: block; }
    .main { margin-left: 0; padding: 64px 20px 60px; }
  }
  /* Dark mode */
  [data-theme="dark"] {
    --bg: #1a1a2e; --bg-sidebar: #16213e; --bg-sidebar-hover: #1a2744;
    --text: #e0e0e0; --text-muted: #aaa; --text-light: #888;
    --border: #333; --table-stripe: #1e2a3a; --code-bg: #2d2d2d;
    --accent: #e8878e; --accent-light: #2a1f2f;
  }
  [data-theme="dark"] .exp-section { background: #1e2a3a; }
  [data-theme="dark"] .exp-section code { background: #3a3a4a; }

  .exp-result-card{background:linear-gradient(135deg,var(--accent-light),#fff);border:1px solid var(--accent);border-radius:10px;padding:20px 24px;margin-bottom:28px;display:grid;grid-template-columns:1fr auto;gap:16px;align-items:center;}
  [data-theme="dark"] .exp-result-card{background:linear-gradient(135deg,#2a1f2e,#1a1a2e);}
  .exp-result-card .result-label{font-size:12px;text-transform:uppercase;letter-spacing:1px;color:var(--text-muted);margin-bottom:4px;}
  .exp-result-card .result-value{font-size:20px;font-weight:700;color:var(--text);}
  .exp-result-card .result-finding{font-size:14px;color:var(--text-muted);line-height:1.5;margin-top:8px;}
  .exp-result-card .result-meta{text-align:right;}
  .exp-result-card .result-cost{font-size:24px;font-weight:800;color:var(--accent);}
  .exp-result-card .result-time{font-size:12px;color:var(--text-light);}
  .exp-section.section-why{border-left-color:#6366f1;}
  .exp-section.section-what{border-left-color:#0d9488;}
  .exp-section.section-how{border-left-color:#d97706;}
  .exp-section.section-impact{border-left-color:#16a34a;}
  [data-theme="dark"] .exp-section.section-why{border-left-color:#818cf8;}
  [data-theme="dark"] .exp-section.section-what{border-left-color:#2dd4bf;}
  [data-theme="dark"] .exp-section.section-how{border-left-color:#fbbf24;}
  [data-theme="dark"] .exp-section.section-impact{border-left-color:#4ade80;}
  .exp-section.section-why h3{color:#6366f1;}
  .exp-section.section-what h3{color:#0d9488;}
  .exp-section.section-how h3{color:#d97706;}
  .exp-section.section-impact h3{color:#16a34a;}
  .exp-title-area{margin-bottom:24px;}
  .exp-title-area h2{font-size:28px;font-weight:800;line-height:1.2;margin-bottom:8px;}
  .cat-narrative{max-width:800px;}
  .cat-narrative p{margin-bottom:16px;line-height:1.8;font-size:16px;}
  .cat-narrative p:first-child{font-size:18px;}
  .key-insight{background:linear-gradient(135deg,var(--accent-light),#fff);border:2px solid var(--accent);border-radius:12px;padding:24px 28px;margin:32px 0;position:relative;}
  [data-theme="dark"] .key-insight{background:linear-gradient(135deg,#2a1f2e,#1a1a2e);}
  .key-insight::before{content:"KEY INSIGHT";position:absolute;top:-10px;left:20px;background:var(--accent);color:#fff;font-size:10px;font-weight:700;letter-spacing:1px;padding:2px 10px;border-radius:4px;}
  .key-insight p{font-size:16px;line-height:1.7;color:var(--text);margin:0;}
  .cat-header{margin-bottom:32px;}
  .cat-header h2{font-size:28px;font-weight:800;margin-bottom:4px;}

  /* Filter pills */
  .filter-pills { display:flex; gap:8px; flex-wrap:wrap; margin-bottom:16px; }
  .filter-pill { border:1px solid var(--border); background:var(--bg); color:var(--text-muted); padding:4px 14px; border-radius:16px; font-size:12px; cursor:pointer; transition:all 0.2s; }
  .filter-pill:hover { border-color:var(--accent); }
  .filter-pill.active { background:var(--accent); color:#fff; border-color:var(--accent); }

  /* Progress bar */
  .progress-bar { margin-bottom:16px; }
  .progress-label { font-size:12px; color:var(--text-light); margin-bottom:4px; }
  .progress-track { height:3px; background:var(--border); border-radius:2px; }
  .progress-fill { height:100%; background:var(--accent); border-radius:2px; transition:width 0.3s; }

  /* Cost bar */
  .cost-bar { height:6px; background:#eee; border-radius:3px; margin-top:4px; width:100px; display:inline-block; vertical-align:middle; }
  .cost-bar-fill { height:100%; background:var(--accent); border-radius:3px; }
  [data-theme="dark"] .cost-bar { background:#333; }

  /* Journey nav */
  .journey-nav { display:flex; justify-content:space-between; align-items:center; padding:16px 0; border-top:1px solid var(--border); margin-top:32px; }
  .journey-chapter { font-size:14px; color:var(--text-muted); }
  .journey-btn { background:var(--accent); color:#fff; padding:8px 20px; border-radius:4px; text-decoration:none; font-size:14px; }
  .journey-btn:hover { opacity:0.9; }

  /* Dark mode toggle */
  .theme-toggle { background:none; border:1px solid var(--border); border-radius:4px; padding:4px 8px; cursor:pointer; font-size:16px; float:right; margin-top:-2px; }
  .theme-toggle:hover { border-color:var(--accent); }

  @media (max-width: 600px) {
    .hero h1 { font-size: 24px; }
    .hero .subtitle { font-size: 15px; }
    .stats { gap: 20px; }
    .stat-value { font-size: 22px; }
    table { font-size: 12px; }
    thead th, tbody td { padding: 6px 8px; }
    .cat-cards { grid-template-columns: 1fr; }
  }
</style>
{{end}}

{{define "sidebar"}}
<button class="hamburger" onclick="document.getElementById('sidebar').classList.toggle('open');document.getElementById('overlay').classList.toggle('open')" aria-label="Toggle navigation">&#9776;</button>
<div class="overlay" id="overlay" onclick="document.getElementById('sidebar').classList.toggle('open');this.classList.toggle('open')"></div>
<aside class="sidebar" id="sidebar">
  <div class="sidebar-header">
    <a href="/"><h2>Research Index</h2></a>
    <button class="theme-toggle" id="themeToggle" title="Toggle dark mode" aria-label="Toggle dark mode">&#9789;</button>
  </div>
  <nav>
    <ul>
      {{range .Categories}}
      <li{{if eq $.ActiveCat .Slug}} class="expanded"{{end}}>
        <a href="/category/{{.Slug}}"{{if eq $.ActiveCat .Slug}} class="active"{{end}}>{{.Name}}</a>
        <ul>
          {{range .Experiments}}{{if and .HasDetail (gt .NumID 0)}}
          <li><a href="/exp/{{.Num}}"{{if eq $.ActiveExp .NumID}} class="active"{{end}}>{{.Num}}: {{.Focus}}</a></li>
          {{end}}{{end}}
        </ul>
      </li>
      {{end}}
      <li style="margin-top:12px; border-top:1px solid var(--border); padding-top:8px;">
        <a href="/discovery"{{if eq $.PageType "discovery"}} class="active"{{end}}>Discovery Graph</a>
        <a href="/timeline"{{if eq $.PageType "timeline"}} class="active"{{end}}>Timeline</a>
        <a href="/compare/"{{if or (eq $.PageType "compare") (eq $.PageType "compare-detail")}} class="active"{{end}}>Visual Comparison</a>
        <a href="/roadmap"{{if eq $.PageType "roadmap"}} class="active"{{end}}>Roadmap</a>
        <a href="/pipeline/stats"{{if eq $.PageType "pipeline-stats"}} class="active"{{end}}>Pipeline Stats</a>
      </li>
    </ul>
  </nav>
</aside>
{{end}}

{{define "breadcrumbs"}}
<div class="breadcrumbs">
  {{range $i, $bc := .Breadcrumbs}}{{if $i}}<span class="sep">&rsaquo;</span>{{end}}{{if $bc.URL}}<a href="{{$bc.URL}}">{{$bc.Label}}</a>{{else}}<span>{{$bc.Label}}</span>{{end}}{{end}}
</div>
{{end}}

{{define "footer"}}
<div class="footer">
  <p>AI Code Orchestration Research &mdash; 114+ experiments, 11+ models, ~$12 total cost</p>
  <p>Built with Go, served on port 9094</p>
</div>
{{end}}

{{define "layout_start"}}
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/line-numbers/prism-line-numbers.min.css">
{{template "style"}}
</head>
<body>
{{template "sidebar" .}}
<div class="main">
<div class="main-inner">
{{template "breadcrumbs" .}}
{{end}}

{{define "layout_end"}}
{{template "footer"}}
</div></div>
<div class="diagram-overlay" id="diagramOverlay">
  <button class="diagram-overlay-close" id="diagramClose">&times;</button>
  <div class="diagram-overlay-content" id="diagramContent"></div>
</div>
<script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
<script>
mermaid.initialize({startOnLoad:true, securityLevel:'loose', theme:'neutral', themeVariables:{primaryColor:'#f5e1e3', primaryBorderColor:'#d4727a', lineColor:'#999'}});
// Re-render mermaid diagrams on back/forward navigation (bfcache)
window.addEventListener('pageshow', function(e) {
  if (e.persisted) {
    document.querySelectorAll('.mermaid[data-processed]').forEach(function(el) {
      el.removeAttribute('data-processed');
      el.innerHTML = el.textContent;
    });
    mermaid.run();
  }
});
</script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-python.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-bash.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-go.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/line-numbers/prism-line-numbers.min.js"></script>
<script>
// Diagram lightbox — click any mermaid diagram to view full-screen
const overlay = document.getElementById('diagramOverlay');
const overlayContent = document.getElementById('diagramContent');
const overlayClose = document.getElementById('diagramClose');

document.querySelectorAll('.mermaid').forEach(el => {
  // Skip diagrams inside the overlay itself
  if (el.closest('.diagram-overlay')) return;
  el.style.cursor = 'pointer';
  el.title = 'Click to view full size';
  el.addEventListener('click', (e) => {
    // If clicking a link inside the diagram (Mermaid click directives), let it navigate
    if (e.target.closest('a')) return;
    // Otherwise show lightbox
    overlayContent.innerHTML = el.innerHTML;
    overlay.classList.add('active');
  });
});

overlayClose.addEventListener('click', () => overlay.classList.remove('active'));
overlay.addEventListener('click', (e) => {
  if (e.target === overlay) overlay.classList.remove('active');
});
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') overlay.classList.remove('active');
});

// Code toggle
document.querySelectorAll('.code-toggle').forEach(btn => {
  const pre = btn.nextElementSibling;
  if (pre) { pre.style.display = 'none'; }
  btn.addEventListener('click', () => {
    if (pre) {
      const hidden = pre.style.display === 'none';
      pre.style.display = hidden ? 'block' : 'none';
      btn.textContent = hidden ? 'Hide source' : 'Show full source';
    }
  });
});

// Dark mode toggle
(function() {
  var saved = localStorage.getItem('theme');
  if (saved) document.documentElement.setAttribute('data-theme', saved);
  var btn = document.getElementById('themeToggle');
  if (!btn) return;
  function updateBtn() {
    var isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    btn.textContent = isDark ? '\u2600' : '\u263D';
  }
  updateBtn();
  btn.addEventListener('click', function() {
    var isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    var next = isDark ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    updateBtn();
    // Re-init mermaid with appropriate theme
    if (typeof mermaid !== 'undefined') {
      var mTheme = next === 'dark' ? 'dark' : 'neutral';
      mermaid.initialize({startOnLoad:false, securityLevel:'loose', theme:mTheme, themeVariables:{primaryColor:'#f5e1e3', primaryBorderColor:'#d4727a', lineColor:'#999'}});
    }
  });
})();

// Filter pills
(function() {
  document.querySelectorAll('.filter-pills').forEach(function(container) {
    var pills = container.querySelectorAll('.filter-pill');
    var tableId = container.getAttribute('data-table');
    var table = document.getElementById(tableId);
    if (!table) return;
    var rows = table.querySelectorAll('tbody tr');
    pills.forEach(function(pill) {
      pill.addEventListener('click', function() {
        pills.forEach(function(p){ p.classList.remove('active'); });
        pill.classList.add('active');
        var filter = pill.getAttribute('data-filter');
        rows.forEach(function(row) {
          if (filter === 'all') { row.style.display = ''; return; }
          if (filter === 'success') { row.style.display = row.getAttribute('data-status') === 'success' ? '' : 'none'; return; }
          if (filter === 'failure') { row.style.display = row.getAttribute('data-status') === 'failure' ? '' : 'none'; return; }
          if (filter === 'has-detail') { row.style.display = row.getAttribute('data-detail') === 'yes' ? '' : 'none'; return; }
        });
      });
    });
  });
})();

// Search functionality
if (document.getElementById('searchInput')) {
  var searchInput = document.getElementById('searchInput');
  var searchResults = document.getElementById('searchResults');
  var debounceTimer;

  searchInput.addEventListener('input', function() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(function() {
      var q = searchInput.value.toLowerCase().trim();
      if (q.length < 2) { searchResults.style.display = 'none'; return; }
      var matches = SEARCH_DATA.filter(function(e) {
        return e.focus.toLowerCase().indexOf(q) >= 0 || e.finding.toLowerCase().indexOf(q) >= 0 || e.result.toLowerCase().indexOf(q) >= 0 || e.num.indexOf(q) >= 0;
      }).slice(0, 12);
      if (matches.length === 0) { searchResults.style.display = 'none'; return; }
      var html = '';
      matches.forEach(function(m) {
        var title = m.focus.replace(new RegExp('(' + q.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + ')', 'gi'), '<mark>$1</mark>');
        html += '<a class="search-result-item" href="/exp/' + m.num + '">';
        html += '<div class="sr-num">' + (m.icon ? m.icon + ' ' : '') + 'Exp ' + m.num + ' &middot; ' + m.catName + '</div>';
        html += '<div class="sr-title">' + title + '</div>';
        html += '<div class="sr-finding">' + m.finding.substring(0, 100) + '</div></a>';
      });
      searchResults.innerHTML = html;
      searchResults.style.display = 'block';
    }, 150);
  });

  document.addEventListener('click', function(e) {
    if (!e.target.closest('.search-container')) searchResults.style.display = 'none';
  });

  searchInput.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') { searchResults.style.display = 'none'; searchInput.blur(); }
  });
}

// Grid filtering
var gridPills = document.querySelectorAll('#gridFilterPills .filter-pill');
gridPills.forEach(function(pill) {
  pill.addEventListener('click', function() {
    gridPills.forEach(function(p) { p.classList.remove('active'); });
    pill.classList.add('active');
    var cat = pill.getAttribute('data-cat');
    document.querySelectorAll('.exp-grid .exp-block').forEach(function(block) {
      if (cat === 'all' || block.getAttribute('data-category') === cat) {
        block.style.display = '';
      } else {
        block.style.display = 'none';
      }
    });
  });
});

// Grid sorting
var gridSort = document.getElementById('gridSort');
if (gridSort) {
  gridSort.addEventListener('change', function() {
    var grid = document.querySelector('.exp-grid');
    var blocks = Array.from(grid.querySelectorAll('.exp-block'));
    blocks.sort(function(a, b) {
      if (gridSort.value === 'number') return (parseInt(a.dataset.num) || 0) - (parseInt(b.dataset.num) || 0);
      if (gridSort.value === 'cost-asc') return (parseFloat(a.dataset.cost) || 0) - (parseFloat(b.dataset.cost) || 0);
      if (gridSort.value === 'cost-desc') return (parseFloat(b.dataset.cost) || 0) - (parseFloat(a.dataset.cost) || 0);
      return 0;
    });
    blocks.forEach(function(b) { grid.appendChild(b); });
  });
}
</script>
</body>
</html>
{{end}}

{{define "home"}}
{{template "layout_start" .}}
<div class="hero">
  <h1>AI Code Orchestration Research</h1>
  <p class="subtitle">Can AI go from a one-line idea to a running, tested, reviewed product?</p>
  <div class="stats">
    <div class="stat"><span class="stat-value">114+</span><span class="stat-label">Experiments</span></div>
    <div class="stat"><span class="stat-value">10+</span><span class="stat-label">Working Apps</span></div>
    <div class="stat"><span class="stat-value">$0.01&ndash;$2.75</span><span class="stat-label">Cost per Product</span></div>
    <div class="stat"><span class="stat-value">~$12</span><span class="stat-label">Total Research Cost</span></div>
    <div class="stat"><span class="stat-value">11+</span><span class="stat-label">Models Tested</span></div>
    <div class="stat"><span class="stat-value">{{.NarrativeComplete}}/{{.NarrativeTotal}}</span><span class="stat-label">Full Narratives</span></div>
  </div>
  <div class="completeness-bar">
    <div style="font-size:12px; color:var(--text-muted); margin-bottom:4px;">{{.NarrativeComplete}} of {{.NarrativeTotal}} experiments have complete narratives</div>
    <div class="completeness-track"><div class="completeness-fill" style="width:{{progressPct .NarrativeComplete .NarrativeTotal}}%"></div></div>
  </div>
  <div style="margin-top:24px;">
    <a href="/journey/1" style="font-size:14px; color:var(--accent); font-weight:600; text-decoration:none;">Start the story &rarr;</a>
  </div>
  <div class="search-container" style="margin-top:24px;">
    <input type="text" id="searchInput" class="search-input" placeholder="Search 114+ experiments..." autocomplete="off" style="max-width:100%;">
    <div id="searchResults" class="search-results" style="display:none;"></div>
  </div>
  <script>var SEARCH_DATA = {{expSearchJSON}};</script>
</div>
<h2 style="font-size:20px; font-weight:700; margin-bottom:16px;">Research Categories</h2>
<div class="cat-cards">
  {{range .Categories}}
  <a class="cat-card" href="/category/{{.Slug}}">
    <h3>{{if .Icon}}<span style="margin-right:6px;">{{.Icon}}</span>{{end}}{{.Name}}</h3>
    <div class="range">{{.ExpRange}}</div>
    <div class="count">{{len .Experiments}} entries</div>
  </a>
  {{end}}
</div>

<div style="display:flex; justify-content:space-between; align-items:center; flex-wrap:wrap; gap:12px; margin:32px 0 8px;">
  <h2 style="font-size:20px; font-weight:700; margin:0;">All Experiments</h2>
  <a href="/api/experiments.json" download="experiments.json" style="padding:6px 16px; border:1px solid var(--accent); border-radius:6px; color:var(--accent); font-size:13px; font-weight:600; text-decoration:none;">Download JSON</a>
</div>
<p style="margin-bottom:16px; color:var(--text-muted);">Every experiment at a glance. Click any card to read the full narrative.</p>
<div style="display:flex; justify-content:space-between; align-items:center; flex-wrap:wrap; gap:8px; margin-bottom:12px;">
  <div class="filter-pills" id="gridFilterPills">
    <span class="filter-pill active" data-cat="all">All</span>
    {{range .Categories}}<span class="filter-pill" data-cat="{{.Slug}}">{{if .Icon}}{{.Icon}} {{end}}{{.Name}}</span>{{end}}
  </div>
  <select id="gridSort" class="grid-sort">
    <option value="number">Sort: Number</option>
    <option value="cost-asc">Cost (low&#8594;high)</option>
    <option value="cost-desc">Cost (high&#8594;low)</option>
  </select>
</div>
<div class="exp-grid">
  {{range $ci, $cat := .Categories}}
  {{range $cat.Experiments}}{{if and .HasDetail (gt .NumID 0)}}
  <a class="exp-block" href="/exp/{{.Num}}" data-category="{{$cat.Slug}}" data-cost="{{.CostFloat}}" data-num="{{.NumID}}">
    <div class="exp-block-num">{{if .Icon}}<span class="exp-icon">{{.Icon}}</span> {{end}}{{.Num}}</div>
    <div class="exp-block-title">{{.Focus}}</div>
    {{if .Result}}<div class="exp-block-result">{{.Result}}</div>{{end}}
    {{if gt .Score 0.0}}<span class="score-badge" style="background:{{scoreColor .Score}}">{{printf "%.3f" .Score}}</span>{{end}}
    {{if .Finding}}<div class="exp-block-finding">{{.Finding}}</div>{{end}}
  </a>
  {{end}}{{end}}
  {{end}}
</div>

{{template "layout_end" .}}
{{end}}

{{define "category"}}
{{template "layout_start" .}}
<div class="section">
  <div class="cat-header">
    <h2>{{.Category.Name}}</h2>
    <p class="exp-range">{{.Category.ExpRange}}</p>
  </div>
  <div class="cat-narrative">{{range .Category.Narrative}}<p>{{.}}</p>{{end}}</div>
  <div class="key-insight"><p>{{.Category.KeyInsight}}</p></div>

  {{if eq .Category.TableType "standard"}}
  <div class="filter-pills" data-table="standard-table">
    <span class="filter-pill active" data-filter="all">All</span>
    <span class="filter-pill" data-filter="success">Success</span>
    <span class="filter-pill" data-filter="failure">Failure</span>
    <span class="filter-pill" data-filter="has-detail">Has Detail</span>
  </div>
  <table id="standard-table">
    <thead><tr><th>Exp</th><th>Focus</th><th>Result</th><th class="num">Cost</th><th>Key Finding</th></tr></thead>
    <tbody>
      {{range .Category.Experiments}}
      <tr data-status="{{if .IsFailure}}failure{{else}}success{{end}}" data-detail="{{if and .HasDetail (gt .NumID 0)}}yes{{else}}no{{end}}">
        <td>{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Num}}</a>{{else}}{{.Num}}{{end}}</td>
        <td>{{.Focus}}</td><td>{{.Result}}</td>
        <td class="num">{{.Cost}}{{if gt .CostFloat 0.0}} <div class="cost-bar"><div class="cost-bar-fill" style="width:{{costBarPct .CostFloat $.Category.MaxCost}}%"></div></div>{{end}}</td>
        <td>{{.Finding}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>
  {{end}}

  {{if eq .Category.TableType "website-clone"}}
  <table>
    <thead><tr><th>Rank</th><th>Model</th><th class="num">SSIM</th><th>Pass Rate</th><th class="num">Cost</th><th>Visual Quality</th></tr></thead>
    <tbody>
      {{range .Category.Experiments}}
      <tr><td>{{.Num}}</td><td>{{.Focus}}</td><td class="num">{{.Category}}</td><td>{{.Result}}</td><td class="num">{{.Cost}}</td><td>{{.Finding}}</td></tr>
      {{end}}
    </tbody>
  </table>
  {{end}}

  {{if eq .Category.TableType "graphql"}}
  <table>
    <thead><tr><th>Exp</th><th>Approach</th><th>Status</th><th>Key Question</th></tr></thead>
    <tbody>
      {{range .Category.Experiments}}
      <tr>
        <td>{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Num}}</a>{{else}}{{.Num}}{{end}}</td>
        <td>{{.Focus}}</td><td>{{.Result}}</td><td>{{.Finding}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>
  {{end}}

  {{if eq .Category.TableType "model-comparison"}}
  <table>
    <thead><tr><th>Model</th><th>Code Pass Rate</th><th class="num">SSIM</th><th class="num">Cost/Call</th><th>Visual Quality</th></tr></thead>
    <tbody>
      {{range .Category.Experiments}}
      <tr><td>{{.Focus}}</td><td>{{.Result}}</td><td class="num">{{.Category}}</td><td class="num">{{.Cost}}</td><td>{{.Finding}}</td></tr>
      {{end}}
    </tbody>
  </table>
  {{end}}

  {{if eq .Category.TableType "business-pipeline"}}
  <table>
    <thead><tr><th>Exp</th><th>Focus</th><th>Category</th><th>Key Output</th></tr></thead>
    <tbody>
      {{range .Category.Experiments}}
      <tr>
        <td>{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Num}}</a>{{else}}{{.Num}}{{end}}</td>
        <td>{{.Focus}}</td><td>{{.Category}}</td><td>{{.Finding}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>
  {{end}}
</div>
{{template "layout_end" .}}
{{end}}

{{define "experiment"}}
{{template "layout_start" .}}
{{if gt .TotalExps 0}}
<div class="progress-bar">
  <div class="progress-label">Experiment {{.ExpIndex}} of {{.TotalExps}}</div>
  <div class="progress-track"><div class="progress-fill" style="width:{{progressPct .ExpIndex .TotalExps}}%"></div></div>
</div>
{{end}}
<div class="section">
  <div class="exp-title-area">
    <h2>{{.Experiment.Focus}}</h2>
  </div>

  <div class="exp-result-card">
    <div>
      <div class="result-label">{{if .Experiment.IsFailure}}What Went Wrong{{else}}Key Finding{{end}}</div>
      <div class="result-value">{{.Experiment.Result}}</div>
      {{if .Experiment.Finding}}<div class="result-finding"><strong>Key finding:</strong> {{.Experiment.Finding}}</div>{{end}}
    </div>
    <div class="result-meta">
      {{if .Experiment.Cost}}<div class="result-cost">{{.Experiment.Cost}}</div>{{end}}
      {{if .Experiment.Category}}<div class="result-time">{{.Experiment.Category}}</div>{{end}}
    </div>
  </div>

  {{if .Experiment.CloneShot}}
  <div class="exp-screenshots">
    <div class="exp-screenshot-pair">
      {{if .Experiment.RefShot}}
      <div class="exp-screenshot">
        <div class="exp-screenshot-label">Original</div>
        <a href="{{.Experiment.RefShot}}" target="_blank"><img src="{{.Experiment.RefShot}}" alt="Original website" loading="lazy"></a>
      </div>
      {{end}}
      <div class="exp-screenshot">
        <div class="exp-screenshot-label">AI Clone</div>
        <a href="{{.Experiment.CloneShot}}" target="_blank"><img src="{{.Experiment.CloneShot}}" alt="AI clone" loading="lazy"></a>
      </div>
    </div>
  </div>
  {{end}}

  <div class="exp-detail">
    <div class="meta">
      {{if .Experiment.Result}}<div class="meta-item"><div class="meta-label">Result</div><div class="meta-value">{{.Experiment.Result}}</div></div>{{end}}
      {{if .Experiment.Cost}}<div class="meta-item"><div class="meta-label">Cost</div><div class="meta-value">{{.Experiment.Cost}}</div></div>{{end}}
      {{if .Experiment.Category}}<div class="meta-item"><div class="meta-label">Category</div><div class="meta-value">{{.Experiment.Category}}</div></div>{{end}}
      {{if gt .Experiment.ReadingTime 0}}<div class="meta-item"><div class="meta-label">Reading Time</div><div class="meta-value">{{.Experiment.ReadingTime}} min</div></div>{{end}}
    </div>
  </div>

  {{if .Experiment.Why}}
  <div class="exp-section section-why">
    <h3>Why</h3>
    <p>{{.Experiment.Why}}</p>
  </div>
  {{end}}

  {{if .Experiment.What}}
  <div class="exp-section section-what">
    <h3>What</h3>
    <p>{{.Experiment.What}}</p>
  </div>
  {{end}}

  {{if .Experiment.How}}
  <div class="exp-section section-how">
    <h3>How</h3>
    <p>{{.Experiment.How}}</p>
  </div>
  {{end}}

  {{if .Experiment.Result}}
  {{if .Experiment.IsFailure}}
  <div class="exp-section" style="border-left-color:#dc3545; background:#fdf2f3;">
    <h3 style="color:#dc3545;">Result</h3>
    <p style="font-size:18px; font-weight:600; color:#dc3545;">{{.Experiment.Result}}</p>
  </div>
  {{else}}
  <div class="exp-section" style="border-left-color:#28a745; background:#f0fdf4;">
    <h3 style="color:#28a745;">Result</h3>
    <p style="font-size:18px; font-weight:600;">{{.Experiment.Result}}</p>
  </div>
  {{end}}
  {{end}}

  {{if .Experiment.Impact}}
  <div class="exp-section section-impact">
    <h3>Impact on Pipeline</h3>
    <p>{{.Experiment.Impact}}</p>
  </div>
  {{end}}

  {{if .Experiment.Related}}
  <div class="exp-section">
    <h3>Related Experiments</h3>
    <div class="related-cards">
      {{range .Experiment.Related}}
      {{$rel := getExp .}}
      {{if $rel}}
      <a class="related-card" href="/exp/{{.}}">
        <div class="related-card-num">{{if $rel.Icon}}{{$rel.Icon}} {{end}}Exp {{$rel.Num}}</div>
        <div class="related-card-title">{{$rel.Focus}}</div>
        <div class="related-card-finding">{{$rel.Finding}}</div>
      </a>
      {{else}}
      <a href="/exp/{{.}}" style="margin-right:12px;">Experiment {{.}}</a>
      {{end}}
      {{end}}
    </div>
  </div>
  {{end}}

  {{range .Experiment.Description}}
  <p style="margin-top:16px; line-height:1.7;">{{.}}</p>
  {{end}}

  {{if .SourceCode}}
  <div class="source-section">
    <h3>Source Code</h3>
    <div class="source-file">{{.SourceName}}</div>
    <button class="code-toggle">Show full source</button>
    <pre class="source-pre line-numbers"><code class="language-{{.SourceLang}}">{{.SourceCode}}</code></pre>
  </div>
  {{end}}
</div>
<div class="exp-nav">
  {{if prevExp .Experiment.NumID}}<a href="/exp/{{prevExp .Experiment.NumID}}">&larr; Experiment {{prevExp .Experiment.NumID}}</a>{{else}}<span></span>{{end}}
  <a href="/category/{{.Category.Slug}}">&uarr; {{.Category.Name}}</a>
  {{if nextExp .Experiment.NumID}}<a href="/exp/{{nextExp .Experiment.NumID}}">Experiment {{nextExp .Experiment.NumID}} &rarr;</a>{{else}}<span></span>{{end}}
</div>
{{if .JourneyMode}}
<div class="journey-nav">
  {{if gt .ExpIndex 1}}<a class="journey-btn" href="/journey/{{sub .ExpIndex 1}}">&larr; Chapter {{sub .ExpIndex 1}}</a>{{else}}<span></span>{{end}}
  <span class="journey-chapter">Chapter {{.ExpIndex}} of {{.TotalExps}}</span>
  {{if lt .ExpIndex .TotalExps}}<a class="journey-btn" href="/journey/{{add .ExpIndex 1}}">Chapter {{add .ExpIndex 1}} &rarr;</a>{{else}}<a class="journey-btn" href="/">Finish</a>{{end}}
</div>
{{end}}
{{template "layout_end" .}}
{{end}}

{{define "discovery"}}
{{template "layout_start" .}}
<div class="section">
  <h2>Discovery Graph</h2>
  <p style="margin-bottom:24px; color:var(--text-muted);">How experiments connect to each other. Drag nodes to rearrange, scroll to zoom, click any node to read the experiment.</p>

  <div id="d3graph" style="width:100%; height:800px; border:1px solid var(--border); border-radius:8px; overflow:hidden; background:#fafafa;"></div>

  <details style="margin-top:24px;">
    <summary style="cursor:pointer; color:var(--accent); font-weight:600;">Static Mermaid View (for export)</summary>
    <div style="margin-top:12px; overflow-x:auto;">
      <pre class="mermaid" style="font-size:10px;">
{{.DiscoveryGraph}}
      </pre>
    </div>
  </details>
</div>
<script src="https://d3js.org/d3.v7.min.js"></script>
<script>
(function() {
  // Parse the mermaid source for nodes and links
  var mermaidEl = document.querySelector("details .mermaid");
  if (!mermaidEl) return;
  var src = mermaidEl.textContent || "";
  var nodes = [], nodeMap = {}, links = [];

  // Category colors
  var catColors = {
    "code": "#4A90D9", "design": "#7B68EE", "review": "#E8878E",
    "test": "#F5A623", "clone": "#50C878", "business": "#999",
    "graphql": "#E535AB", "pipeline": "#D4727A"
  };
  function getColor(id) {
    if (id >= 901) return catColors.pipeline;
    if (id >= 91) return catColors.graphql;
    if (id >= 55) return catColors.business;
    if (id >= 40) return catColors.test;
    if (id >= 27) return catColors.review;
    if (id >= 22) return catColors.design;
    return catColors.code;
  }

  var nodeRe = /(\d+)\["(\d+): ([^"]+)"\]/g;
  var m;
  while ((m = nodeRe.exec(src)) !== null) {
    var nid = parseInt(m[1]);
    if (!nodeMap[nid]) {
      var n = {id:nid, label:m[3], num:m[2], color:getColor(nid)};
      nodes.push(n);
      nodeMap[nid] = n;
    }
  }
  var linkRe = /\b(\d+)\["[^"]*"\]\s*-->\s*(\d+)\[/g;
  while ((m = linkRe.exec(src)) !== null) {
    var s = parseInt(m[1]), t = parseInt(m[2]);
    if (nodeMap[s] && nodeMap[t]) links.push({source:s, target:t});
  }
  if (nodes.length === 0) return;

  // Only keep nodes that have at least one edge
  var connectedIds = new Set();
  links.forEach(function(l) { connectedIds.add(l.source); connectedIds.add(l.target); });
  nodes = nodes.filter(function(n) { return connectedIds.has(n.id); });

  var container = document.getElementById("d3graph");
  var width = container.clientWidth, height = 800;
  var svg = d3.select("#d3graph").append("svg")
    .attr("width", width).attr("height", height)
    .style("font-family", "-apple-system, sans-serif");
  var g = svg.append("g");

  // Zoom
  svg.call(d3.zoom().scaleExtent([0.15, 4]).on("zoom", function(e) {
    g.attr("transform", e.transform);
  }));

  // Force simulation with more spacing
  var sim = d3.forceSimulation(nodes)
    .force("link", d3.forceLink(links).id(function(d){return d.id;}).distance(220))
    .force("charge", d3.forceManyBody().strength(-800))
    .force("center", d3.forceCenter(width/2, height/2))
    .force("collision", d3.forceCollide().radius(80));

  // Links
  var link = g.append("g").selectAll("line").data(links).enter().append("line")
    .attr("stroke", "#ccc").attr("stroke-opacity", 0.5).attr("stroke-width", 1.5)
    .attr("marker-end", "url(#arrow)");

  // Arrow marker
  svg.append("defs").append("marker")
    .attr("id", "arrow").attr("viewBox", "0 0 10 10")
    .attr("refX", 20).attr("refY", 5)
    .attr("markerWidth", 6).attr("markerHeight", 6)
    .attr("orient", "auto")
    .append("path").attr("d", "M0,0 L10,5 L0,10 Z").attr("fill", "#ccc");

  // Nodes
  var node = g.append("g").selectAll("g").data(nodes).enter().append("g")
    .attr("cursor", "pointer")
    .call(d3.drag()
      .on("start", function(e,d) { if(!e.active) sim.alphaTarget(0.3).restart(); d.fx=d.x; d.fy=d.y; })
      .on("drag", function(e,d) { d.fx=e.x; d.fy=e.y; })
      .on("end", function(e,d) { if(!e.active) sim.alphaTarget(0); d.fx=null; d.fy=null; }));

  // Node circles
  node.append("circle")
    .attr("r", 10)
    .attr("fill", function(d) { return d.color; })
    .attr("stroke", "#fff").attr("stroke-width", 2);

  // Node number badge
  node.append("text")
    .text(function(d) { return d.num; })
    .attr("text-anchor", "middle").attr("dy", "0.35em")
    .attr("font-size", "8px").attr("font-weight", "700").attr("fill", "#fff");

  // Node label
  node.append("text")
    .text(function(d) { return d.label.length > 25 ? d.label.substring(0,22) + "..." : d.label; })
    .attr("x", 16).attr("y", 4)
    .attr("font-size", "11px").attr("fill", "#333");

  // Hover: show full label
  node.append("title").text(function(d) { return "Exp " + d.num + ": " + d.label; });

  // Click: navigate
  node.on("click", function(e, d) { window.location.href = "/exp/" + d.id; });

  // Highlight on hover
  node.on("mouseover", function() {
    d3.select(this).select("circle").transition().duration(200).attr("r", 14);
  }).on("mouseout", function() {
    d3.select(this).select("circle").transition().duration(200).attr("r", 10);
  });

  sim.on("tick", function() {
    link.attr("x1", function(d){return d.source.x;}).attr("y1", function(d){return d.source.y;})
      .attr("x2", function(d){return d.target.x;}).attr("y2", function(d){return d.target.y;});
    node.attr("transform", function(d){return "translate("+d.x+","+d.y+")";});
  });

  // Legend
  var legend = svg.append("g").attr("transform", "translate(20, 20)");
  var cats = [
    {label:"Code Gen (1-21)", color:catColors.code},
    {label:"Design (22-26)", color:catColors.design},
    {label:"Review (27-39)", color:catColors.review},
    {label:"Testing (40-52)", color:catColors.test},
    {label:"Business (55-84)", color:catColors.business},
    {label:"Pipeline Steps", color:catColors.pipeline},
  ];
  cats.forEach(function(c, i) {
    legend.append("circle").attr("cx", 8).attr("cy", i*20+8).attr("r", 6).attr("fill", c.color);
    legend.append("text").attr("x", 20).attr("y", i*20+12).text(c.label)
      .attr("font-size", "11px").attr("fill", "#666");
  });
})();
</script>
{{template "layout_end" .}}
{{end}}

{{define "timeline"}}
{{template "layout_start" .}}
<div class="section">
  <h2>Research Timeline</h2>
  <p style="color:var(--text-muted); margin-bottom:24px;">The chronological journey from "can AI write code?" to a complete autonomous pipeline. Click any experiment to read the full story.</p>

  <style>
    .timeline { position:relative; padding-left:40px; }
    .timeline::before { content:''; position:absolute; left:15px; top:0; bottom:0; width:2px; background:var(--border); }
    .tl-phase { margin-bottom:40px; }
    .tl-phase-title { font-size:18px; font-weight:700; margin-bottom:16px; position:relative; }
    .tl-phase-title::before { content:''; position:absolute; left:-33px; top:6px; width:14px; height:14px; border-radius:50%; border:3px solid var(--accent); background:var(--bg); z-index:1; }
    .tl-item { display:flex; gap:16px; align-items:flex-start; margin-bottom:12px; padding:12px 16px; background:var(--table-stripe); border-radius:8px; transition:all 0.2s; position:relative; }
    .tl-item:hover { background:var(--accent-light); }
    .tl-item::before { content:''; position:absolute; left:-29px; top:18px; width:8px; height:8px; border-radius:50%; background:var(--accent); }
    .tl-num { font-size:13px; font-weight:700; color:var(--accent); min-width:36px; }
    .tl-content { flex:1; }
    .tl-title { font-weight:600; font-size:14px; }
    .tl-title a { color:var(--text); }
    .tl-title a:hover { color:var(--accent); }
    .tl-finding { font-size:13px; color:var(--text-muted); margin-top:2px; }
    .tl-meta { display:flex; gap:12px; font-size:12px; color:var(--text-light); margin-top:4px; }
    .tl-cost { font-weight:600; }
    .tl-result { padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600; }
    .tl-pass { background:#d4edda; color:#155724; }
    .tl-fail { background:#fce4ec; color:#b71c1c; }
    .tl-milestone { margin:24px 0; padding:16px 20px; background:var(--accent-light); border-left:3px solid var(--accent); border-radius:0 8px 8px 0; position:relative; }
    .tl-milestone::before { content:'⭐'; position:absolute; left:-37px; top:14px; font-size:16px; }
    .tl-milestone strong { color:var(--accent); }
  </style>

  <div class="timeline">
    <!-- Phase 1: Can AI Write Code? -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 1: Can AI Write Code?</div>
      {{range .Categories}}{{if eq .Slug "code-generation"}}
      {{range .Experiments}}
      <div class="tl-item">
        <div class="tl-num">{{.Num}}</div>
        <div class="tl-content">
          <div class="tl-title">{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Focus}}</a>{{else}}{{.Focus}}{{end}}</div>
          <div class="tl-finding">{{.Finding}}</div>
          <div class="tl-meta">
            {{if .Cost}}<span class="tl-cost">{{.Cost}}</span>{{end}}
            {{if .Result}}<span class="tl-result {{if .IsFailure}}tl-fail{{else}}tl-pass{{end}}">{{.Result}}</span>{{end}}
          </div>
        </div>
      </div>
      {{end}}
      {{end}}{{end}}

      <div class="tl-milestone"><strong>Milestone:</strong> Prompt wording > model choice. Cheapest model hits 100% with V4 prompts. Auto-fix recovers 40-60% of failures for free.</div>
    </div>

    <!-- Phase 2: Can AI Design Products? -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 2: Can AI Design Products?</div>
      {{range .Categories}}{{if eq .Slug "product-design"}}
      {{range .Experiments}}
      <div class="tl-item">
        <div class="tl-num">{{.Num}}</div>
        <div class="tl-content">
          <div class="tl-title">{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Focus}}</a>{{else}}{{.Focus}}{{end}}</div>
          <div class="tl-finding">{{.Finding}}</div>
          <div class="tl-meta">
            {{if .Cost}}<span class="tl-cost">{{.Cost}}</span>{{end}}
            {{if .Result}}<span class="tl-result {{if .IsFailure}}tl-fail{{else}}tl-pass{{end}}">{{.Result}}</span>{{end}}
          </div>
        </div>
      </div>
      {{end}}
      {{end}}{{end}}

      <div class="tl-milestone"><strong>Milestone:</strong> $0.04 of design turns a code generator into a product builder. Persona interviews find features the brief missed.</div>
    </div>

    <!-- Phase 3: Can AI Review Its Own Work? -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 3: Can AI Review Its Own Work?</div>
      {{range .Categories}}{{if eq .Slug "review-pipeline"}}
      {{range .Experiments}}
      <div class="tl-item">
        <div class="tl-num">{{.Num}}</div>
        <div class="tl-content">
          <div class="tl-title">{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Focus}}</a>{{else}}{{.Focus}}{{end}}</div>
          <div class="tl-finding">{{.Finding}}</div>
          <div class="tl-meta">
            {{if .Cost}}<span class="tl-cost">{{.Cost}}</span>{{end}}
            {{if .Result}}<span class="tl-result {{if .IsFailure}}tl-fail{{else}}tl-pass{{end}}">{{.Result}}</span>{{end}}
          </div>
        </div>
      </div>
      {{end}}
      {{end}}{{end}}

      <div class="tl-milestone"><strong>Milestone:</strong> 20 reviewers, no diminishing returns. Domain expert found ALL 10 missing features. Progressive enhancement: ZERO regressions.</div>
    </div>

    <!-- Phase 4: Can AI Test and Secure? -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 4: Can AI Test and Secure?</div>
      {{range .Categories}}{{if eq .Slug "testing-security"}}
      {{range .Experiments}}
      <div class="tl-item">
        <div class="tl-num">{{.Num}}</div>
        <div class="tl-content">
          <div class="tl-title">{{if and .HasDetail (gt .NumID 0)}}<a href="/exp/{{.Num}}">{{.Focus}}</a>{{else}}{{.Focus}}{{end}}</div>
          <div class="tl-finding">{{.Finding}}</div>
          <div class="tl-meta">
            {{if .Cost}}<span class="tl-cost">{{.Cost}}</span>{{end}}
            {{if .Result}}<span class="tl-result {{if .IsFailure}}tl-fail{{else}}tl-pass{{end}}">{{.Result}}</span>{{end}}
          </div>
        </div>
      </div>
      {{end}}
      {{end}}{{end}}

      <div class="tl-milestone"><strong>Milestone:</strong> Playwright catches 4 bugs that 26 tests + 10 reviewers missed. TDD: 90.3% coverage. Go server survives all chaos.</div>
    </div>

    <!-- Phase 5: The Full Business Pipeline -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 5: The Full Business Pipeline</div>
      <div class="tl-item">
        <div class="tl-num">55-84</div>
        <div class="tl-content">
          <div class="tl-title"><a href="/category/business-pipeline">30 experiments: Docker, TypeScript, PostgreSQL, Auth, Marketing, Legal, Support, Localisation</a></div>
          <div class="tl-finding">Complete business pipeline from idea discovery to legal documents, all under $1 total</div>
          <div class="tl-meta"><span class="tl-cost">~$1.00 total</span><span class="tl-result tl-pass">30/30</span></div>
        </div>
      </div>

      <div class="tl-milestone"><strong>Milestone:</strong> Every non-code artifact a startup needs — marketing, legal, support — can be AI-generated as a solid first draft.</div>
    </div>

    <!-- Phase 6: Visual Cloning -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 6: Can AI Clone Websites?</div>
      <div class="tl-item">
        <div class="tl-num">54</div>
        <div class="tl-content">
          <div class="tl-title"><a href="/category/website-cloning">9-Model Website Clone Comparison</a></div>
          <div class="tl-finding">Opus best visual ($2.75), Llama 4 Scout best value (~$0.02). SSIM pixel metrics misleading.</div>
          <div class="tl-meta"><span class="tl-cost">$0.02-$2.75</span><span class="tl-result tl-pass">9 models tested</span></div>
        </div>
      </div>
      <div class="tl-item">
        <div class="tl-num">95</div>
        <div class="tl-content">
          <div class="tl-title"><a href="/exp/95">AI Image Generation — Nano Banana 2 vs GPT-5</a></div>
          <div class="tl-finding">Nano Banana 2: 5x faster, 70x cheaper. Photorealistic London apartment photos.</div>
          <div class="tl-meta"><span class="tl-cost">$0.03 total</span><span class="tl-result tl-pass">9/9 images</span></div>
        </div>
      </div>

      <div class="tl-milestone"><strong>Milestone:</strong> AI generates layout + AI generates images = complete product page with zero real photography.</div>
    </div>

    <!-- Phase 7: What's Next -->
    <div class="tl-phase">
      <div class="tl-phase-title">Phase 7: The Composable Pipeline</div>
      <div class="tl-item">
        <div class="tl-num">91</div>
        <div class="tl-content">
          <div class="tl-title"><a href="/exp/91">GraphQL Without Codegen</a></div>
          <div class="tl-finding">Minimal hand-rolled (231 lines) beats gqlgen (7 files, 5x cost)</div>
          <div class="tl-meta"><span class="tl-cost">$0.045</span><span class="tl-result tl-pass">BUILD PASS</span></div>
        </div>
      </div>
      <div class="tl-item">
        <div class="tl-num">99-101</div>
        <div class="tl-content">
          <div class="tl-title"><a href="/category/composable-architecture">Composable Block System (planned)</a></div>
          <div class="tl-finding">YAML block definitions + Go orchestrator + visual editor. The system builds its own blocks.</div>
          <div class="tl-meta"><span class="tl-result tl-pass">Design complete</span></div>
        </div>
      </div>

      <div class="tl-milestone"><strong>Next:</strong> Every experiment is a tested block. The pipeline that builds blocks is itself a pipeline of blocks.</div>
    </div>
  </div>
</div>
{{template "layout_end" .}}
{{end}}

{{define "compare"}}
{{template "layout_start" .}}
<div class="section">
  <h2>Visual Comparison</h2>
  <p style="color:var(--text-muted); margin-bottom:8px;">20 websites cloned with Opus across 20 industries. Click any card to see the full comparison.</p>
  <p style="color:var(--text-muted); margin-bottom:32px; font-size:14px;">Universal spacing (4-40px) and component patterns validated across all clones.</p>

  <div style="display:grid; grid-template-columns:repeat(auto-fill, minmax(260px, 1fr)); gap:20px; margin-bottom:32px;">
    {{range .CloneSites}}
    <a href="/compare/{{.Slug}}" style="text-decoration:none; color:inherit; border:1px solid var(--border); border-radius:12px; overflow:hidden; transition:box-shadow 0.2s ease, transform 0.2s ease; display:block;">
      <div style="height:180px; overflow:hidden; background:var(--table-stripe);">
        <img src="{{.ClonePrefix}}/desktop.png" style="width:100%; height:180px; object-fit:cover; object-position:top;" alt="{{.Name}} clone">
      </div>
      <div style="padding:12px 14px;">
        <strong style="font-size:15px;">{{.Name}}</strong>
        <span style="color:var(--text-muted); font-size:13px; margin-left:6px;">{{.Category}}</span>
        <div style="margin-top:6px; display:flex; align-items:center; gap:8px; font-size:13px; color:var(--text-light);">
          <span style="display:inline-block; width:12px; height:12px; border-radius:50%; background:{{.PrimaryColor}};"></span>
          <span>{{.Lines}} lines · {{.Iterations}} iter</span>
          {{if .Method}}<span style="padding:1px 8px; border-radius:8px; font-size:11px; font-weight:600; {{if eq .Method "Hybrid"}}background:#d4edda; color:#155724;{{else if eq .Method "Screenshot-Guided"}}background:#cce5ff; color:#004085;{{else}}background:#f8f9fa; color:#6c757d;{{end}}">{{.Method}}</span>{{end}}
        </div>
      </div>
    </a>
    {{end}}
  </div>
</div>
<style>
  .section a[href^="/compare/"]:hover {
    box-shadow: 0 4px 16px rgba(0,0,0,0.10);
    transform: translateY(-2px);
  }
</style>
{{template "layout_end" .}}
{{end}}

{{define "compare-detail"}}
{{template "layout_start" .}}
<div class="section">
  <h2 style="margin-bottom:8px;">{{.CloneSite.Name}}</h2>
  <div style="display:flex; flex-wrap:wrap; align-items:center; gap:12px; margin-bottom:28px; font-size:14px; color:var(--text-muted);">
    <span style="display:inline-flex; align-items:center; gap:6px;">
      <span style="display:inline-block; width:14px; height:14px; border-radius:50%; background:{{.CloneSite.PrimaryColor}};"></span>
      {{.CloneSite.Category}}
    </span>
    <span style="color:var(--border);">|</span>
    <span>{{.CloneSite.Lines}} lines</span>
    <span style="color:var(--border);">|</span>
    <span>{{.CloneSite.Iterations}} iterations</span>
    <span style="color:var(--border);">|</span>
    <span>Port {{.CloneSite.Port}}</span>
    {{if .CloneSite.Method}}
    <span style="color:var(--border);">|</span>
    <span style="padding:2px 10px; border-radius:10px; font-size:12px; font-weight:600; {{if eq .CloneSite.Method "Hybrid"}}background:#d4edda; color:#155724;{{else if eq .CloneSite.Method "Screenshot-Guided"}}background:#cce5ff; color:#004085;{{else}}background:#f8f9fa; color:#6c757d;{{end}}">{{.CloneSite.Method}}</span>
    {{end}}
  </div>

  <div style="display:flex; gap:8px; margin-bottom:16px; align-items:center;">
    {{if .CloneSite.RefPrefix}}
    <button class="compare-tab active" onclick="showTab(this, 'tab-original')" style="padding:8px 20px; border:1px solid var(--border); border-radius:6px; background:var(--accent); color:#fff; cursor:pointer; font-size:14px; font-weight:600;">Original</button>
    <button class="compare-tab" onclick="showTab(this, 'tab-clone')" style="padding:8px 20px; border:1px solid var(--border); border-radius:6px; background:var(--bg); color:var(--text); cursor:pointer; font-size:14px; font-weight:600;">AI Clone</button>
    <button class="compare-tab" onclick="showTab(this, 'tab-sidebyside')" style="padding:8px 20px; border:1px solid var(--border); border-radius:6px; background:var(--bg); color:var(--text); cursor:pointer; font-size:14px; font-weight:600;">Side by Side</button>
    {{else}}
    <button class="compare-tab active" onclick="showTab(this, 'tab-clone')" style="padding:8px 20px; border:1px solid var(--border); border-radius:6px; background:var(--accent); color:#fff; cursor:pointer; font-size:14px; font-weight:600;">AI Clone</button>
    {{end}}
    <span style="flex:1;"></span>
    {{if gt .CloneSite.Port 0}}
    <button class="compare-tab" onclick="showTab(this, 'tab-live')" style="padding:8px 20px; border:1px solid var(--border); border-radius:6px; background:var(--bg); color:var(--text); cursor:pointer; font-size:14px; font-weight:600;">Live Preview</button>
    <span style="flex:1;"></span>
    <a href="http://localhost:{{.CloneSite.Port}}" target="_blank" style="padding:8px 20px; border:1px solid var(--accent); border-radius:6px; color:var(--accent); font-size:14px; font-weight:600; text-decoration:none;">Open in New Tab &rarr;</a>
    {{end}}
  </div>

  {{if .CloneSite.RefPrefix}}
  <div id="tab-original" class="compare-panel" style="margin-bottom:36px;">
    <a href="{{.CloneSite.RefPrefix}}/desktop.png" target="_blank">
      <img src="{{.CloneSite.RefPrefix}}/desktop.png" style="width:100%; max-width:900px; border:1px solid var(--border); border-radius:10px;" alt="Original {{.CloneSite.Name}}">
    </a>
  </div>
  <div id="tab-clone" class="compare-panel" style="margin-bottom:36px; display:none;">
    <a href="{{.CloneSite.ClonePrefix}}/desktop.png" target="_blank">
      <img src="{{.CloneSite.ClonePrefix}}/desktop.png" style="width:100%; max-width:900px; border:1px solid var(--border); border-radius:10px;" alt="{{.CloneSite.Name}} clone">
    </a>
  </div>
  <div id="tab-sidebyside" class="compare-panel" style="margin-bottom:36px; display:none;">
    <div style="display:grid; grid-template-columns:1fr 1fr; gap:16px;">
      <div>
        <div style="font-size:12px; font-weight:600; color:var(--text-muted); margin-bottom:8px; text-align:center;">Original</div>
        <a href="{{.CloneSite.RefPrefix}}/desktop.png" target="_blank">
          <img src="{{.CloneSite.RefPrefix}}/desktop.png" style="width:100%; border:1px solid var(--border); border-radius:10px;" alt="Original">
        </a>
      </div>
      <div>
        <div style="font-size:12px; font-weight:600; color:var(--text-muted); margin-bottom:8px; text-align:center;">AI Clone</div>
        <a href="{{.CloneSite.ClonePrefix}}/desktop.png" target="_blank">
          <img src="{{.CloneSite.ClonePrefix}}/desktop.png" style="width:100%; border:1px solid var(--border); border-radius:10px;" alt="Clone">
        </a>
      </div>
    </div>
  </div>
  {{else}}
  <div id="tab-clone" class="compare-panel" style="margin-bottom:36px;">
    <a href="{{.CloneSite.ClonePrefix}}/desktop.png" target="_blank">
      <img src="{{.CloneSite.ClonePrefix}}/desktop.png" style="width:100%; max-width:900px; border:1px solid var(--border); border-radius:10px;" alt="{{.CloneSite.Name}} clone">
    </a>
  </div>
  {{end}}

  {{if gt .CloneSite.Port 0}}
  <div id="tab-live" class="compare-panel" style="margin-bottom:36px; display:none;">
    <div style="background:var(--table-stripe); border:1px solid var(--border); border-radius:10px; overflow:hidden;">
      <div style="padding:8px 16px; background:var(--bg); border-bottom:1px solid var(--border); display:flex; align-items:center; gap:8px; font-size:13px; color:var(--text-muted);">
        <span style="display:flex; gap:4px;"><span style="width:10px; height:10px; border-radius:50%; background:#ff5f57;"></span><span style="width:10px; height:10px; border-radius:50%; background:#febc2e;"></span><span style="width:10px; height:10px; border-radius:50%; background:#28c840;"></span></span>
        <span style="flex:1; text-align:center;">localhost:{{.CloneSite.Port}}</span>
      </div>
      <iframe src="http://localhost:{{.CloneSite.Port}}" style="width:100%; height:700px; border:none;" loading="lazy"></iframe>
    </div>
    <div style="font-size:12px; color:var(--text-light); margin-top:8px; text-align:center;">Live clone running on port {{.CloneSite.Port}}. If blank, the server may not be running.</div>
  </div>
  {{end}}

  <script>
  function showTab(btn, panelId) {
    document.querySelectorAll('.compare-tab').forEach(function(t) {
      t.style.background = 'var(--bg)'; t.style.color = 'var(--text)';
    });
    btn.style.background = 'var(--accent)'; btn.style.color = '#fff';
    document.querySelectorAll('.compare-panel').forEach(function(p) { p.style.display = 'none'; });
    var panel = document.getElementById(panelId);
    if (panel) panel.style.display = 'block';
  }
  </script>

  {{if .CloneSite.AIImages}}
  <h3 style="margin:36px 0 16px; padding-bottom:8px; border-bottom:1px solid var(--border);">AI-Generated Images</h3>
  <div style="display:grid; grid-template-columns:repeat(auto-fill, minmax(200px, 1fr)); gap:14px; margin-bottom:36px;">
    {{range .CloneSite.AIImages}}
    <a href="{{$.CloneSite.AIPrefix}}/{{.}}" target="_blank" style="display:block;">
      <img src="{{$.CloneSite.AIPrefix}}/{{.}}" style="width:100%; height:160px; object-fit:cover; border-radius:10px; border:1px solid var(--border);" alt="{{.}}">
    </a>
    {{end}}
  </div>
  {{end}}

  <div style="display:flex; justify-content:space-between; align-items:center; padding:20px 0; border-top:1px solid var(--border); margin-top:24px;">
    <div>
      {{if prevClone .CloneSite.Slug}}
      <a href="/compare/{{prevClone .CloneSite.Slug}}" style="display:inline-flex; align-items:center; gap:6px; font-size:14px;">
        <span>&larr;</span> {{prevCloneName .CloneSite.Slug}}
      </a>
      {{end}}
    </div>
    <a href="/compare/" style="font-size:14px; color:var(--text-muted);">All Clones</a>
    <div>
      {{if nextClone .CloneSite.Slug}}
      <a href="/compare/{{nextClone .CloneSite.Slug}}" style="display:inline-flex; align-items:center; gap:6px; font-size:14px;">
        {{nextCloneName .CloneSite.Slug}} <span>&rarr;</span>
      </a>
      {{end}}
    </div>
  </div>
</div>
{{template "layout_end" .}}
{{end}}

{{define "score-chart"}}
{{template "layout_start" .}}
<div class="section">
  <h2>SSIM Score Comparison</h2>
  <p style="margin-bottom:24px; color:var(--text-muted);">All website cloning experiments ranked by overall SSIM score.</p>
  <div style="overflow-x:auto;">
  {{$scores := sortedScores}}
  {{range $i, $e := $scores}}
  <div style="display:flex; align-items:center; gap:12px; margin-bottom:8px;">
    <a href="/exp/{{$e.Num}}" style="width:200px; font-size:13px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; color:var(--text);">{{$e.Name}}</a>
    <div style="flex:1; height:24px; background:var(--border); border-radius:4px; overflow:hidden; max-width:400px;">
      <div style="height:100%; width:{{printf "%.1f" (mul $e.Score 100.0)}}%; background:{{$e.Color}}; border-radius:4px; transition:width 0.3s;"></div>
    </div>
    <span style="font-size:13px; font-weight:600; min-width:50px; color:{{$e.Color}};">{{printf "%.3f" $e.Score}}</span>
    <span style="font-size:12px; color:var(--text-light);">{{$e.Cost}}</span>
  </div>
  {{end}}
  </div>
</div>
{{template "layout_end" .}}
{{end}}

{{define "roadmap"}}
{{template "layout_start" .}}
{{template "breadcrumbs" .}}
<div class="hero" style="border-bottom:none; padding-bottom:16px;">
  <h1>Experiment Roadmap</h1>
  <p class="subtitle">Priority experiments for the website cloning pipeline</p>
</div>
<div class="section">
  <table>
    <thead>
      <tr>
        <th style="width:40px">#</th>
        <th>Experiment</th>
        <th>Description</th>
        <th style="width:80px">Difficulty</th>
        <th style="width:70px">Impact</th>
        <th style="width:90px">Status</th>
        <th>Related</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="num">1</td>
        <td style="font-weight:600;">Multi-page Clone</td>
        <td>Clone an entire site (3-5 pages) with shared nav, consistent design tokens, and inter-page navigation. Currently we clone single pages only.</td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/165">165</a>, <a href="/exp/168">168</a></td>
      </tr>
      <tr>
        <td class="num">2</td>
        <td style="font-weight:600;">Responsive Clone (Mobile + Tablet)</td>
        <td>Extract and reproduce media queries. Clone should look correct at 375px, 768px, and 1440px viewports.</td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/133">133</a>, <a href="/exp/152">152</a></td>
      </tr>
      <tr>
        <td class="num">3</td>
        <td style="font-weight:600;">Interactive Clone (JS)</td>
        <td>Reproduce JavaScript behaviors: dropdowns, accordions, scroll animations, modals. Move beyond static HTML/CSS.</td>
        <td><span style="background:#ef4444; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Hard</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/166">166</a>, <a href="/exp/167">167</a></td>
      </tr>
      <tr>
        <td class="num">4</td>
        <td style="font-weight:600;">Hover State Matching</td>
        <td>Use captured hover data (Exp 166-167) to reproduce exact transition properties, durations, and color values on hover/focus/active.</td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:var(--accent); color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">In Progress</span></td>
        <td><a href="/exp/166">166</a>, <a href="/exp/167">167</a></td>
      </tr>
      <tr>
        <td class="num">5</td>
        <td style="font-weight:600;">Carousel Reproduction</td>
        <td>Clone carousel/slider components with correct slide content, navigation dots/arrows, auto-play timing, and transition effects.</td>
        <td><span style="background:#ef4444; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Hard</span></td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/159">159</a>, <a href="/exp/160">160</a></td>
      </tr>
      <tr>
        <td class="num">6</td>
        <td style="font-weight:600;">Full-Stack Clone (API + DB + Auth)</td>
        <td>Clone not just the frontend but also the backend: API endpoints, database schema, and authentication flow.</td>
        <td><span style="background:#ef4444; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Hard</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/152">152</a>, <a href="/exp/165">165</a></td>
      </tr>
      <tr>
        <td class="num">7</td>
        <td style="font-weight:600;">Performance Scoring (Lighthouse)</td>
        <td>Run Lighthouse on clones and originals. Measure performance, accessibility, SEO, and best practices scores side-by-side.</td>
        <td><span style="background:#22c55e; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Easy</span></td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/162">162</a></td>
      </tr>
      <tr>
        <td class="num">8</td>
        <td style="font-weight:600;">Component Extraction (Web Components)</td>
        <td>Extract reusable Web Components from clones. Each component self-contained with Shadow DOM, scoped CSS, and slot-based content.</td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/119">119</a>, <a href="/exp/130">130</a></td>
      </tr>
      <tr>
        <td class="num">9</td>
        <td style="font-weight:600;">Design System Tokens (Interaction Tokens)</td>
        <td>Extend token extraction (Exp 118) to include interaction tokens: hover durations, easing curves, animation keyframes, scroll triggers.</td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:#fbbf24; color:#000; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Medium</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/118">118</a>, <a href="/exp/166">166</a></td>
      </tr>
      <tr>
        <td class="num">10</td>
        <td style="font-weight:600;">White-Label Pipeline</td>
        <td>One-command pipeline: input a URL, output a white-label clone with swapped branding, colors, content, and AI images. Design transfer at scale.</td>
        <td><span style="background:#ef4444; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Hard</span></td>
        <td><span style="background:#16a34a; color:#fff; padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">High</span></td>
        <td><span style="background:var(--border); color:var(--text-muted); padding:2px 8px; border-radius:10px; font-size:11px; font-weight:600;">Planned</span></td>
        <td><a href="/exp/163">163</a>, <a href="/exp/165">165</a></td>
      </tr>
    </tbody>
  </table>
</div>
{{template "layout_end" .}}
{{end}}

{{define "pipeline-stats"}}
{{template "layout_start" .}}
{{template "breadcrumbs" .}}
<div class="hero" style="border-bottom:none; padding-bottom:16px;">
  <h1>Pipeline Stats</h1>
  <p class="subtitle">Validated pipeline performance across all cloned sites</p>
</div>
<div class="section">
  <h2>Clone Performance Table</h2>
  <div style="overflow-x:auto; margin-top:16px;">
    <table>
      <thead>
        <tr>
          <th>Site</th>
          <th>Exp</th>
          <th class="num">Blocks</th>
          <th class="num">Lines</th>
          <th class="num">SSIM</th>
          <th class="num">Cost</th>
          <th class="num">Images</th>
          <th class="num">Iterations</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><strong>Nando's</strong></td>
          <td><a href="/exp/152">152</a></td>
          <td class="num">196</td>
          <td class="num">1,032</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.575}}">0.575</span></td>
          <td class="num">~$0.60</td>
          <td class="num">0</td>
          <td class="num">1</td>
        </tr>
        <tr>
          <td><strong>Nando's</strong> (refined)</td>
          <td><a href="/exp/155">155</a></td>
          <td class="num">196</td>
          <td class="num">&mdash;</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.622}}">0.622</span></td>
          <td class="num">~$1.00</td>
          <td class="num">&mdash;</td>
          <td class="num">3</td>
        </tr>
        <tr>
          <td><strong>Figma</strong></td>
          <td><a href="/exp/156">156</a></td>
          <td class="num">293</td>
          <td class="num">516</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.700}}">0.700</span></td>
          <td class="num">~$0.60</td>
          <td class="num">16</td>
          <td class="num">1</td>
        </tr>
        <tr>
          <td><strong>Stripe</strong></td>
          <td><a href="/exp/157">157</a></td>
          <td class="num">358</td>
          <td class="num">1,087</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.571}}">0.571</span></td>
          <td class="num">~$0.60</td>
          <td class="num">0</td>
          <td class="num">1</td>
        </tr>
        <tr>
          <td><strong>Linear</strong></td>
          <td><a href="/exp/161">161</a></td>
          <td class="num">221</td>
          <td class="num">&mdash;</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.328}}">0.328</span></td>
          <td class="num">~$0.20</td>
          <td class="num">4</td>
          <td class="num">1</td>
        </tr>
        <tr>
          <td><strong>SaaS Tools</strong></td>
          <td><a href="/exp/162">162</a></td>
          <td class="num">178</td>
          <td class="num">821</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.789}}">0.789</span></td>
          <td class="num">~$0.40</td>
          <td class="num">5</td>
          <td class="num">1</td>
        </tr>
        <tr>
          <td><strong>Vercel</strong></td>
          <td><a href="/exp/165">165</a></td>
          <td class="num">255</td>
          <td class="num">893</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.787}}">0.787</span></td>
          <td class="num">~$0.60</td>
          <td class="num">12</td>
          <td class="num">3</td>
        </tr>
        <tr>
          <td><strong>Notion</strong></td>
          <td><a href="/exp/168">168</a></td>
          <td class="num">275</td>
          <td class="num">1,277</td>
          <td class="num"><span class="score-badge" style="background:{{scoreColor 0.725}}">0.725</span></td>
          <td class="num">~$0.10</td>
          <td class="num">5</td>
          <td class="num">2</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>

<div class="section">
  <h2>SSIM Scores by Site</h2>
  <p style="color:var(--text-muted); margin-bottom:20px;">Sorted highest to lowest. Higher is better (1.0 = pixel-perfect match).</p>
  <svg viewBox="0 0 700 280" style="max-width:700px; width:100%; font-family:inherit;">
    <defs>
      <linearGradient id="bar-grad" x1="0%" y1="0%" x2="100%" y2="0%">
        <stop offset="0%" style="stop-color:var(--accent);stop-opacity:0.8"/>
        <stop offset="100%" style="stop-color:var(--accent);stop-opacity:1"/>
      </linearGradient>
    </defs>
    <!-- Grid lines -->
    <line x1="140" y1="10" x2="140" y2="270" stroke="#e0e0e0" stroke-width="1"/>
    <line x1="280" y1="10" x2="280" y2="270" stroke="#e0e0e0" stroke-width="0.5" stroke-dasharray="4"/>
    <line x1="420" y1="10" x2="420" y2="270" stroke="#e0e0e0" stroke-width="0.5" stroke-dasharray="4"/>
    <line x1="560" y1="10" x2="560" y2="270" stroke="#e0e0e0" stroke-width="0.5" stroke-dasharray="4"/>
    <line x1="690" y1="10" x2="690" y2="270" stroke="#e0e0e0" stroke-width="1"/>
    <!-- Axis labels -->
    <text x="140" y="278" text-anchor="middle" fill="#888" font-size="10">0</text>
    <text x="280" y="278" text-anchor="middle" fill="#888" font-size="10">0.25</text>
    <text x="420" y="278" text-anchor="middle" fill="#888" font-size="10">0.50</text>
    <text x="560" y="278" text-anchor="middle" fill="#888" font-size="10">0.75</text>
    <text x="690" y="278" text-anchor="middle" fill="#888" font-size="10">1.0</text>
    <!-- 1: SaaS Tools 0.789 -->
    <text x="135" y="34" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">SaaS Tools</text>
    <rect x="140" y="20" width="434" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="579" y="35" fill="#fff" font-size="11" font-weight="700">0.789</text>
    <!-- 2: Vercel 0.787 -->
    <text x="135" y="64" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Vercel</text>
    <rect x="140" y="50" width="433" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="578" y="65" fill="#fff" font-size="11" font-weight="700">0.787</text>
    <!-- 3: Notion 0.725 -->
    <text x="135" y="94" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Notion</text>
    <rect x="140" y="80" width="399" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="544" y="95" fill="#fff" font-size="11" font-weight="700">0.725</text>
    <!-- 4: Figma 0.700 -->
    <text x="135" y="124" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Figma</text>
    <rect x="140" y="110" width="385" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="530" y="125" fill="#fff" font-size="11" font-weight="700">0.700</text>
    <!-- 5: Nando's (refined) 0.622 -->
    <text x="135" y="154" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Nando's (v2)</text>
    <rect x="140" y="140" width="342" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="487" y="155" fill="#fff" font-size="11" font-weight="700">0.622</text>
    <!-- 6: Nando's 0.575 -->
    <text x="135" y="184" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Nando's (v1)</text>
    <rect x="140" y="170" width="316" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="461" y="185" fill="#fff" font-size="11" font-weight="700">0.575</text>
    <!-- 7: Stripe 0.571 -->
    <text x="135" y="214" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Stripe</text>
    <rect x="140" y="200" width="314" height="22" rx="4" fill="url(#bar-grad)"/>
    <text x="459" y="215" fill="#fff" font-size="11" font-weight="700">0.571</text>
    <!-- 8: Linear 0.328 -->
    <text x="135" y="244" text-anchor="end" fill="currentColor" font-size="12" font-weight="600">Linear</text>
    <rect x="140" y="230" width="180" height="22" rx="4" fill="url(#bar-grad)" opacity="0.7"/>
    <text x="325" y="245" fill="#fff" font-size="11" font-weight="700">0.328</text>
  </svg>
</div>
{{template "layout_end" .}}
{{end}}
`

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:             "AI Code Orchestration Research",
		PageType:          "home",
		Categories:        categories,
		NarrativeComplete: narrativeComplete,
		NarrativeTotal:    narrativeTotal,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "home", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/category/")
	slug = strings.TrimSuffix(slug, "/")
	cat, ok := categoryBySlug[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:      cat.Name + " — AI Code Orchestration Research",
		PageType:   "category",
		Categories: categories,
		Category:   cat,
		ActiveCat:  slug,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: cat.Name},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "category", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func experimentHandler(w http.ResponseWriter, r *http.Request) {
	numStr := strings.TrimPrefix(r.URL.Path, "/exp/")
	numStr = strings.TrimSuffix(numStr, "/")
	var exp *Experiment
	var cat *Category
	var num int
	// Try integer lookup first, then string-based for sub-experiments like "172b"
	if n, err := strconv.Atoi(numStr); err == nil {
		if e, ok := expByNum[n]; ok {
			exp = e
			cat = expCategory[n]
			num = n
		}
	}
	if exp == nil {
		if e, ok := expByStr[numStr]; ok {
			exp = e
			cat = expCatByStr[numStr]
			num = exp.NumID
		}
	}
	if exp == nil {
		http.NotFound(w, r)
		return
	}
	// Find this experiment's 1-based index in sorted list
	expIdx := 0
	for i, id := range sortedExpIDs {
		if id == num {
			expIdx = i + 1
			break
		}
	}
	data := PageData{
		Title:      fmt.Sprintf("Exp %s: %s — AI Code Orchestration Research", exp.Num, exp.Focus),
		PageType:   "experiment",
		Categories: categories,
		Category:   cat,
		Experiment: exp,
		ActiveCat:  cat.Slug,
		ActiveExp:  num,
		TotalExps:  len(sortedExpIDs),
		ExpIndex:   expIdx,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: cat.Name, URL: "/category/" + cat.Slug},
			{Label: fmt.Sprintf("Exp %s: %s", exp.Num, exp.Focus)},
		},
	}
	// Load source code if available
	if exp.SourceFile != "" {
		path := scriptsDir + exp.SourceFile
		if content, err := os.ReadFile(path); err == nil {
			data.SourceCode = string(content)
			data.SourceName = exp.SourceFile
			if strings.HasSuffix(exp.SourceFile, ".py") {
				data.SourceLang = "python"
			} else if strings.HasSuffix(exp.SourceFile, ".sh") {
				data.SourceLang = "bash"
			} else if strings.HasSuffix(exp.SourceFile, ".go") {
				data.SourceLang = "go"
			}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "experiment", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func journeyHandler(w http.ResponseWriter, r *http.Request) {
	idxStr := strings.TrimPrefix(r.URL.Path, "/journey/")
	idxStr = strings.TrimSuffix(idxStr, "/")
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 1 || idx > len(sortedExpIDs) {
		http.NotFound(w, r)
		return
	}
	num := sortedExpIDs[idx-1]
	exp := expByNum[num]
	cat := expCategory[num]
	data := PageData{
		Title:       fmt.Sprintf("Chapter %d: %s — Journey", idx, exp.Focus),
		PageType:    "experiment",
		Categories:  categories,
		Category:    cat,
		Experiment:  exp,
		ActiveCat:   cat.Slug,
		ActiveExp:   num,
		TotalExps:   len(sortedExpIDs),
		ExpIndex:    idx,
		JourneyMode: true,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Journey", URL: "/journey/1"},
			{Label: fmt.Sprintf("Chapter %d", idx)},
		},
	}
	// Load source code if available
	if exp.SourceFile != "" {
		path := scriptsDir + exp.SourceFile
		if content, err := os.ReadFile(path); err == nil {
			data.SourceCode = string(content)
			data.SourceName = exp.SourceFile
			if strings.HasSuffix(exp.SourceFile, ".py") {
				data.SourceLang = "python"
			} else if strings.HasSuffix(exp.SourceFile, ".sh") {
				data.SourceLang = "bash"
			} else if strings.HasSuffix(exp.SourceFile, ".go") {
				data.SourceLang = "go"
			}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "experiment", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:          "Discovery Graph — AI Code Orchestration Research",
		PageType:       "discovery",
		Categories:     categories,
		DiscoveryGraph: discoveryGraph,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Discovery Graph"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "discovery", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func timelineHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Research Timeline — AI Code Orchestration Research",
		PageType:   "timeline",
		Categories: categories,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Timeline"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "timeline", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func compareHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/compare/")
	slug = strings.TrimSuffix(slug, "/")

	if slug == "" || slug == "compare" {
		// Overview grid
		data := PageData{
			Title:      "Visual Comparison — AI Code Orchestration Research",
			PageType:   "compare",
			Categories: categories,
			CloneSites: cloneSites,
			Breadcrumbs: []Breadcrumb{{Label: "Home", URL: "/"}, {Label: "Visual Comparison"}},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "compare", data)
		return
	}

	// Detail page
	site, ok := cloneSiteBySlug[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:      site.Name + " Clone — Visual Comparison",
		PageType:   "compare-detail",
		Categories: categories,
		CloneSite:  site,
		CloneSites: cloneSites,
		Breadcrumbs: []Breadcrumb{{Label: "Home", URL: "/"}, {Label: "Visual Comparison", URL: "/compare/"}, {Label: site.Name}},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "compare-detail", data)
}

func apiExperimentsHandler(w http.ResponseWriter, r *http.Request) {
	type ExpJSON struct {
		Num       string  `json:"num"`
		NumID     int     `json:"numID"`
		Focus     string  `json:"focus"`
		Result    string  `json:"result"`
		Cost      string  `json:"cost"`
		Finding   string  `json:"finding"`
		Icon      string  `json:"icon,omitempty"`
		Score     float64 `json:"score,omitempty"`
		HasDetail bool    `json:"hasDetail"`
	}
	type CatJSON struct {
		Slug        string    `json:"slug"`
		Name        string    `json:"name"`
		Experiments []ExpJSON `json:"experiments"`
	}
	var cats []CatJSON
	for _, cat := range categories {
		c := CatJSON{Slug: cat.Slug, Name: cat.Name}
		for _, exp := range cat.Experiments {
			c.Experiments = append(c.Experiments, ExpJSON{
				Num: exp.Num, NumID: exp.NumID, Focus: exp.Focus,
				Result: exp.Result, Cost: exp.Cost, Finding: exp.Finding,
				Icon: exp.Icon, Score: exp.Score, HasDetail: exp.HasDetail,
			})
		}
		cats = append(cats, c)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(cats)
}

func scoreChartHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "SSIM Score Comparison — AI Code Orchestration Research",
		PageType:   "score-chart",
		Categories: categories,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Score Comparison"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "score-chart", data)
}

func roadmapHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Roadmap — AI Code Orchestration Research",
		PageType:   "roadmap",
		Categories: categories,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Roadmap"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "roadmap", data)
}

func pipelineStatsHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Pipeline Stats — AI Code Orchestration Research",
		PageType:   "pipeline-stats",
		Categories: categories,
		Breadcrumbs: []Breadcrumb{
			{Label: "Home", URL: "/"},
			{Label: "Pipeline Stats"},
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "pipeline-stats", data)
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	// Serve screenshot files
	screenshotDirs := map[string]string{
		"/screenshots/airbnb-ref/":    "/tmp/ai-code-orchestration-research/experiments/exp-54-airbnb/",
		"/screenshots/airbnb-nb2/":    "/tmp/ai-code-orchestration-research/experiments/exp-95-ai-images-nano-banana-2/clone-screenshots/",
		"/screenshots/airbnb-gpt5/":   "/tmp/ai-code-orchestration-research/experiments/exp-95-ai-images/clone-screenshots/",
		"/screenshots/airbnb-grey/":   "/tmp/ai-code-orchestration-research/experiments/exp-54-airbnb/clone-screenshots/",
		"/screenshots/each-opus/":     "/tmp/ai-code-orchestration-research/experiments/exp-54-website-clone-opus/font-screenshots/",
		"/screenshots/each-sonnet/":   "/tmp/ai-code-orchestration-research/experiments/exp-54-website-clone-sonnet/",
		"/screenshots/each-llama/":    "/tmp/ai-code-orchestration-research/experiments/exp-54-website-clone-llama4-scout/",
		"/screenshots/each-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-54-website-clone-v2/",
		"/screenshots/nb2-images/":    "/tmp/ai-code-orchestration-research/experiments/exp-95-ai-images-nano-banana-2/",
		"/screenshots/stripe/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-stripe-com/screenshots/",
		"/screenshots/tailwind/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-tailwindcss-com/screenshots/",
		"/screenshots/medium/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-medium-com/screenshots/",
		"/screenshots/linear/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-linear-app/screenshots/",
		"/screenshots/nandos/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-nandos-co-uk/screenshots/",
		"/screenshots/bbc-news/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-bbc-news/screenshots/",
		"/screenshots/producthunt/":   "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-producthunt/screenshots/",
		"/screenshots/github/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-github-repo/screenshots/",
		"/screenshots/spotify/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-spotify-com/screenshots/",
		"/screenshots/airbnb-exp/":    "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-airbnb-experiences/screenshots/",
		"/screenshots/notion/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-notion-so/screenshots/",
		"/screenshots/vercel/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-vercel-com/screenshots/",
		"/screenshots/figma/":         "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-figma-com/screenshots/",
		"/screenshots/slack/":         "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-slack-com/screenshots/",
		"/screenshots/uber-eats/":     "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-uber-eats/screenshots/",
		"/screenshots/duolingo/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-duolingo-com/screenshots/",
		"/screenshots/tesla/":         "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-tesla-com/screenshots/",
		"/screenshots/wise/":          "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-wise-com/screenshots/",
		// Reference screenshots (short names used by first capture batch)
		"/screenshots/ref-stripe/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-stripe/reference/",
		"/screenshots/ref-tailwind/":    "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-tailwind/reference/",
		"/screenshots/ref-medium/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-medium/reference/",
		"/screenshots/ref-linear/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-linear/reference/",
		"/screenshots/ref-nandos/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-nandos/reference/",
		"/screenshots/ref-bbc/":         "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-bbc-news/reference/",
		"/screenshots/ref-producthunt/": "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-producthunt/reference/",
		"/screenshots/ref-github/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-github-repo/reference/",
		// AI generated images (short names used by image gen agents)
		"/screenshots/ai-eachandother/": "/tmp/ai-code-orchestration-research/experiments/exp-54-website-clone-opus/ai-images/",
		"/screenshots/ai-stripe/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-stripe/ai-images/",
		"/screenshots/ai-tailwind/":     "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-tailwind/ai-images/",
		"/screenshots/ai-medium/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-medium/ai-images/",
		"/screenshots/ai-linear/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-linear/ai-images/",
		"/screenshots/ai-nandos/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-nandos/ai-images/",
		"/screenshots/ai-bbc/":          "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-bbc-news/ai-images/",
		"/screenshots/ai-producthunt/":  "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-producthunt/ai-images/",
		"/screenshots/ai-github/":       "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-github-repo/ai-images/",
		"/screenshots/ai-airbnb-exp/":   "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-airbnb-experiences/ai-images/",
		"/screenshots/ai-figma/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-figma-com/ai-images/",
		"/screenshots/ai-slack/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-slack-com/ai-images/",
		"/screenshots/ai-uber-eats/":    "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-uber-eats/ai-images/",
		"/screenshots/ai-duolingo/":     "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-duolingo-com/ai-images/",
		"/screenshots/ai-tesla/":        "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-tesla-com/ai-images/",
		"/screenshots/ai-wise/":         "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-wise-com/ai-images/",
		// Exp 98: Image-guided cloning screenshots
		"/screenshots/exp98-nandos-B/":  "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/nandos-B/screenshots/",
		"/screenshots/exp98-stripe-B/":  "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/stripe-B/screenshots/",
		"/screenshots/exp98-bbc-B/":     "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/bbc-B/screenshots/",
		"/screenshots/exp98-tesla-B/":   "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/tesla-B/screenshots/",
		"/screenshots/exp98-nandos-C/":  "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/nandos-C/screenshots/",
		"/screenshots/exp98-stripe-C/":  "/tmp/ai-code-orchestration-research/experiments/exp-98-image-guided/stripe-C/screenshots/",
		// Bounding box clone experiments (Exp 152-168)
		"/screenshots/exp152-nandos/":   "/tmp/ai-code-orchestration-research/experiments/exp-152-full-blocks/clone-screenshots/",
		"/screenshots/exp152-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-97-clone-nandos/reference/",
		"/screenshots/exp155-nandos/":   "/tmp/ai-code-orchestration-research/experiments/exp-155-logical-sections/clone-screenshots/",
		"/screenshots/exp156-figma/":    "/tmp/ai-code-orchestration-research/experiments/exp-156-clone-figma/clone-screenshots/",
		"/screenshots/exp156-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-156-clone-figma/reference/",
		"/screenshots/exp157-stripe/":   "/tmp/ai-code-orchestration-research/experiments/exp-157-clone-stripe/clone-screenshots/",
		"/screenshots/exp157-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-157-clone-stripe/reference/",
		"/screenshots/exp161-linear/":   "/tmp/ai-code-orchestration-research/experiments/exp-161-clone-linear/clone-screenshots/",
		"/screenshots/exp161-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-161-clone-linear/reference/",
		"/screenshots/exp162-saastools/":"/tmp/ai-code-orchestration-research/experiments/exp-162-clone-to-product/clone-screenshots/",
		"/screenshots/exp162-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-162-clone-to-product/reference/",
		"/screenshots/exp163-transfer/": "/tmp/ai-code-orchestration-research/experiments/exp-163-design-transfer/clone-screenshots/",
		"/screenshots/exp164-previews/": "/tmp/ai-code-orchestration-research/experiments/exp-164-design-transfer-images/clone-screenshots/",
		"/screenshots/exp165-vercel/":   "/tmp/ai-code-orchestration-research/experiments/exp-165-optimized-pipeline/clone-screenshots/",
		"/screenshots/exp165-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-165-optimized-pipeline/reference/",
		"/screenshots/exp168-notion/":   "/tmp/ai-code-orchestration-research/experiments/exp-168-clone-notion/clone-screenshots/",
		"/screenshots/exp168-ref/":      "/tmp/ai-code-orchestration-research/experiments/exp-168-clone-notion/reference/",
	}
	for prefix, dir := range screenshotDirs {
		http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/category/", categoryHandler)
	http.HandleFunc("/exp/", experimentHandler)
	http.HandleFunc("/journey/", journeyHandler)
	http.HandleFunc("/discovery", discoveryHandler)
	http.HandleFunc("/api/experiments.json", apiExperimentsHandler)
	http.HandleFunc("/compare/scores", scoreChartHandler)
	http.HandleFunc("/compare/", compareHandler)
	http.HandleFunc("/timeline", timelineHandler)
	http.HandleFunc("/roadmap", roadmapHandler)
	http.HandleFunc("/pipeline/stats", pipelineStatsHandler)

	addr := ":9094"
	fmt.Printf("Research index server running at http://localhost%s\n", addr)
	fmt.Println("Routes:")
	fmt.Printf("  GET %-30s Home dashboard\n", "/")
	for _, cat := range categories {
		fmt.Printf("  GET %-30s %s\n", "/category/"+cat.Slug, cat.Name)
	}
	fmt.Printf("  GET %-30s Individual experiment\n", "/exp/{num}")
	fmt.Printf("  GET %-30s Journey mode\n", "/journey/{n}")
	fmt.Printf("  GET %-30s Discovery graph\n", "/discovery")
	fmt.Printf("  GET %-30s Roadmap\n", "/roadmap")
	fmt.Printf("  GET %-30s Pipeline Stats\n", "/pipeline/stats")
	fmt.Printf("  %d experiments across %d categories\n", len(sortedExpIDs), len(categories))
	log.Fatal(http.ListenAndServe(addr, nil))
}
