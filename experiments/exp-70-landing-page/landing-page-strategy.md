# Landing Page Strategy

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FreelanceFlow | CRM for Freelancers</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        }
        body {
            background-color: #f8f9fa;
            color: #333;
            line-height: 1.6;
        }
        .container {
            width: 90%;
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }
        section {
            padding: 80px 0;
        }
        .text-center {
            text-align: center;
        }
        .btn {
            display: inline-block;
            background: #4361ee;
            color: white;
            padding: 15px 30px;
            border-radius: 50px;
            text-decoration: none;
            font-weight: 600;
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
            font-size: 16px;
        }
        .btn:hover {
            background: #3a56d4;
            transform: translateY(-3px);
            box-shadow: 0 10px 20px rgba(0,0,0,0.1);
        }
        .btn-outline {
            background: transparent;
            border: 2px solid #4361ee;
            color: #4361ee;
        }
        .btn-outline:hover {
            background: #4361ee;
            color: white;
        }
        /* Hero Section */
        .hero {
            background: linear-gradient(135deg, #4361ee 0%, #3a0ca3 100%);
            color: white;
            padding: 120px 0;
        }
        .hero h1 {
            font-size: 3.5rem;
            margin-bottom: 20px;
            line-height: 1.2;
        }
        .hero p {
            font-size: 1.2rem;
            margin-bottom: 30px;
            opacity: 0.9;
            max-width: 600px;
            margin-left: auto;
            margin-right: auto;
        }
        .hero-buttons {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 30px;
        }
        /* Features Section */
        .features {
            background: white;
        }
        .section-title {
            font-size: 2.5rem;
            margin-bottom: 15px;
            color: #2b2d42;
        }
        .section-subtitle {
            font-size: 1.1rem;
            color: #666;
            max-width: 600px;
            margin: 0 auto 50px;
        }
        .features-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
        }
        .feature-card {
            text-align: center;
            padding: 30px;
            border-radius: 10px;
            transition: all 0.3s ease;
        }
        .feature-card:hover {
            transform: translateY(-10px);
            box-shadow: 0 15px 30px rgba(0,0,0,0.1);
        }
        .feature-icon {
            font-size: 3rem;
            color: #4361ee;
            margin-bottom: 20px;
        }
        .feature-card h3 {
            margin-bottom: 15px;
            color: #2b2d42;
        }
        /* Pricing Section */
        .pricing {
            background: #f8f9fa;
        }
        .pricing-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
        }
        .pricing-card {
            background: white;
            border-radius: 10px;
            padding: 40px 30px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.05);
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }
        .pricing-card:hover {
            transform: translateY(-10px);
            box-shadow: 0 15px 30px rgba(0,0,0,0.1);
        }
        .pricing-card.popular {
            border: 2px solid #4361ee;
            transform: scale(1.05);
        }
        .popular-tag {
            position: absolute;
            top: 20px;
            right: -30px;
            background: #4361ee;
            color: white;
            padding: 5px 30px;
            transform: rotate(45deg);
            font-size: 0.8rem;
            font-weight: bold;
        }
        .price {
            font-size: 3rem;
            font-weight: 700;
            color: #2b2d42;
            margin: 20px 0;
        }
        .price span {
            font-size: 1rem;
            color: #666;
        }
        .pricing-features {
            list-style: none;
            margin: 30px 0;
            text-align: left;
        }
        .pricing-features li {
            padding: 10px 0;
            border-bottom: 1px solid #eee;
        }
        .pricing-features li:last-child {
            border-bottom: none;
        }
        /* Testimonials */
        .testimonials {
            background: white;
        }
        .testimonials-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
        }
        .testimonial-card {
            background: #f8f9fa;
            padding: 30px;
            border-radius: 10px;
            position: relative;
        }
        .testimonial-card:before {
            content: """;
            position: absolute;
            top: 20px;
            left: 20px;
            font-size: 5rem;
            color: #4361ee;
            opacity: 0.2;
            font-family: Georgia, serif;
        }
        .testimonial-content {
            margin-top: 20px;
            font-style: italic;
            color: #555;
        }
        .testimonial-author {
            display: flex;
            align-items: center;
            margin-top: 20px;
        }
        .author-avatar {
            width: 50px;
            height: 50px;
            border-radius: 50%;
            background: #4361ee;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            margin-right: 15px;
        }
        .author-info h4 {
            color: #2b2d42;
        }
        .author-info p {
            color: #666;
            font-size: 0.9rem;
        }
        /* FAQ */
        .faq {
            background: #f8f9fa;
        }
        .faq-container {
            max-width: 800px;
            margin: 0 auto;
        }
        .faq-item {
            margin-bottom: 20px;
            border-radius: 10px;
            overflow: hidden;
        }
        .faq-question {
            background: white;
            padding: 20px;
            font-weight: 600;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .faq-answer {
            background: white;
            padding: 20px;
            display: none;
        }
        .faq-answer.show {
            display: block;
        }
        /* Footer */
        .footer {
            background: #2b2d42;
            color: white;
            padding: 60px 0 30px;
        }
        .footer-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 30px;
            margin-bottom: 40px;
        }
        .footer-column h3 {
            margin-bottom: 20px;
            font-size: 1.3rem;
        }
        .footer-column ul {
            list-style: none;
        }
        .footer-column ul li {
            margin-bottom: 10px;
        }
        .footer-column ul li a {
            color: #ddd;
            text-decoration: none;
            transition: color 0.3s ease;
        }
        .footer-column ul li a:hover {
            color: #4361ee;
        }
        .copyright {
            text-align: center;
            padding-top: 30px;
            border-top: 1px solid #4a4c6a;
            color: #aaa;
            font-size: 0.9rem;
        }
        /* Responsive */
        @media (max-width: 768px) {
            .hero h1 {
                font-size: 2.5rem;
            }
            .hero-buttons {
                flex-direction: column;
                align-items: center;
            }
            .section-title {
                font-size: 2rem;
            }
        }
    </style>
</head>
<body>
    <!-- Hero Section -->
    <section class="hero">
        <div class="container text-center">
            <h1>Manage Your Freelance Business Like a Pro</h1>
            <p>FreelanceFlow helps you organize clients, track projects, and manage invoices - all in one simple platform designed for freelancers.</p>
            <div class="hero-buttons">
                <a href="#" class="btn">Start Free Trial</a>
                <a href="#" class="btn btn-outline">View Demo</a>
            </div>
        </div>
    </section>

    <!-- Features Section -->
    <section class="features">
        <div class="container">
            <h2 class="section-title text-center">Powerful Features for Freelancers</h2>
            <p class="section-subtitle text-center">Everything you need to run your freelance business efficiently</p>
            
            <div class="features-grid">
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-users"></i>
                    </div>
                    <h3>Client Management</h3>
                    <p>Store all client information, communication history, and project details in one place.</p>
                </div>
                
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-tasks"></i>
                    </div>
                    <h3>Project Tracking</h3>
                    <p>Monitor project progress, deadlines, and deliverables with our intuitive dashboard.</p>
                </div>
                
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-file-invoice-dollar"></i>
                    </div>
                    <h3>Smart Invoicing</h3>
                    <p>Create professional invoices automatically and track payment status in real-time.</p>
                </div>
                
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-calendar-alt"></i>
                    </div>
                    <h3>Time Tracking</h3>
                    <p>Log hours effortlessly and generate detailed reports for your clients.</p>
                </div>
                
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-chart-line"></i>
                    </div>
                    <h3>Financial Insights</h3>
                    <p>Get clear analytics on your income, expenses, and business performance.</p>
                </div>
                
                <div class="feature-card">
                    <div class="feature-icon">
                        <i class="fas fa-mobile-alt"></i>
                    </div>
                    <h3>Mobile App</h3>
                    <p>Access your CRM on the go with our fully functional mobile application.</p>
                </div>
            </div>
        </div>
    </section>

    <!-- Pricing Section -->
    <section class="pricing">
        <div class="container">
            <h2 class="section-title text-center">Simple, Transparent Pricing</h2>
            <p class="section-subtitle text-center">Choose the plan that works best for your freelance business</p>
            
            <div class="pricing-grid">
                <div class="pricing-card">
                    <h3>Starter</h3>
                    <div class="price">$0<span>/month</span></div>
                    <p>Perfect for getting started</p>
                    <ul class="pricing-features">
                        <li>Up to 5 clients</li>
                        <li>Basic project tracking</li>
                        <li>Simple invoicing</li>
                        <li>Email support</li>
                    </ul>
                    <a href="#" class="btn btn-outline">Get Started</a>
                </div>
                
                <div class="pricing-card popular">
                    <div class="popular-tag">POPULAR</div>
                    <h3>Pro</h3>
                    <div class="price">$15<span>/month</span></div>
                    <p>For growing freelancers</p>
                    <ul class="pricing-features">
                        <li>Unlimited clients</li>
                        <li>Advanced project tracking</li>
                        <li>Smart invoicing & payments</li>
                        <li>Time tracking</li>
                        <li>Financial reports</li>
                        <li>Mobile app access</li>
                    </ul>
                    <a href="#" class="btn">Start Free Trial</a>
                </div>
                
                <div class="pricing-card">
                    <h3>Business</h3>
                    <div class="price">$49<span>/month</span></div>
                    <p>For teams and agencies</p>
                    <ul class="pricing-features">
                        <li>All Pro features</li>
                        <li>Team collaboration</li>
                        <li>Multi-user accounts</li>
                        <li>Advanced analytics</li>
                        <li>Priority support</li>
                        <li>Custom branding</li>
                    </ul>
                    <a href="#" class="btn btn-outline">Contact Sales</a>
                </div>
            </div>
        </div>
    </section>

    <!-- Testimonials Section -->
    <section class="testimonials">
        <div class="container">
            <h2 class="section-title text-center">Trusted by Thousands of Freelancers</h2>
            <p class="section-subtitle text-center">See what our users have to say about FreelanceFlow</p>
            
            <div class="testimonials-grid">
                <div class="testimonial-card">
                    <div class="testimonial-content">
                        FreelanceFlow has completely transformed how I manage my client relationships. I've saved over 10 hours per week on administrative tasks and my client retention has improved by 40%.
                    </div>
                    <div class="testimonial-author">
                        <div class="author-avatar">MJ</div>
                        <div class="author-info">
                            <h4>Michael Johnson</h4>
                            <p>UI/UX Designer</p>
                        </div>
                    </div>
                </div>
                
                <div class="testimonial-card">
                    <div class="testimonial-content">
                        As a freelance writer with multiple clients, keeping track of deadlines and payments was a nightmare. FreelanceFlow's invoicing system alone has saved me countless hours and improved my cash flow significantly.
                    </div>
                    <div class="testimonial-author">
                        <div class="author-avatar">SR</div>
                        <div class="author-info">
                            <h4>Sarah Rodriguez</h4>
                            <p>Content Writer</p>
                        </div>
                    </div>
                </div>
                
                <div class="testimonial-card">
                    <div class="testimonial-content">
                        The time tracking feature alone is worth the price. I can now accurately bill my clients and see exactly where my time goes. The mobile app is a game-changer for field work.
                    </div>
                    <div class="testimonial-author">
                        <div class="author-avatar">DT</div>
                        <div class="author-info">
                            <h4>David Thompson</h4>
                            <p>Web Developer</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- FAQ Section -->
    <section class="faq">
        <div class="container">
            <h2 class="section-title text-center">Frequently Asked Questions</h2>
            <p class="section-subtitle text-center">Everything you need to know about FreelanceFlow</p>
            
            <div class="faq-container">
                <div class="faq-item">
                    <div class="faq-question">
                        How does the free trial work?
                        <i class="fas fa-chevron-down"></i>
                    </div>
                    <div class="faq-answer">
                        <p>Our 14-day free trial gives you full access to all Pro features. No credit card required to start. You can cancel anytime during the trial period.</p>
                    </div>
                </div>
                
                <div class="faq-item">
                    <div class="faq-question">
                        Can I upgrade or downgrade my plan?
                        <i class="fas fa-chevron-down"></i>
                    </div>
                    <div class="faq-answer">
                        <p>Yes, you can upgrade or downgrade your plan at any time. Changes take effect immediately, and you'll be charged or credited accordingly.</p>
                    </div>
                </div>
                
                <div class="faq-item">
                    <div class="faq-question">
                        Is my data secure?
                        <i class="fas fa-chevron-down"></i>
                    </div>
                    <div class="faq-answer">
                        <p>Absolutely. We use industry-standard encryption and security measures to protect your data. Your information is stored securely in the cloud with regular backups.</p>
                    </