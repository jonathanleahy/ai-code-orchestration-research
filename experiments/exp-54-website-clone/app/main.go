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

		html {
			scroll-behavior: smooth;
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

		h1, h2, h3, h4 {
			font-family: Georgia, "Times New Roman", serif;
			font-weight: normal;
			letter-spacing: -0.02em;
		}

		h1 {
			font-size: 80px;
			line-height: 1.2;
			margin-bottom: 24px;
		}

		h2 {
			font-size: 40px;
			line-height: 1.2;
			margin-bottom: 32px;
		}

		h3 {
			font-size: 24px;
			margin-bottom: 16px;
		}

		h4 {
			font-size: 18px;
			margin-bottom: 16px;
		}

		p {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
		}

		.section-label {
			text-transform: uppercase;
			font-size: 12px;
			letter-spacing: 3px;
			color: #d4727a;
			font-weight: 500;
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
			box-shadow: 0 1px 0 rgba(0,0,0,0.05);
		}

		.nav-logo {
			font-size: 24px;
			font-weight: bold;
			letter-spacing: -0.5px;
			font-family: Georgia, 'Times New Roman', serif;
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
			padding: 150px 40px;
			text-align: left;
			border-bottom: 1px solid #eee;
		}

		.hero-label {
			text-transform: uppercase;
			font-size: 12px;
			letter-spacing: 3px;
			color: #d4727a;
			margin-bottom: 32px;
			display: block;
			font-weight: 500;
		}

		.hero-title {
			font-size: 80px;
			font-family: Georgia, 'Times New Roman', serif;
			line-height: 0.95;
			margin-bottom: 48px;
			max-width: 700px;
		}

		.hero-body {
			font-size: 16px;
			line-height: 1.7;
			max-width: 550px;
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
			letter-spacing: 3px;
			color: #d4727a;
			margin-bottom: 48px;
			font-weight: 500;
		}

		.logos {
			display: flex;
			justify-content: space-between;
			align-items: center;
			flex-wrap: wrap;
			gap: 60px;
		}

		.logo-item {
			font-size: 14px;
			font-weight: 700;
			color: #999;
			text-transform: uppercase;
			letter-spacing: 1.5px;
		}

		/* Services */
		.services {
			padding: 150px 40px;
		}

		.service-grid {
			display: flex;
			flex-direction: column;
			gap: 0;
		}

		.service-item {
			padding: 24px 0;
			border-bottom: 1px solid #eee;
		}

		.service-item:last-child {
			border-bottom: none;
		}

		.service-title {
			font-size: 20px;
			font-family: Georgia, serif;
			margin-bottom: 0;
			display: flex;
			align-items: center;
			gap: 16px;
			transition: color 0.2s;
		}

		.service-title:hover {
			color: #d4727a;
		}

		.service-title::after {
			content: "→";
			color: #d4727a;
			margin-left: auto;
		}

		.service-description {
			display: none;
		}

		/* How Section */
		.how-section {
			padding: 150px 40px;
			background-color: #fafafa;
		}

		.how-section h2 {
			font-size: 40px;
			font-weight: bold;
			margin-bottom: 32px;
		}

		.how-intro {
			font-size: 16px;
			line-height: 1.7;
			max-width: 700px;
			margin-bottom: 80px;
			color: #333;
		}

		.how-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 48px;
		}

		.how-item h3 {
			font-size: 24px;
			margin-bottom: 16px;
		}

		.how-item p {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
		}

		/* CTA */
		.cta-section {
			padding: 150px 40px;
			text-align: center;
			background-color: #d4727a;
			color: #fff;
		}

		.cta-title {
			font-size: 40px;
			font-family: Georgia, serif;
			margin-bottom: 48px;
			color: #fff;
		}

		.cta-link {
			display: inline-flex;
			align-items: center;
			gap: 12px;
			text-decoration: none;
			color: #fff;
			font-size: 16px;
			font-weight: 600;
			border: 2px solid #fff;
			padding: 16px 32px;
			border-radius: 4px;
			transition: all 0.2s;
		}

		.cta-link:hover {
			background-color: #fff;
			color: #d4727a;
		}

		.cta-arrow {
			font-size: 20px;
		}

		/* Case Studies */
		.case-studies {
			padding: 150px 40px;
			border-top: 1px solid #f0f0f0;
		}

		.case-studies-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 40px;
		}

		.case-study {
			padding: 32px;
			border-left: 4px solid #d4727a;
			border-radius: 8px;
		}

		.case-study-label {
			text-transform: uppercase;
			font-size: 12px;
			letter-spacing: 3px;
			color: #d4727a;
			font-weight: 500;
			margin-bottom: 16px;
			display: block;
		}

		.case-study-title {
			font-size: 24px;
			font-family: Georgia, serif;
			margin-bottom: 16px;
			font-weight: bold;
		}

		.case-study-text {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
			margin-bottom: 20px;
		}

		.case-study-link {
			text-decoration: none;
			color: #d4727a;
			font-size: 14px;
			font-weight: 600;
			display: flex;
			align-items: center;
			gap: 8px;
			transition: color 0.2s;
		}

		.case-study-link:hover {
			color: #d4727a;
		}

		.case-study-link::after {
			content: "→";
		}

		/* Blog */
		.blog-section {
			padding: 150px 40px;
			background-color: #fafafa;
			border-top: 1px solid #f0f0f0;
		}

		/* Learn More Section */
		.learn-more-section {
			padding: 150px 40px;
			background-color: #1a1a1a;
			color: #fff;
		}

		.learn-more-section .section-label {
			color: #999;
		}

		.learn-more-section h2 {
			font-size: 40px;
			color: #fff;
			margin-bottom: 48px;
		}

		.learn-more-links {
			display: flex;
			flex-direction: column;
			gap: 24px;
			max-width: 600px;
		}

		.learn-more-link {
			text-decoration: none;
			color: #d4727a;
			font-size: 18px;
			font-family: Georgia, serif;
			display: flex;
			align-items: center;
			gap: 12px;
			transition: color 0.2s;
		}

		.learn-more-link::after {
			content: "→";
			margin-left: auto;
		}

		.learn-more-link:hover {
			color: #d4727a;
		}

		.learn-more-cta {
			margin-top: 48px;
			display: flex;
			align-items: center;
			gap: 24px;
		}

		.learn-more-button {
			display: inline-flex;
			align-items: center;
			gap: 12px;
			text-decoration: none;
			color: #fff;
			font-size: 16px;
			font-weight: 600;
			border: 2px solid #fff;
			padding: 12px 32px;
			border-radius: 4px;
			background-color: transparent;
			transition: all 0.2s;
		}

		.learn-more-button:hover {
			border-color: #d4727a;
			color: #d4727a;
		}

		.learn-more-email {
			color: #d4727a;
			text-decoration: none;
			font-size: 16px;
			transition: color 0.2s;
		}

		.learn-more-email:hover {
			color: #d4727a;
		}

		.blog-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 40px;
		}

		.blog-card {
			background: #fff;
			padding: 24px;
			border-top: 1px solid #eee;
		}

		.blog-date {
			font-size: 12px;
			color: #d4727a;
			margin-bottom: 12px;
			text-transform: uppercase;
			letter-spacing: 3px;
			font-weight: 500;
		}

		.blog-title {
			font-size: 28px;
			font-family: Georgia, serif;
			margin-bottom: 12px;
			font-weight: bold;
		}

		.blog-excerpt {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
			margin-bottom: 16px;
		}

		.blog-read-more {
			text-decoration: none;
			color: #d4727a;
			font-size: 14px;
			font-weight: 600;
			display: inline-block;
			transition: color 0.2s;
		}

		.blog-read-more:hover {
			color: #d4727a;
		}

		/* Benefits Section */
		.benefits-section {
			padding: 120px 40px;
			background-color: #fff;
		}

		.benefits-grid {
			display: grid;
			grid-template-columns: repeat(4, 1fr);
			gap: 40px;
			max-width: 1200px;
			margin: 0 auto;
		}

		.benefit-item {
			text-align: left;
			border-top: 1px solid #ddd;
			padding-top: 24px;
		}

		.benefit-item h3 {
			font-size: 18px;
			font-family: Georgia, serif;
			font-weight: bold;
			margin-bottom: 16px;
		}

		.benefit-item p {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
		}

		/* Approach Section */
		.approach-section {
			padding: 120px 40px;
			background-color: #fafafa;
		}

		.approach-content {
			max-width: 900px;
			margin: 0 auto;
		}

		.approach-section h2 {
			font-size: 40px;
			margin-bottom: 32px;
		}

		.approach-text {
			font-size: 16px;
			line-height: 1.7;
			color: #333;
		}

		/* Learn More Cards Section */
		.learn-more-cards-section {
			padding: 120px 40px;
			background-color: #fff;
		}

		.learn-more-cards-container {
			max-width: 1200px;
			margin: 0 auto;
		}

		.learn-more-cards-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 40px;
			margin-top: 60px;
		}

		.learn-more-card {
			background-color: #f9fafb;
			padding: 40px;
			border-radius: 8px;
			display: flex;
			flex-direction: column;
		}

		.learn-more-card h3 {
			font-size: 24px;
			font-family: Georgia, serif;
			margin-bottom: 24px;
		}

		.learn-more-card-arrow {
			color: #d4727a;
			font-size: 24px;
			margin-top: auto;
		}

		/* Footer */
		footer {
			background-color: #1a1a1a;
			color: #fff;
			padding: 100px 40px 60px;
		}

		footer::before {
			content: "&";
			display: block;
			font-size: 120px;
			line-height: 1;
			margin-bottom: 60px;
			font-family: Georgia, serif;
			color: #d4727a;
		}

		.footer-content {
			display: grid;
			grid-template-columns: repeat(4, 1fr);
			gap: 80px;
			margin-bottom: 80px;
			max-width: 1200px;
			margin-left: auto;
			margin-right: auto;
		}

		.footer-column h4 {
			font-size: 12px;
			font-weight: 500;
			margin-bottom: 24px;
			text-transform: uppercase;
			letter-spacing: 3px;
			color: #d4727a;
		}

		.footer-column ul {
			list-style: none;
		}

		.footer-column li {
			margin-bottom: 16px;
		}

		.footer-column a {
			color: #ccc;
			text-decoration: none;
			font-size: 14px;
			transition: color 0.2s;
		}

		.footer-column a[href^="mailto"] {
			color: #d4727a;
			font-family: 'Courier New', monospace;
		}

		.footer-column a:hover {
			color: #d4727a;
		}

		.footer-column p {
			font-size: 14px;
			line-height: 1.8;
			color: #999;
		}

		.footer-copyright {
			border-top: 1px solid #333;
			padding-top: 40px;
		}

		.footer-copyright-top {
			text-align: center;
			font-size: 13px;
			color: #666;
			margin-bottom: 24px;
		}

		.footer-copyright-top p {
			margin: 4px 0;
			color: #666;
		}

		.footer-copyright-bottom {
			display: flex;
			justify-content: center;
			gap: 32px;
			font-size: 13px;
			color: #666;
			flex-wrap: wrap;
		}

		.footer-copyright-bottom a {
			color: #999;
			text-decoration: none;
			transition: color 0.2s;
		}

		.footer-copyright-bottom a:hover {
			color: #d4727a;
		}

		.social-icons {
			display: flex;
			gap: 16px;
			margin-top: 24px;
		}

		.social-icon {
			display: inline-flex;
			align-items: center;
			justify-content: center;
			width: 40px;
			height: 40px;
			border: 1px solid #555;
			border-radius: 4px;
			text-decoration: none;
			color: #ccc;
			font-size: 18px;
			transition: all 0.2s;
		}

		.social-icon:hover {
			border-color: #d4727a;
			color: #d4727a;
		}

		.office-address {
			margin-top: 24px;
			font-size: 14px;
			line-height: 1.8;
			color: #999;
			font-family: 'Courier New', monospace;
		}

		.office-address p {
			color: #999;
			margin: 0;
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

			.benefits-grid {
				grid-template-columns: 1fr;
			}

			.learn-more-cards-grid {
				grid-template-columns: 1fr;
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
			<li style="margin-left: 16px; font-size: 14px;">EN / DE</li>
		</ul>
	</nav>

	<!-- Hero -->
	<section class="hero">
		<div class="hero-label">A GLOBAL PRODUCT & UX DESIGN COMPANY</div>
		<h1 class="hero-title">Competitive<br>Advantage<br>by <span class="pink">Design</span>.</h1>
		<p class="hero-body">We partner with world-class businesses to transform <span class="pink">customer experience</span> through strategic product design. We're obsessed with understanding <span class="pink">user needs</span> and creating solutions that are <span class="pink">brilliantly easy to use</span>. Our team delivers measurable impact that drives competitive advantage.</p>
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
			<h2>Putting the <span class="pink">customer at the centre</span></h2>
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

	<!-- Benefits Section -->
	<section class="benefits-section">
		<div class="container">
			<div class="benefits-grid">
				<div class="benefit-item">
					<h3>Ship quality products</h3>
					<p>Design with rigor and precision. We ensure every element serves a purpose and delivers measurable business value.</p>
				</div>
				<div class="benefit-item">
					<h3>Create happier users</h3>
					<p>User satisfaction drives loyalty. Our research-backed approach ensures your customers feel heard and valued.</p>
				</div>
				<div class="benefit-item">
					<h3>Reduce project risk</h3>
					<p>Validate early, iterate often. We identify potential issues before they become costly problems in production.</p>
				</div>
				<div class="benefit-item">
					<h3>Deliver results that matter</h3>
					<p>Every design decision is grounded in data and strategy. We track impact and optimize for outcomes that matter.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Our Approach Section -->
	<section class="approach-section">
		<div class="approach-content">
			<h2>We follow a UX process because it works.</h2>
			<p class="approach-text">Our battle-tested process combines strategic thinking with creative excellence. We start with deep customer research to understand the landscape, define clear objectives, and identify opportunities. Through iterative design and continuous validation, we ensure every decision moves closer to measurable business outcomes. This methodical approach has helped world-class companies transform their products and scale their impact.</p>
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
					<span class="case-study-label">Product Transformation</span>
					<h3 class="case-study-title">Google</h3>
					<p class="case-study-text">Redesigned core product flows to improve user engagement and retention. Strategic research informed multi-quarter product roadmap.</p>
					<a href="#case-studies" class="case-study-link">Learn more</a>
				</div>
				<div class="case-study">
					<span class="case-study-label">Design System</span>
					<h3 class="case-study-title">Logitech</h3>
					<p class="case-study-text">Led cross-functional design transformation across connected hardware ecosystem. Built design system supporting global go-to-market.</p>
					<a href="#case-studies" class="case-study-link">Learn more</a>
				</div>
				<div class="case-study">
					<span class="case-study-label">Digital Banking</span>
					<h3 class="case-study-title">BNP Paribas</h3>
					<p class="case-study-text">Modernized digital banking experience serving millions of users. Implemented customer-centric design operating model for enterprise organization.</p>
					<a href="#case-studies" class="case-study-link">Learn more</a>
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
					<a href="#blog" class="blog-read-more">Read more →</a>
				</div>
				<div class="blog-card">
					<div class="blog-date">February 2025</div>
					<h3 class="blog-title">Research-Driven Product Strategy</h3>
					<p class="blog-excerpt">Why deep customer research should inform every strategic decision. We explore proven methodologies that uncover hidden opportunities and reduce execution risk.</p>
					<a href="#blog" class="blog-read-more">Read more →</a>
				</div>
				<div class="blog-card">
					<div class="blog-date">January 2025</div>
					<h3 class="blog-title">Building Design Systems at Scale</h3>
					<p class="blog-excerpt">A practical guide to establishing design systems that scale with your organization. Learn how to implement, govern, and evolve design systems successfully.</p>
					<a href="#blog" class="blog-read-more">Read more →</a>
				</div>
			</div>
		</div>
	</section>

	<!-- Learn More Cards Section -->
	<section class="learn-more-cards-section">
		<div class="learn-more-cards-container">
			<div class="section-label">Learn more about –</div>
			<div class="learn-more-cards-grid">
				<div class="learn-more-card">
					<h3>UX, UI & Product Design</h3>
					<div class="learn-more-card-arrow">→</div>
				</div>
				<div class="learn-more-card">
					<h3>Product Discovery & Research</h3>
					<div class="learn-more-card-arrow">→</div>
				</div>
				<div class="learn-more-card">
					<h3>Fractional Design Leadership</h3>
					<div class="learn-more-card-arrow">→</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Learn More -->
	<section class="learn-more-section">
		<div class="container">
			<h2>Learn more on what we do</h2>
			<div class="learn-more-links">
				<a href="#services" class="learn-more-link">Our Approach</a>
				<a href="#services" class="learn-more-link">UX/UI & Product Design</a>
				<a href="#services" class="learn-more-link">UX Research</a>
				<a href="#services" class="learn-more-link">Design Systems</a>
				<a href="#case-studies" class="learn-more-link">Case Studies</a>
			</div>
			<div class="learn-more-cta">
				<a href="mailto:hello@eachandother.com" class="learn-more-button">Contact us</a>
				<a href="mailto:hello@eachandother.com" class="learn-more-email">hello@eachandother.com</a>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer>
		<div class="footer-content">
			<div class="footer-column">
				<h4>Each&Other</h4>
				<p>Global product & UX design consultancy serving world-class businesses.</p>
				<p style="margin-top: 24px;"><a href="mailto:hello@eachandother.com">hello@eachandother.com</a></p>
				<div class="office-address">
					<p><strong>Dublin, Ireland</strong><br>Unit 1, The Dock<br>Dublin, Ireland</p>
					<p style="margin-top: 16px;"><strong>London, UK</strong><br>42 Greek Street<br>London, UK</p>
				</div>
				<div class="social-icons">
					<a href="https://linkedin.com/company/eachandother" class="social-icon" title="LinkedIn">in</a>
					<a href="https://instagram.com/eachandother" class="social-icon" title="Instagram">📷</a>
				</div>
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
			<div class="footer-copyright-top">
				<p>Copyright © 2026 Each&Other Ltd.</p>
				<p>Registered in Ireland. No. 545982</p>
			</div>
			<div class="footer-copyright-bottom">
				<span>© 2026 Each&Other. All rights reserved.</span>
				<a href="#privacy">Privacy Policy</a>
			</div>
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
