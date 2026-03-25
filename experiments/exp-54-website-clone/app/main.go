package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const html = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Each&Other - Global Product & UX Design Company</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		html, body {
			width: 100%;
			height: 100%;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
			background-color: #ffffff;
			color: #000000;
			line-height: 1.6;
			font-size: 16px;
		}

		.container {
			max-width: 1200px;
			margin: 0 auto;
			padding: 0 40px;
		}

		h1, h2, h3 {
			font-family: Georgia, "Times New Roman", serif;
			font-weight: normal;
			letter-spacing: -0.02em;
		}

		h1 {
			font-size: 64px;
			line-height: 1.2;
			margin-bottom: 24px;
		}

		h2 {
			font-size: 48px;
			line-height: 1.2;
			margin-bottom: 32px;
		}

		h3 {
			font-size: 24px;
			margin-bottom: 16px;
		}

		p {
			font-size: 16px;
			line-height: 1.8;
			color: #333;
		}

		.section-label {
			text-transform: uppercase;
			font-size: 12px;
			letter-spacing: 2px;
			color: #d4727a;
			font-weight: 600;
			margin-bottom: 48px;
		}

		.pink {
			color: #d4727a;
		}

		/* Navigation */
		nav {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 32px 40px;
			border-bottom: 1px solid #f0f0f0;
		}

		.nav-logo {
			font-size: 18px;
			font-weight: 600;
			letter-spacing: -0.5px;
		}

		.nav-links {
			display: flex;
			gap: 32px;
			align-items: center;
			list-style: none;
		}

		.nav-links a {
			text-decoration: none;
			color: #000;
			font-size: 14px;
			transition: color 0.2s;
		}

		.nav-links a:hover {
			color: #d4727a;
		}

		.phone-icon {
			font-size: 18px;
			color: #000;
		}

		/* Hero */
		.hero {
			padding: 120px 40px;
			text-align: center;
		}

		.hero-label {
			text-transform: uppercase;
			font-size: 13px;
			letter-spacing: 2px;
			color: #d4727a;
			margin-bottom: 48px;
			display: inline-block;
		}

		.hero-title {
			font-size: 72px;
			font-family: Georgia, serif;
			line-height: 1.1;
			margin-bottom: 48px;
			max-width: 900px;
			margin-left: auto;
			margin-right: auto;
		}

		.hero-body {
			font-size: 18px;
			line-height: 1.8;
			max-width: 600px;
			margin: 0 auto;
			color: #333;
		}

		/* Logo Bar */
		.logo-bar {
			padding: 100px 40px;
			border-top: 1px solid #f0f0f0;
			border-bottom: 1px solid #f0f0f0;
		}

		.logo-bar-label {
			text-transform: uppercase;
			font-size: 12px;
			letter-spacing: 2px;
			color: #d4727a;
			margin-bottom: 48px;
		}

		.logos {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
			gap: 40px;
			align-items: center;
		}

		.logo-item {
			font-size: 18px;
			font-weight: 500;
			color: #666;
			text-align: center;
		}

		/* Services */
		.services {
			padding: 120px 40px;
		}

		.service-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 60px;
		}

		.service-item {
			padding-bottom: 24px;
		}

		.service-title {
			font-size: 24px;
			font-family: Georgia, serif;
			margin-bottom: 16px;
			display: flex;
			align-items: center;
			gap: 16px;
		}

		.service-title::after {
			content: "→";
			color: #d4727a;
		}

		.service-description {
			font-size: 15px;
			line-height: 1.8;
			color: #555;
		}

		/* How Section */
		.how-section {
			padding: 120px 40px;
			background-color: #fafafa;
		}

		.how-intro {
			font-size: 18px;
			line-height: 1.8;
			max-width: 700px;
			margin-bottom: 80px;
			color: #333;
		}

		.how-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 60px;
		}

		.how-item h3 {
			font-size: 24px;
			margin-bottom: 16px;
		}

		.how-item p {
			font-size: 15px;
			line-height: 1.8;
			color: #555;
		}

		/* CTA */
		.cta-section {
			padding: 120px 40px;
			text-align: center;
		}

		.cta-title {
			font-size: 48px;
			font-family: Georgia, serif;
			margin-bottom: 32px;
		}

		.cta-link {
			display: inline-flex;
			align-items: center;
			gap: 12px;
			text-decoration: none;
			color: #d4727a;
			font-size: 16px;
			font-weight: 600;
			transition: color 0.2s;
		}

		.cta-link:hover {
			color: #c05a62;
		}

		.cta-arrow {
			font-size: 20px;
		}

		/* Case Studies */
		.case-studies {
			padding: 120px 40px;
			border-top: 1px solid #f0f0f0;
		}

		.case-studies-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 40px;
		}

		.case-study {
			padding: 0;
		}

		.case-study-title {
			font-size: 24px;
			font-family: Georgia, serif;
			margin-bottom: 12px;
		}

		.case-study-text {
			font-size: 15px;
			line-height: 1.8;
			color: #555;
		}

		/* Blog */
		.blog-section {
			padding: 120px 40px;
			background-color: #fafafa;
			border-top: 1px solid #f0f0f0;
		}

		.blog-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 40px;
		}

		.blog-card {
			background: #fff;
			padding: 24px;
		}

		.blog-date {
			font-size: 13px;
			color: #999;
			margin-bottom: 12px;
		}

		.blog-title {
			font-size: 20px;
			font-family: Georgia, serif;
			margin-bottom: 12px;
		}

		.blog-excerpt {
			font-size: 14px;
			line-height: 1.7;
			color: #666;
		}

		/* Footer */
		footer {
			background-color: #000;
			color: #fff;
			padding: 80px 40px 40px;
		}

		.footer-content {
			display: grid;
			grid-template-columns: repeat(4, 1fr);
			gap: 60px;
			margin-bottom: 60px;
			max-width: 1200px;
			margin-left: auto;
			margin-right: auto;
		}

		.footer-column h4 {
			font-size: 14px;
			font-weight: 600;
			margin-bottom: 24px;
			text-transform: uppercase;
			letter-spacing: 1px;
		}

		.footer-column ul {
			list-style: none;
		}

		.footer-column li {
			margin-bottom: 12px;
		}

		.footer-column a {
			color: #ccc;
			text-decoration: none;
			font-size: 14px;
			transition: color 0.2s;
		}

		.footer-column a:hover {
			color: #fff;
		}

		.footer-column p {
			font-size: 14px;
			line-height: 1.8;
			color: #aaa;
		}

		.footer-copyright {
			text-align: center;
			border-top: 1px solid #333;
			padding-top: 40px;
			font-size: 13px;
			color: #999;
		}

		@media (max-width: 768px) {
			h1 {
				font-size: 42px;
			}

			h2 {
				font-size: 32px;
			}

			.hero-title {
				font-size: 42px;
			}

			.hero {
				padding: 80px 20px;
			}

			.how-grid {
				grid-template-columns: 1fr;
			}

			.footer-content {
				grid-template-columns: 1fr;
				gap: 40px;
			}

			.nav-links {
				gap: 16px;
			}

			.nav-links a {
				font-size: 12px;
			}
		}
	</style>
</head>
<body>
	<!-- Navigation -->
	<nav>
		<div class="nav-logo">Each&Other</div>
		<ul class="nav-links">
			<li><a href="#services">UX Services</a></li>
			<li><a href="#cta">Head of UX</a></li>
			<li><a href="#services">UX Capability</a></li>
			<li><a href="#case-studies">About Us</a></li>
			<li><a href="#blog">Views</a></li>
			<li><span class="phone-icon">☎</span></li>
		</ul>
	</nav>

	<!-- Hero -->
	<section class="hero">
		<div class="hero-label">A GLOBAL PRODUCT & UX DESIGN COMPANY</div>
		<h1 class="hero-title">Competitive Advantage by <span class="pink">Design</span>.</h1>
		<p class="hero-body">We partner with world-class businesses to transform customer experience through strategic product design and user research. Our team delivers measurable impact that drives competitive advantage.</p>
	</section>

	<!-- Logo Bar -->
	<section class="logo-bar">
		<div class="logo-bar-label">Who we work with –</div>
		<div class="logos">
			<div class="logo-item">Google</div>
			<div class="logo-item">Zurich</div>
			<div class="logo-item">Stripe</div>
			<div class="logo-item">SSE Airtricity</div>
			<div class="logo-item">Coinbase</div>
			<div class="logo-item">BNP Paribas</div>
		</div>
	</section>

	<!-- Services -->
	<section class="services" id="services">
		<div class="section-label">What we do –</div>
		<div class="service-grid">
			<div class="service-item">
				<div class="service-title">UX Research</div>
				<p class="service-description">Deep customer insights through qualitative and quantitative research methods. We uncover user needs and validate opportunities.</p>
			</div>
			<div class="service-item">
				<div class="service-title">Fractional Head of Design</div>
				<p class="service-description">Strategic design leadership for teams building digital products. We shape vision, build process, and mentor growing design teams.</p>
			</div>
			<div class="service-item">
				<div class="service-title">Team Augmentation</div>
				<p class="service-description">Specialized designers and researchers embedded with your team. We scale capability and deliver exceptional output alongside your crew.</p>
			</div>
		</div>
	</section>

	<!-- How We Do It -->
	<section class="how-section">
		<div class="container">
			<div class="section-label">How we do it –</div>
			<p class="how-intro">We believe exceptional product experiences start with a relentless focus on the customer. Every decision we make, every design we craft, and every recommendation we offer is grounded in deep customer understanding and measurable impact.</p>
			<div class="how-grid">
				<div class="how-item">
					<h3>Strategy</h3>
					<p>We work collaboratively to define clear product strategies, identify market opportunities, and establish design foundations that guide successful execution.</p>
				</div>
				<div class="how-item">
					<h3>Design & Delivery</h3>
					<p>From concept through to launch, we deliver high-quality design that balances user needs with business objectives. Iterative, evidence-based refinement is our approach.</p>
				</div>
				<div class="how-item">
					<h3>Long-term Support</h3>
					<p>We partner beyond the initial engagement. Ongoing optimization, team mentoring, and strategic guidance ensure sustained competitive advantage.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA -->
	<section class="cta-section" id="cta">
		<h2 class="cta-title">Let's create something great together.</h2>
		<a href="mailto:hello@eachandother.com" class="cta-link">
			Get in touch
			<span class="cta-arrow">→</span>
		</a>
	</section>

	<!-- Case Studies -->
	<section class="case-studies" id="case-studies">
		<div class="container">
			<div class="section-label">Trusted by world-class businesses –</div>
			<div class="case-studies-grid">
				<div class="case-study">
					<h3 class="case-study-title">Google</h3>
					<p class="case-study-text">Redesigned core product flows to improve user engagement and retention. Strategic research informed multi-quarter product roadmap.</p>
				</div>
				<div class="case-study">
					<h3 class="case-study-title">Logitech</h3>
					<p class="case-study-text">Led cross-functional design transformation across connected hardware ecosystem. Built design system supporting global go-to-market.</p>
				</div>
				<div class="case-study">
					<h3 class="case-study-title">BNP Paribas</h3>
					<p class="case-study-text">Modernized digital banking experience serving millions of users. Implemented customer-centric design operating model for enterprise organization.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Blog -->
	<section class="blog-section" id="blog">
		<div class="container">
			<div class="section-label">Our recent publishing –</div>
			<div class="blog-grid">
				<div class="blog-card">
					<div class="blog-date">March 2025</div>
					<h3 class="blog-title">The Future of Design Leadership</h3>
					<p class="blog-excerpt">How modern design leaders are reshaping organizations and driving competitive advantage through customer-centric practices and cross-functional collaboration.</p>
				</div>
				<div class="blog-card">
					<div class="blog-date">February 2025</div>
					<h3 class="blog-title">Research-Driven Product Strategy</h3>
					<p class="blog-excerpt">Why deep customer research should inform every strategic decision. We explore proven methodologies that uncover hidden opportunities and reduce execution risk.</p>
				</div>
				<div class="blog-card">
					<div class="blog-date">January 2025</div>
					<h3 class="blog-title">Building Design Systems at Scale</h3>
					<p class="blog-excerpt">A practical guide to establishing design systems that scale with your organization. Learn how to implement, govern, and evolve design systems successfully.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer>
		<div class="footer-content">
			<div class="footer-column">
				<h4>Each&Other</h4>
				<p>Global product & UX design consultancy serving world-class businesses.</p>
				<p style="margin-top: 16px;">San Francisco, CA<br>London, UK<br>Berlin, Germany</p>
			</div>
			<div class="footer-column">
				<h4>Services</h4>
				<ul>
					<li><a href="#services">UX Research</a></li>
					<li><a href="#services">Fractional Head of Design</a></li>
					<li><a href="#services">Team Augmentation</a></li>
				</ul>
			</div>
			<div class="footer-column">
				<h4>About</h4>
				<ul>
					<li><a href="#case-studies">Case Studies</a></li>
					<li><a href="#blog">Our Work</a></li>
					<li><a href="#about">About Us</a></li>
					<li><a href="#careers">Careers</a></li>
				</ul>
			</div>
			<div class="footer-column">
				<h4>News</h4>
				<ul>
					<li><a href="#blog">Blog</a></li>
					<li><a href="#newsletter">Newsletter</a></li>
					<li><a href="#social">Social</a></li>
				</ul>
			</div>
		</div>
		<div class="footer-copyright">
			© 2025 Each&Other. All rights reserved.
		</div>
	</footer>
</body>
</html>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
