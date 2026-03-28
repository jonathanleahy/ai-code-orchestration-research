package main

import (
	"fmt"
	"html/template"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// Data model
// ---------------------------------------------------------------------------

type Experiment struct {
	Num         string // "01", "40–42", etc.
	NumID       int    // primary number for linking (1, 40, etc.)
	Focus       string
	Result      string
	Cost        string
	Finding     string
	Category    string          // extra column for business-pipeline / model-comparison
	HasDetail   bool            // whether /exp/N should resolve
	Icon        string          // 1-3 emoji icons identifying the experiment
	Description []template.HTML // detailed content (legacy, still rendered if present)
	SourceFile  string          // path to script relative to scriptsDir
	// Structured sections for experiment detail pages
	Why     template.HTML   // Why was this experiment run?
	What    template.HTML   // Setup — models, inputs, configuration
	How     template.HTML   // Approach — steps taken
	Impact  template.HTML   // What changed in the pipeline because of this
	Related []string        // experiment numbers that relate (e.g. "Exp 2", "Exp 17")
	// Computed fields (populated in init)
	CostFloat    float64 // parsed from Cost string
	IsFailure    bool    // true if Result contains FAIL or 0%
	WordCount    int     // word count of narrative sections
	ReadingTime  int     // WordCount / 200, min 1
	Score        float64 // SSIM score for cloning experiments
	Thumbnail    string  // URL path to screenshot thumbnail
	PrimaryColor string  // brand color for visual cards
	CloneShot    string  // URL path to clone desktop screenshot
	RefShot      string  // URL path to reference desktop screenshot
}

const scriptsDir = "/home/jon/work/dark-factory/scripts/dev-spike-v3/"

type Category struct {
	Slug        string
	Name        string
	Icon        string // emoji icon for the category
	ExpRange    string
	Narrative   []template.HTML
	KeyInsight  template.HTML
	Experiments []Experiment
	TableType   string // "standard", "website-clone", "graphql", "model-comparison", "business-pipeline"
	MaxCost     float64 // computed: max CostFloat in this category (for cost bars)
}

type PageData struct {
	Title          string
	PageType       string // "home", "category", "experiment", "discovery"
	Categories     []Category
	Category       *Category
	Experiment     *Experiment
	Breadcrumbs    []Breadcrumb
	ActiveCat      string // slug of expanded sidebar category
	ActiveExp      int    // experiment number if on exp page
	SourceCode     string // source code content for experiment pages
	SourceLang     string // "python", "bash", "go"
	SourceName     string // filename for display
	TotalExps          int    // total routable experiments
	NarrativeComplete  int    // experiments with full Why/What/How/Impact
	NarrativeTotal     int    // experiments with detail pages
	JourneyMode        bool   // true if rendering in journey mode
	ExpIndex       int    // 1-based position in sorted experiments
	DiscoveryGraph template.HTML   // mermaid graph for discovery page
	CloneSites     []CloneSite    // for compare pages
	CloneSite      *CloneSite     // for compare detail page
}

type Breadcrumb struct {
	Label string
	URL   string
}

type CloneSite struct {
	Slug         string
	Name         string
	Category     string
	PrimaryColor string
	Lines        int
	Port         int
	Iterations   int
	Method       string // "Text Description", "Screenshot-Guided", "Hybrid"
	RefPrefix    string   // screenshot URL prefix for original
	ClonePrefix  string   // screenshot URL prefix for clone
	AIPrefix     string   // screenshot URL prefix for AI images
	AIImages     []string // filenames in AI images dir
}

var cloneSites = []CloneSite{
	{Slug: "airbnb", Name: "Airbnb", Category: "E-commerce listing", PrimaryColor: "#FF385C", Lines: 757, Port: 8095, Iterations: 10, RefPrefix: "/screenshots/airbnb-ref", ClonePrefix: "/screenshots/airbnb-nb2", AIPrefix: "/screenshots/nb2-images", AIImages: []string{"hero.png", "bedroom1.png", "kitchen.png", "bedroom2.png", "bathroom.png", "neighbourhood.png", "host-avatar.png"}},
	{Slug: "eachandother", Name: "Each&Other", Category: "Agency portfolio", PrimaryColor: "#d4727a", Lines: 1024, Port: 8094, Iterations: 8, ClonePrefix: "/screenshots/each-opus", AIPrefix: "/screenshots/ai-eachandother", AIImages: []string{"hero.png", "team.png"}},
	{Slug: "stripe", Name: "Stripe", Category: "SaaS landing", PrimaryColor: "#635BFF", Lines: 948, Port: 8101, Iterations: 3, RefPrefix: "/screenshots/ref-stripe", ClonePrefix: "/screenshots/stripe", AIPrefix: "/screenshots/ai-stripe", AIImages: []string{"hero.png", "dashboard.png"}},
	{Slug: "tailwind", Name: "Tailwind CSS", Category: "Developer docs", PrimaryColor: "#38BDF8", Lines: 616, Port: 8102, Iterations: 1, RefPrefix: "/screenshots/ref-tailwind", ClonePrefix: "/screenshots/tailwind", AIPrefix: "/screenshots/ai-tailwind", AIImages: []string{"hero.png", "preview.png"}},
	{Slug: "medium", Name: "Medium", Category: "Blog article", PrimaryColor: "#1A8917", Lines: 1137, Port: 8103, Iterations: 1, RefPrefix: "/screenshots/ref-medium", ClonePrefix: "/screenshots/medium", AIPrefix: "/screenshots/ai-medium", AIImages: []string{"hero.png", "author.png"}},
	{Slug: "linear", Name: "Linear", Category: "Dashboard SaaS", PrimaryColor: "#5E6AD2", Lines: 725, Port: 8104, Iterations: 1, RefPrefix: "/screenshots/ref-linear", ClonePrefix: "/screenshots/linear", AIPrefix: "/screenshots/ai-linear", AIImages: []string{"hero.png", "feature.png"}},
	{Slug: "nandos", Name: "Nando's", Category: "Restaurant", PrimaryColor: "#ED1C24", Lines: 780, Port: 8105, Iterations: 11, RefPrefix: "/screenshots/ref-nandos", ClonePrefix: "/screenshots/nandos", AIPrefix: "/screenshots/ai-nandos", AIImages: []string{"hero.png", "interior.png"}},
	{Slug: "bbc-news", Name: "BBC News", Category: "News media", PrimaryColor: "#BB1919", Lines: 339, Port: 8106, Iterations: 1, RefPrefix: "/screenshots/ref-bbc", ClonePrefix: "/screenshots/bbc-news", AIPrefix: "/screenshots/ai-bbc", AIImages: []string{"hero.png", "reporter.png"}},
	{Slug: "producthunt", Name: "Product Hunt", Category: "Marketplace", PrimaryColor: "#DA552F", Lines: 503, Port: 8107, Iterations: 1, RefPrefix: "/screenshots/ref-producthunt", ClonePrefix: "/screenshots/producthunt", AIPrefix: "/screenshots/ai-producthunt", AIImages: []string{"hero.png", "avatar.png"}},
	{Slug: "github", Name: "GitHub", Category: "Developer tool", PrimaryColor: "#1f6feb", Lines: 572, Port: 8108, Iterations: 1, RefPrefix: "/screenshots/ref-github", ClonePrefix: "/screenshots/github", AIPrefix: "/screenshots/ai-github", AIImages: []string{"hero.png", "code.png"}},
	{Slug: "spotify", Name: "Spotify", Category: "Music streaming", PrimaryColor: "#1DB954", Lines: 485, Port: 8111, Iterations: 1, ClonePrefix: "/screenshots/spotify"},
	{Slug: "airbnb-exp", Name: "Airbnb Experiences", Category: "Experiences", PrimaryColor: "#FF385C", Lines: 389, Port: 8112, Iterations: 1, ClonePrefix: "/screenshots/airbnb-exp", AIPrefix: "/screenshots/ai-airbnb-exp", AIImages: []string{"hero.png", "activity.png"}},
	{Slug: "notion", Name: "Notion", Category: "Productivity", PrimaryColor: "#000000", Lines: 1308, Port: 8113, Iterations: 1, ClonePrefix: "/screenshots/notion"},
	{Slug: "vercel", Name: "Vercel", Category: "Developer platform", PrimaryColor: "#000000", Lines: 502, Port: 8114, Iterations: 1, ClonePrefix: "/screenshots/vercel"},
	{Slug: "figma", Name: "Figma", Category: "Design tool", PrimaryColor: "#F24E1E", Lines: 544, Port: 8115, Iterations: 1, ClonePrefix: "/screenshots/figma", AIPrefix: "/screenshots/ai-figma", AIImages: []string{"hero.png", "collaboration.png"}},
	{Slug: "slack", Name: "Slack", Category: "Communication", PrimaryColor: "#4A154B", Lines: 608, Port: 8116, Iterations: 1, ClonePrefix: "/screenshots/slack", AIPrefix: "/screenshots/ai-slack", AIImages: []string{"hero.png", "team.png"}},
	{Slug: "uber-eats", Name: "Uber Eats", Category: "Food delivery", PrimaryColor: "#06C167", Lines: 275, Port: 8117, Iterations: 1, ClonePrefix: "/screenshots/uber-eats", AIPrefix: "/screenshots/ai-uber-eats", AIImages: []string{"hero.png", "delivery.png"}},
	{Slug: "duolingo", Name: "Duolingo", Category: "Education", PrimaryColor: "#58CC02", Lines: 459, Port: 8118, Iterations: 1, ClonePrefix: "/screenshots/duolingo", AIPrefix: "/screenshots/ai-duolingo", AIImages: []string{"hero.png", "mascot.png"}},
	{Slug: "tesla", Name: "Tesla", Category: "Automotive", PrimaryColor: "#CC0000", Lines: 223, Port: 8119, Iterations: 1, ClonePrefix: "/screenshots/tesla", AIPrefix: "/screenshots/ai-tesla", AIImages: []string{"hero.png", "interior.png"}},
	{Slug: "wise", Name: "Wise", Category: "Fintech", PrimaryColor: "#9FE870", Lines: 505, Port: 8120, Iterations: 1, ClonePrefix: "/screenshots/wise", AIPrefix: "/screenshots/ai-wise", AIImages: []string{"hero.png", "card.png"}},
	// Screenshot-guided clones (Exp 98)
	{Slug: "nandos-screenshot", Name: "Nando's", Category: "Restaurant", PrimaryColor: "#ED1C24", Lines: 786, Port: 9201, Iterations: 1, Method: "Screenshot-Guided", RefPrefix: "/screenshots/ref-nandos", ClonePrefix: "/screenshots/exp98-nandos-B", AIImages: []string{"hero.png", "interior.png", "food-grid.png", "delivery.png"}},
	{Slug: "stripe-screenshot", Name: "Stripe", Category: "SaaS landing", PrimaryColor: "#635BFF", Lines: 491, Port: 9202, Iterations: 1, Method: "Screenshot-Guided", RefPrefix: "/screenshots/ref-stripe", ClonePrefix: "/screenshots/exp98-stripe-B", AIImages: []string{"hero.png", "dashboard.png", "developers.png"}},
	{Slug: "bbc-screenshot", Name: "BBC News", Category: "News media", PrimaryColor: "#BB1919", Lines: 588, Port: 9203, Iterations: 1, Method: "Screenshot-Guided", RefPrefix: "/screenshots/ref-bbc", ClonePrefix: "/screenshots/exp98-bbc-B", AIImages: []string{"hero.png", "parliament.png", "reporter.png", "weather.png"}},
	{Slug: "tesla-screenshot", Name: "Tesla", Category: "Automotive", PrimaryColor: "#CC0000", Lines: 205, Port: 9204, Iterations: 1, Method: "Screenshot-Guided", ClonePrefix: "/screenshots/exp98-tesla-B", AIImages: []string{"hero.png", "interior.png", "charging.png", "factory.png"}},
	{Slug: "nandos-hybrid", Name: "Nando's", Category: "Restaurant", PrimaryColor: "#ED1C24", Lines: 544, Port: 9211, Iterations: 4, Method: "Hybrid", RefPrefix: "/screenshots/ref-nandos", ClonePrefix: "/screenshots/exp98-nandos-C", AIImages: []string{"hero.png", "menu-items.png", "spice-rack.png", "team.png"}},
	{Slug: "stripe-hybrid", Name: "Stripe", Category: "SaaS landing", PrimaryColor: "#635BFF", Lines: 491, Port: 9212, Iterations: 1, Method: "Hybrid", RefPrefix: "/screenshots/ref-stripe", ClonePrefix: "/screenshots/exp98-stripe-C", AIImages: []string{"hero.png", "analytics.png", "cards.png", "mobile.png"}},
}

var cloneSiteBySlug = map[string]*CloneSite{}

func init() {
	for i := range cloneSites {
		cloneSiteBySlug[cloneSites[i].Slug] = &cloneSites[i]
	}
}

// ---------------------------------------------------------------------------
// All research data
// ---------------------------------------------------------------------------

var categories = []Category{
	{
		Slug:     "the-pipeline",
		Name:     "The Pipeline",
		Icon:     "\U0001f527",
		ExpRange: "Brief to Deployed Product — 13 Steps, $0.49 Total",
		Narrative: []template.HTML{
			`This pipeline takes an 8-word product brief and produces a running, tested, reviewed, deployed application. Every step was validated through 116+ experiments. Total cost: <strong>$0.49</strong> with cheap models, <strong>$2.75</strong> with premium (Opus).`,
			`The pipeline is the central output of the research. Each step was discovered and validated independently, then assembled into a coherent flow. <strong>The order matters</strong>: persona discovery before specification prevents missing features; dev review before code prevents architectural mistakes; progressive enhancement after initial build prevents regressions; Playwright after all fixes catches integration bugs nothing else finds.`,
			`<strong>The key insight is layered quality gates.</strong> No single technique catches everything. Unit tests miss browser bugs. Reviews miss race conditions. Static analysis misses UX problems. But stacked together &mdash; reviewers, then code, then post-code review, then fix, then Playwright, then security &mdash; the pipeline catches virtually every class of defect for under $1.`,
			`The pipeline has three phases: <strong>Design</strong> (steps 1&ndash;5, $0.037) turns a brief into a reviewed, typed specification. <strong>Build</strong> (steps 6&ndash;10, $0.41) produces working, tested code with progressive enhancement. <strong>Verify</strong> (steps 11&ndash;13, $0.03) catches browser bugs, security holes, and deploys to production.`,
		},
		KeyInsight: `The auto-fix pipeline (goimports + gofmt + fix-address-of-const) recovers 40&ndash;60% of build failures for free. Each quality gate catches a different class of defect — stacking them is what makes the pipeline work.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "Step 1", NumID: 901, Focus: "Persona Discovery", Result: "Exp 23", Cost: "$0.01", Finding: "4 personas find features the brief missed", HasDetail: true,
				Why: `A product brief is always incomplete. Ours said "CRM for freelancers with invoicing" &mdash; eight words. A human product manager would spend weeks interviewing users to find the gaps. We had $0.01 and 30 seconds. The question was whether AI-generated personas could do what real user interviews do: surface the features the brief author forgot to mention.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["8-Word Brief"] --> B["Generate 4 Personas"]
    B --> C["Interview Brief"]
    C --> D{"Accept / Reject?"}
    D -->|Approve| E["Approved Features"]
    D -->|Reject + Objections| C
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#d4edda,stroke:#28a745
    click C "/exp/23"
</pre>
We generate 4 realistic user personas representing different archetypes &mdash; power user, casual user, admin, and new user. Each persona evaluates the product brief from their perspective and identifies what is missing. The input is the raw brief; the output is 4 detailed personas with needs, pain points, and concrete feature demands.`),
				How: `Each persona is prompted to <em>interview</em> the brief: What would I need from this product? What is obviously missing? Would I actually use this? The personas don't just list features &mdash; they approve or reject the concept with specific objections. In <a href="/exp/23">Experiment 23</a>, two personas approved the CRM concept and two rejected it. The rejections were more valuable than the approvals &mdash; they surfaced real requirements gaps that the brief never mentioned.`,
				Impact: `Personas found "recurring invoices" was demanded by 3 of 4 personas despite not appearing in the original brief. After incorporating their feedback, all four approved. This became Step 1 of the pipeline because <strong>every subsequent step builds on the feature list</strong> &mdash; if the list is wrong, the wireframes are wrong, the spec is wrong, the code is wrong. For $0.01, we eliminate the most expensive class of bug: building the wrong product.`,
				Related: []string{"23"},
			},
			{Num: "Step 2", NumID: 902, Focus: "MVP Synthesis", Result: "Exp 36", Cost: "$0.005", Finding: "Constrained simplicity filters scope", HasDetail: true,
				Why: `After <a href="/exp/901">Step 1</a>, we have persona demands &mdash; sometimes dozens of features. Shipping all of them in an MVP is expensive and slow. But we discovered in Experiment 34 that a naive "simplicity agent" cuts features the personas demanded. It would see "recurring invoices" and say "too complex, skip it." That is the wrong kind of simplification &mdash; removing scope instead of reducing implementation cost.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Persona Demands"] --> B["Simplicity Agent"]
    B --> C{"Per Feature"}
    C -->|Complex| D["SIMPLIFY How"]
    C -->|Essential| E["KEEP As-Is"]
    C -->|Out of Scope| F["CUT"]
    D --> G["MVP Feature List"]
    E --> G
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style G fill:#d4edda,stroke:#28a745
    style F fill:#fce4ec,stroke:#d4727a
    click B "/exp/36"
</pre>
The constrained simplicity agent synthesises persona demands into a prioritised feature list. Its key instruction: <em>&ldquo;simplify HOW features are implemented, never WHETHER they are included.&rdquo;</em> Input is persona interviews plus the original brief. Output is a feature list split into &ldquo;must have&rdquo; and &ldquo;nice to have,&rdquo; with implementation notes favouring the cheapest approach that still works.`),
				How: `The agent examines each persona demand and asks: can we deliver this with simpler technology? In-memory storage instead of PostgreSQL. Embedded HTML instead of a SPA framework. Simple form validation instead of a validation library. The feature stays; only the implementation complexity drops. <a href="/exp/36">Experiment 36</a> ran this head-to-head against a domain-expert-only approach to measure the difference.`,
				Impact: `Constrained simplicity produced an app with <strong>all features present but 60% cheaper</strong> than the domain-expert-only approach. Every feature worked &mdash; they were just implemented more simply. This step is the reason the pipeline produces full-featured apps for under $0.50 instead of $2+. It feeds directly into <a href="/exp/903">Step 3 (Wireframes)</a>, where the simplified feature list becomes concrete screens.`,
				Related: []string{"34", "36"},
			},
			{Num: "Step 3", NumID: 903, Focus: "Screen Wireframes", Result: "Exp 22, 24", Cost: "$0.008", Finding: "Text wireframes for every screen", HasDetail: true,
				Why: `We had a feature list. We asked the AI to "build a CRM with these features." It produced 493 lines of code with four navigation patterns, two CSS frameworks loaded inline, and a settings page nobody asked for. The AI was not being creative &mdash; it was guessing, because "build a CRM" is ambiguous. Every ambiguity becomes a coin flip, and coin flips compound. We needed to remove the guesswork before code generation.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Feature List"] --> B["Wireframe Generator"]
    B --> C["Text Wireframes"]
    C --> D["Guided Code Gen"]
    D --> E["106 Lines Focused"]
    F["No Wireframes"] --> G["Unconstrained Gen"]
    G --> H["493 Lines Bloated"]
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style E fill:#d4edda,stroke:#28a745
    style H fill:#fce4ec,stroke:#d4727a
    click B "/exp/22"
    click D "/exp/24"
</pre>
Each screen is described as a text wireframe: what elements appear, their layout, user interactions, and navigation paths. Input is the feature list plus personas. Output is a set of text wireframes describing every screen the user will see &mdash; turning an abstract feature list into a concrete UI specification the code generation step can follow precisely.`),
				How: template.HTML("<a href=\"/exp/22\">Experiment 22</a> developed the brief-to-journeys-to-screens pipeline: the AI maps user journeys first, then derives screens from those journeys. <a href=\"/exp/24\">Experiment 24</a> tested wireframe-guided code generation head-to-head against unconstrained generation on the same feature set, measuring lines of code and feature completeness."),
				Impact: `The wireframe version produced <strong>106 lines versus 493 lines</strong> for unconstrained generation &mdash; less than a quarter of the code, with the same feature set. Wireframes constrain the AI's tendency to over-engineer. This is why we generate wireframes before code, not after: they act as guardrails that keep <a href="/exp/907">Step 7 (HTTP Server)</a> focused on exactly what the user needs to see.`,
				Related: []string{"22", "24"},
			},
			{Num: "Step 4", NumID: 904, Focus: "Dev Review Gate", Result: "Exp 25, 32", Cost: "$0.007", Finding: "5-7 reviewers check spec before code", HasDetail: true,
				Why: `We were catching bugs in code. Fixing them cost $0.02&ndash;$0.05 each in API calls, plus rebuild time. Then we realised: these bugs started as spec ambiguities. A wireframe that says "user list" without specifying pagination, sorting, or empty states produces code that handles none of those. Catching complexity at the spec stage costs $0.007. Catching it after code is written costs 10x more.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Spec + Wireframes"] --> B["Dev Architect"]
    A --> C["Product Owner"]
    A --> D["QA Engineer"]
    A --> E["Domain Expert"]
    B --> F{"APPROVE / REJECT?"}
    C --> F
    D --> F
    E --> F
    F -->|Approve| G["Revised Spec"]
    F -->|Reject| A
    style A fill:#f5e1e3,stroke:#d4727a
    style F fill:#fff3cd,stroke:#d4a017
    style G fill:#d4edda,stroke:#28a745
    click E "/exp/32"
    click F "/exp/25"
</pre>
Five to seven AI reviewers examine the specification before any code is written. Each has a specific persona and focus area:<br>&bull; <strong>Dev Architect:</strong> cost, complexity, deployment feasibility<br>&bull; <strong>Product Owner:</strong> CRUD completeness, missing user flows<br>&bull; <strong>QA Engineer:</strong> edge cases, error states, validation gaps<br>&bull; <strong>Market Analyst:</strong> competitive gaps, missing table-stakes features<br>&bull; <strong>Domain Expert:</strong> workflow-specific features for the product's industry`),
				How: `Each reviewer reads the wireframes and feature list independently and produces a structured verdict: approve, approve with changes, or reject with specific objections. <a href="/exp/25">Experiment 25</a> validated the gate in a full pipeline run. The breakout result came from <a href="/exp/32">Experiment 32</a>: a single domain expert prompted with invoicing knowledge found <strong>all 10 missing workflow features</strong> &mdash; print invoice, mark as paid, void, line items, tax calculation, recurring billing, payment terms, late fees, credit notes, and invoice numbering. Generic "code quality" reviewers missed every single one.`,
				Impact: `No code is written until this gate passes. At 20 reviewers ($0.10 total), there were still <strong>no diminishing returns</strong> &mdash; every additional reviewer found unique issues (<a href="/exp/39">Experiment 39</a>). The domain expert alone justified the entire step. This gate feeds a clean, reviewed spec into <a href="/exp/905">Step 5 (Exact Type Signatures)</a>, where it becomes machine-enforceable.`,
				Related: []string{"25", "32", "39"},
			},
			{Num: "Step 5", NumID: 905, Focus: "Exact Type Signatures", Result: "Exp 17", Cost: "$0.007", Finding: "Go types in spec prevent mismatches", HasDetail: true,
				Why: template.HTML("The store layer said <code>CreatedAt time.Time</code>. The HTTP handler said <code>DateCreated string</code>. The HTML template said <code>.Created</code>. Three layers, three different names for the same field, instant failure at compile time. This was the #1 cause of build failures across our experiments &mdash; not logic errors, not missing features, just <em>names that didn't match</em>. Each generation step was reading the spec independently and making its own naming choices."),
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["English Spec"] --> B["Type Extraction"]
    B --> C["Go Type Signatures"]
    C --> D["type Client struct ..."]
    D --> E["Store + HTTP Contract"]
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#d4edda,stroke:#28a745
    click B "/exp/17"
</pre>
The specification now includes actual Go code that every subsequent step must use exactly:<br><code>type Client struct {<br>&nbsp;&nbsp;ID &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;string &nbsp;&nbsp;&nbsp;json:&quot;id&quot;<br>&nbsp;&nbsp;Name &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;string &nbsp;&nbsp;&nbsp;json:&quot;name&quot;<br>&nbsp;&nbsp;Email &nbsp;&nbsp;&nbsp;&nbsp;string &nbsp;&nbsp;&nbsp;json:&quot;email&quot;<br>&nbsp;&nbsp;CreatedAt time.Time json:&quot;created_at&quot;<br>}</code><br>Input is the approved specification from <a href="/exp/904">Step 4</a>. Output is Go type declarations with exact field names, types, and JSON tags.`),
				How: `The <a href="/exp/17">V-Model pattern from Experiment 17</a> proved this approach. A "blueprint" AI writes the spec <em>and</em> hidden acceptance tests in one call &mdash; guaranteeing the types match. Then a separate "executor" AI builds from the spec without seeing the tests. If the executor follows the spec faithfully, the hidden tests pass. If it improvises, they fail with precise errors we can feed back.`,
				Impact: `<strong>100% build success rate</strong> when types are specified exactly, versus ~50% when types are described in natural language. This single change eliminated the most common failure mode across 50+ experiments. The exact types flow into <a href="/exp/906">Step 6 (Store Layer)</a> and <a href="/exp/907">Step 7 (HTTP Server)</a>, acting as a contract that keeps independently-generated layers compatible.`,
				Related: []string{"16", "17"},
			},
			{Num: "Step 6", NumID: 906, Focus: "Store Layer", Result: "Exp 27, 36", Cost: "$0.013", Finding: "Data persistence via cheap model", HasDetail: true,
				Why: `We needed to split code generation into manageable pieces. A single "build everything" prompt produces tangled code where the store logic, HTTP routing, and HTML templates are interleaved. Bugs in one layer cascade into all others. By isolating the store layer &mdash; pure Go, no HTTP, no HTML, no UI &mdash; we get a component that is easy to generate, easy to test, and easy to build on top of.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Type Signatures"] --> B["Qwen3-30B"]
    B --> C["store.go + tests"]
    C --> D{"go test"}
    D -->|Pass| E["Store Ready"]
    D -->|Fail| B
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#d4edda,stroke:#28a745
    click B "/exp/27"
    click C "/exp/16"
</pre>
The store layer implements Create, Get, List, Update, and Delete for each entity type. For MVP, in-memory maps with <code>sync.Mutex</code> protection are sufficient. Input is the type signatures from <a href="/exp/905">Step 5</a> plus CRUD requirements. Output is <code>store.go</code> with implementation and matching tests, generated together in a single call (the pattern proven in <a href="/exp/16">Experiment 16</a>).`),
				How: `Generated by the cheapest available model &mdash; Qwen3-30B via OpenRouter at $0.0005 per call. The type signatures from Step 5 ensure the store's interface matches exactly what the HTTP layer will expect. The store is built first and tested in isolation: <code>go test</code> must pass before we proceed. <a href="/exp/27">Experiment 27</a> ran this as part of a fully automated pipeline; <a href="/exp/36">Experiment 36</a> validated the constrained-simplicity approach to implementation choices.`,
				Impact: `Experiment 27 produced a <strong>720-line CRM with 41/41 tests passing, zero human intervention</strong>. The store layer was the foundation: generated first, tested in isolation, then the HTTP layer was built on top of it in <a href="/exp/907">Step 7</a>. Because the store is pure data logic with no side effects, even the cheapest models generate it reliably. This keeps the cost of the entire Build phase under $0.50.`,
				Related: []string{"16", "27", "36"},
			},
			{Num: "Step 7", NumID: 907, Focus: "HTTP Server + UI", Result: "Exp 27, 37", Cost: "$0.16", Finding: "Full app with REST API + HTML", HasDetail: true,
				Why: `The store layer exists and passes all tests. Now we need an HTTP server that calls it, HTML templates that display the data, form handlers that accept user input, and CSS that makes it look professional. This is the most expensive single step because it touches every layer of the application at once &mdash; but by this point, the types are locked (<a href="/exp/905">Step 5</a>), the store interface is known (<a href="/exp/906">Step 6</a>), and the wireframes are defined (<a href="/exp/903">Step 3</a>). The model has complete context.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Store + Wireframes"] --> B["Haiku / Sonnet"]
    B --> C["main.go + HTML/CSS"]
    C --> D{"go build"}
    D -->|Pass| E["Running App"]
    D -->|Fail| B
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#d4edda,stroke:#28a745
    click B "/exp/27"
    click E "/exp/37"
</pre>
Produces a single Go binary with REST API endpoints, embedded HTML/CSS/JS frontend, form handling, validation, and error responses. The prompt includes wireframes, type signatures, and the store interface. Generated by Haiku (free on Claude subscription) or Sonnet. Every server build prompt also includes SaaS UX patterns: breadcrumbs, toast notifications, confirmation dialogs for destructive actions, responsive tables, empty states, loading states, and error boundaries. These are &ldquo;free&rdquo; &mdash; they cost nothing extra to generate but make the product feel professional.`),
				How: `The model receives three artifacts that eliminate ambiguity: wireframes tell it what each screen looks like, type signatures tell it what data structures to use, and the store interface tells it what functions are available. <a href="/exp/27">Experiment 27</a> proved this works end-to-end in a fully automated pipeline. <a href="/exp/37">Experiment 37</a> added a 10-reviewer panel and measured the quality improvement.`,
				Impact: `Experiment 37 produced a <strong>426-line application with 26/26 tests passing</strong>. The application included SaaS UX patterns (breadcrumbs, toast notifications, confirmation dialogs) because the reviewer panel in <a href="/exp/904">Step 4</a> had requested them. This is the step where the product becomes usable &mdash; but not yet polished. That is what <a href="/exp/908">Step 8 (Post-Code Review)</a> and <a href="/exp/910">Step 10 (Progressive Enhancement)</a> handle.`,
				Related: []string{"27", "37"},
			},
			{Num: "Step 8", NumID: 908, Focus: "Post-Code Review", Result: "Exp 31", Cost: "$0.025", Finding: "UX + A11y + Security review built app", HasDetail: true,
				Why: `The application compiles and all tests pass. But "works" and "good" are different things. A form without ARIA labels works &mdash; but a screen reader cannot use it. A page that renders user content without escaping works &mdash; until someone submits a <code>&lt;script&gt;</code> tag. The pre-code reviewers in <a href="/exp/904">Step 4</a> checked the spec; now we need reviewers who check the <em>actual built artifact</em>.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Built App"] --> B["UI/UX Expert"]
    A --> C["A11y Reviewer"]
    A --> D["OWASP Security"]
    B --> E["Fix List"]
    C --> E
    D --> E
    style A fill:#f5e1e3,stroke:#d4727a
    style E fill:#d4edda,stroke:#28a745
    click B "/exp/31"
    click D "/exp/31"
</pre>
Three specialised reviewers examine the complete application code and rendered HTML:<br>&bull; <strong>UI/UX Expert:</strong> layout hierarchy, visual consistency, responsive design, form usability<br>&bull; <strong>Accessibility Reviewer:</strong> WCAG 2.1 AA compliance, ARIA labels, keyboard navigation, colour contrast<br>&bull; <strong>OWASP Security Reviewer:</strong> XSS prevention, CSRF tokens, CSP headers, input validation, SQL injection<br><br>Input is the complete application code plus rendered HTML output. Output is a fix list with severity ratings.`),
				How: `Each reviewer reads the actual source code and HTML output &mdash; not the spec, not the wireframes. They produce specific, actionable findings with severity ratings and exact code locations. <a href="/exp/31">Experiment 31</a> ran all three reviewers against a freshly-built CRM application and measured what they found versus what the pre-code reviewers had caught.`,
				Impact: `Experiment 31 added ARIA labels to all form elements, fixed two XSS vectors in user-generated content display, and added Content-Security-Policy headers. These are the kinds of issues that slip through functional testing because the app "works" without them &mdash; they only matter for quality, accessibility, and security. The fix list feeds directly into <a href="/exp/909">Step 9 (Fix Cycle)</a>.`,
				Related: []string{"31", "37"},
			},
			{Num: "Step 9", NumID: 909, Focus: "Fix Cycle + Auto-Fix", Result: "Exp 4, 32", Cost: "$0.02", Finding: "Apply fixes, auto-fix pipeline", HasDetail: true,
				Why: template.HTML("The post-code review from <a href=\"/pipeline/908\">Step 8</a> produced a fix list. Now we need to apply those fixes without breaking the application. This sounds simple, but we learned the hard way in <a href=\"/exp/32\">Experiment 32</a>: a security reviewer added <code>script-src 'self'</code> to the CSP header, which blocked ALL inline JavaScript, making the entire UI non-functional. The fix was correct in isolation but catastrophic in context. We needed an automated safety net."),
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Fix List"] --> B["Apply Fix"]
    B --> C["Auto-Fix Pipeline"]
    C --> D{"go build"}
    D -->|Pass| E{"More Fixes?"}
    D -->|Fail| B
    E -->|Yes| B
    E -->|No| F["Fixed App"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#fff3cd,stroke:#d4a017
    style F fill:#d4edda,stroke:#28a745
    click C "/exp/4"
    click B "/exp/32"
</pre>
AI applies each fix from the post-code review, then the auto-fix pipeline runs automatically after every modification:<br><br><code>goimports -w . &nbsp;&nbsp;&nbsp;&nbsp;# Fix imports</code><br><code>gofmt -w . &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;# Format code</code><br><code>fix-address-of-const.py # Fix &amp;constant errors</code><br><code>go vet ./... &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;# Static checks</code><br><code>go build . &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;# Compile gate</code><br><br>Input is the reviewer fix list plus application code. Output is an updated application with all fixes applied and all gates passing.`),
				How: template.HTML("The auto-fix pipeline was born in <a href=\"/exp/4\">Experiment 4</a>, where we discovered that <code>goimports</code> and <code>gofmt</code> alone recover a large percentage of build failures &mdash; no API calls needed. The <code>fix-address-of-const.py</code> script handles a shared blind spot across every model we tested: Go does not allow taking the address of constants, but every model generates <code>&amp;someConstant</code> at some point. After each fix is applied, the full test suite runs. If tests fail, the fix is flagged for manual review rather than silently breaking the app."),
				Impact: `The auto-fix pipeline recovers <strong>40&ndash;60% of build failures for free</strong> &mdash; zero API calls, zero cost. Combined with the AI fix application, this step turns a list of findings into a working, improved application. The critical lesson from Experiment 32: always re-run tests after applying security fixes. A "correct" fix in isolation can be catastrophic in context. The fixed application moves to <a href="/exp/910">Step 10 (Progressive Enhancement)</a>.`,
				Related: []string{"4", "32"},
			},
			{Num: "Step 10", NumID: 910, Focus: "Progressive Enhancement", Result: "Exp 38", Cost: "$0.20", Finding: "ZERO regressions across 32 tests", HasDetail: true,
				Why: `At this point we have a working, reviewed, fixed MVP. But the personas from <a href="/exp/901">Step 1</a> demanded features that did not make the MVP cut. We tried generating all 10 remaining features in one shot (<a href="/exp/32">Experiment 32</a>) &mdash; it cost $0.96 and produced test failures. Features conflicted with each other; CSS from one broke the layout of another. We needed a way to add features <em>without destroying what already worked</em>.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["MVP App"] --> B["Add Feature"]
    B --> C["Run ALL Tests"]
    C --> D{"Pass?"}
    D -->|Yes| E{"More Features?"}
    D -->|No| F["Revert + Retry"]
    E -->|Yes| B
    E -->|No| G["Enhanced App"]
    style A fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#fff3cd,stroke:#d4a017
    style F fill:#fce4ec,stroke:#d4727a
    style G fill:#d4edda,stroke:#28a745
    click B "/exp/38"
</pre>
Rather than generating all features at once, progressive enhancement adds them one at a time. Input is the base application plus a list of features from persona demands not yet implemented. After each feature addition: rebuild the application, run all existing tests, verify zero regressions, then move to the next feature. Each iteration costs approximately $0.04.`),
				How: `<a href="/exp/38">Experiment 38</a> ran 5 iterations of progressive enhancement on a CRM application, adding one feature per iteration while running the full 32-test suite after each addition. We compared this directly against the one-shot approach from Experiment 32 on cost, test stability, and code quality.`,
				Impact: `Progressive enhancement achieved <strong>ZERO regressions across 5 iterations and 32 tests</strong>. The one-shot approach cost $0.96 with test failures; progressive enhancement cost $0.20 with perfect stability. This is the most expensive step after the initial server build, but it is where the product goes from "MVP that works" to "product with all the features personas demanded." The enhanced application then enters the Verify phase, starting with <a href="/exp/911">Step 11 (Playwright Testing)</a>.`,
				Related: []string{"32", "38"},
			},
			{Num: "Step 11", NumID: 911, Focus: "Playwright Testing", Result: "Exp 35, 40", Cost: "FREE", Finding: "4 bugs missed by everything else", HasDetail: true,
				Why: template.HTML("We had 26 passing unit tests and 10 approved reviewers. The application compiled, the store worked, the HTTP handlers returned correct responses. Then we opened it in a browser. A modal overlay with <code>display:flex</code> overrode its <code>display:none</code> and blocked every click on the page. A form using <code>new FormData()</code> sent <code>multipart/form-data</code> but the handler called <code>r.ParseForm()</code>, which only parses URL-encoded data. None of these bugs were visible to unit tests or code reviewers &mdash; they only manifest in a real browser."),
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Running App"] --> B["Browser Agent"]
    B --> C["Click / Type / Navigate"]
    C --> D{"Bugs Found?"}
    D -->|Yes| E["Bug Report"]
    D -->|No| F["All Journeys Pass"]
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#fce4ec,stroke:#d4727a
    style F fill:#d4edda,stroke:#28a745
    click B "/exp/35"
    click C "/exp/40"
</pre>
AI personas <em>use</em> the application in a real browser via Playwright. This is not unit testing &mdash; it is integration testing at the UI level, where a simulated user clicks buttons, fills forms, and navigates between pages. Input is a running application. Output is test results, screenshots, and bug reports. The cost is free because Playwright runs locally.`),
				How: template.HTML("<a href=\"/exp/35\">Experiment 35</a> developed the persona-driven journey approach: each AI persona from Step 1 gets a Playwright script that walks through their primary user journey. <a href=\"/exp/40\">Experiment 40</a> combined Playwright with adversarial testing. The bugs Playwright found across these experiments:<br>&bull; <strong>CSP blocking JavaScript</strong> (Exp 32): security reviewer added <code>script-src 'self'</code> which blocked all inline scripts<br>&bull; <strong>Multipart form parsing</strong> (Exp 32): <code>r.ParseForm()</code> does not parse <code>multipart/form-data</code><br>&bull; <strong>JSON error responses</strong> (Exp 32): <code>http.Error()</code> returns plain text but the frontend expects JSON<br>&bull; <strong>Modal overlay blocking clicks</strong> (Exp 37): CSS <code>display:flex</code> overrode <code>display:none</code>"),
				Impact: `<strong>4 bugs found that 26 passing tests + 10 approved reviewers missed.</strong> Every one of these bugs would have shipped to production without Playwright. The testing is free and non-negotiable &mdash; it catches an entire class of defects invisible to every other technique in the pipeline. Any bugs found here feed back through <a href="/exp/909">Step 9 (Fix Cycle)</a> before proceeding to <a href="/exp/912">Step 12 (Security Testing)</a>.`,
				Related: []string{"32", "35", "37", "40"},
			},
			{Num: "Step 12", NumID: 912, Focus: "Security Testing", Result: "Exp 41-45", Cost: "$0.03", Finding: "Pen test + chaos + static analysis", HasDetail: true,
				Why: `The OWASP reviewer in <a href="/exp/908">Step 8</a> checked for common vulnerabilities in the source code. But reading code and <em>attacking</em> a running application are fundamentally different skills. A code reviewer sees that an endpoint exists; a pen tester discovers that endpoint has no authentication. We needed to actually try to break the application, not just read it.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Running App"] --> B["Static Analysis"]
    A --> C["Adversarial Agent"]
    A --> D["Pen Test Agent"]
    B --> E["Vulnerability Report"]
    C --> E
    D --> E
    style A fill:#f5e1e3,stroke:#d4727a
    style E fill:#d4edda,stroke:#28a745
    click B "/exp/43"
    click C "/exp/40"
    click D "/exp/42"
</pre>
Three layers of security testing, each catching different classes of vulnerabilities:<br>&bull; <strong>Static analysis</strong> (FREE): <code>gosec</code> and <code>govulncheck</code> scan code for known vulnerability patterns and dependency CVEs<br>&bull; <strong>Adversarial testing</strong> ($0.01): AI agent sends malformed inputs, boundary values, state manipulation attempts, and unexpected content types<br>&bull; <strong>Pen test agent</strong> ($0.02): Attempts auth bypass, IDOR (accessing other users&#39; data by changing IDs), SQL injection, XSS, rate-limit abuse, and privilege escalation<br><br>Input is a running application. Output is a vulnerability report with severity ratings.`),
				How: `<a href="/exp/40">Experiments 40&ndash;42</a> developed the AI testing and pen test agents. <a href="/exp/43">Experiments 43&ndash;45</a> added static analysis tools and a chaos agent that subjects the server to extreme conditions: 50 concurrent connections, 1MB request bodies, malformed HTTP headers, and slow client attacks. GDPR review (<a href="/exp/50">Experiment 50</a>) and multi-tenant isolation review (<a href="/exp/52">Experiment 52</a>) were later additions that check compliance-specific concerns.`,
				Impact: `Experiment 42 found <strong>4 High severity vulnerabilities</strong> including missing authentication on admin endpoints, IDOR on client data, and no rate limiting. The chaos agent (Experiment 45) threw 25 attack scenarios at a Go HTTP server &mdash; it survived all 25, proving Go's built-in HTTP server is remarkably resilient. GDPR review found 10 non-compliant items in a standard CRUD app; multi-tenant review found 19 of 20 expected isolation features missing. Any critical or high findings block <a href="/exp/913">Step 13 (Deploy)</a>.`,
				Related: []string{"40", "42", "43", "45", "50", "52"},
			},
			{Num: "Step 13", NumID: 913, Focus: "Deploy", Result: "Exp 55", Cost: "—", Finding: "Docker → VPS → Caddy + TLS", HasDetail: true,
				Why: `An application that passes all tests but never reaches users is a waste of pipeline time. The final step must take a green build and put it on a public URL with HTTPS, automatically. We run 19 deployed products on this infrastructure, so the deploy pattern is proven and repeatable &mdash; but it took <a href="/exp/55">Experiment 55</a> to codify it into a step that the pipeline can execute without human intervention.`,
				What: template.HTML(`<pre class="mermaid">
flowchart TD
    A["Passing App"] --> B["Dockerfile"]
    B --> C["docker build"]
    C --> D["docker run on VPS"]
    D --> E["Caddy + TLS"]
    E --> F["Live HTTPS Service"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#f5e1e3,stroke:#d4727a
    style E fill:#f5e1e3,stroke:#d4727a
    style F fill:#d4edda,stroke:#28a745
    click A "/exp/55"
</pre>
Deployment follows a 5-step pattern:<br>1. <strong>Dockerfile:</strong> Multi-stage Go build (builder + scratch/alpine) producing a minimal container<br>2. <strong>Docker build:</strong> <code>docker build -t product .</code> on the VPS<br>3. <strong>Docker run:</strong> Container started with environment variables for configuration<br>4. <strong>Caddy reverse proxy:</strong> Automatic TLS certificate provisioning via Let&#39;s Encrypt<br>5. <strong>Domain:</strong> DNS pointed to VPS, Caddy handles HTTPS termination<br><br>Input is a passing application with all tests green. Output is a running production service with HTTPS.`),
				How: template.HTML("Go's single-binary deployment makes containerization trivial &mdash; the Dockerfile is typically 10 lines. The multi-stage build compiles in a full Go image, then copies only the binary into a scratch or alpine image. No runtime dependencies, no package managers, no version conflicts. Caddy handles TLS automatically via Let's Encrypt &mdash; no certificate management required.<br><br><strong>Critical lesson from production:</strong> <code>docker restart</code> does NOT reload environment variables. If you change <code>.env</code>, you must <code>docker stop</code>, <code>docker rm</code>, then <code>docker run</code> again. This cost us hours of debugging before we codified it as a rule."),
				Impact: `The entire infrastructure is proven across <strong>19 deployed products</strong> in the Dark Factory system. From an 8-word brief to a running HTTPS service, the pipeline costs $0.49 with cheap models and $2.75 with premium (Opus). The deploy step is the payoff: everything before it &mdash; personas, wireframes, reviewers, tests, fixes &mdash; exists to ensure that what we deploy actually works and is worth deploying.`,
				Related: []string{"55"},
			},
		},
	},
	{
		Slug:     "code-generation",
		Name:     "Code Generation",
		Icon:     "\U0001f4bb",
		ExpRange: "Experiments 1 – 21",
		Narrative: []template.HTML{
			`The foundational question: can an AI model take a specification and produce compilable, tested code? Across roughly 50 experiments &mdash; spanning code generation, product design, review panels, testing, and security &mdash; we explored escalation strategies, sub-task granularity, auto-fix pipelines, model routing, progressive enhancement, and multi-service architectures. The answer is unambiguously yes &mdash; but the path there is littered with subtlety.`,
			`The single most important finding is that <strong>prompt wording matters more than model choice</strong>. The same model can go from 0% to 100% success rate with a better prompt. Experiment 3 demonstrated this with V4-style prompts ("just output the file, no explanation") outperforming verbose instructions on every model. Experiment 15 confirmed it: a tiered escalation system was unnecessary because the cheapest model passed every task once the prompt was fixed.`,
			`The <strong>V-Model pattern</strong> (Experiment 10, 17) proved that embedding exact type signatures in the specification prevents the most common class of failure: type mismatches between generated modules. When the spec says <code>func NewStore(path string) *Store</code>, the model produces that exact signature. When it says "create a store", every model invents a different API. Precision in the spec is free and eliminates an entire failure category.`,
		},
		KeyInsight: `Parser v2 (Experiment 13) handles 8 file output formats reliably. The auto-fix pipeline &mdash; goimports, gofmt, fix-address-of-const &mdash; recovers 40&ndash;60% of build failures without any API calls.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "01", NumID: 1, Focus: "Escalation (Cheap to Strong)", Result: "0/5 FAIL", Cost: "$0.13", Finding: "Escalation does not fix shared blind spots", HasDetail: true, Icon: "\U0001f504",
				Why:     `If a cheap model fails to produce compilable code, can we fix it by escalating to a more expensive model? This is the intuitive approach — try cheap first, fall back to premium.`,
				What:    `Three-tier escalation chain: <strong>Qwen3-30B</strong> ($0.0005/call) → <strong>MiniMax M2.7</strong> ($0.007/call) → <strong>Sonnet</strong> ($0.012/call). Five Go code generation tasks from the task-board application. When a tier fails <code>go build</code>, the error message is passed to the next tier with "fix this error."`,
				How:     `Each tier receives: (1) the original prompt, (2) the previous tier's output, (3) the compiler error. The stronger model attempts to fix the specific error. If all three tiers fail, the task is marked as failed. Auto-fix pipeline (goimports/gofmt) runs between each tier.`,
				Impact:  `<strong>Killed the escalation strategy.</strong> Proved that model capability isn't the bottleneck — shared blind spots are. This redirected research toward better prompts (Exp 3) and deterministic auto-fix scripts (Exp 4) instead of more expensive models.`,
				Related: []string{"3", "4", "5", "14", "15"},
			},
			{Num: "02", NumID: 2, Focus: "Sub-Task Granularity", Result: "Partial", Cost: "$0.04", Finding: "1 file per task is optimal for quality", HasDetail: true,
				Why:     `How many files should each AI call generate? One file per call is slow but focused. All files in one call is fast but the model loses coherence across file boundaries.`,
				What:    `Three granularity levels tested on the task-board app: <strong>(a)</strong> 1 file per call (5 calls), <strong>(b)</strong> 2 files per call (3 calls), <strong>(c)</strong> all files in one call (1 call). Model: Qwen3-30B. Gate: <code>go vet</code> + <code>go build</code> per file.`,
				How:     `Each configuration runs 3 times. The architecture spec is provided to every call. For multi-file calls, the model outputs file blocks (<code>--- FILE: path ---</code>). Parser v1 extracts files. Auto-fix pipeline runs after extraction.`,
				Impact:  `Established <strong>1 file per sub-task</strong> as the default. Later refined in Exp 16 to 2 files when the files are tightly coupled (e.g. implementation + tests).`,
				Related: []string{"16"},
			},
			{Num: "03", NumID: 3, Focus: "V1 Re-Run (Improved Prompts)", Result: "100% pass", Cost: "$0.05", Finding: `V4 prompt: "just output the file" wins`, HasDetail: true,
				Why:     `Experiments 1 and 2 both failed. Before blaming the models, test whether the <em>prompt</em> is the problem. Can the same cheap model succeed with better instructions?`,
				What:    `Re-ran all failing tasks from Exp 1-2 on the same model (<strong>Qwen3-30B</strong>, $0.0005/call) with four prompt variations:<br><strong>V1:</strong> "Generate a Go file that implements..."<br><strong>V2:</strong> "Write the following Go file..."<br><strong>V3:</strong> "Create this file. Output Go code only."<br><strong>V4:</strong> "Output only the complete file contents for {path}. No explanation. No markdown fences. Start with package {pkg}."`,
				How:     `Each prompt variant runs on the same 5 sub-tasks, 3 times each. Gate: <code>go build</code> passes. The V4 prompt was designed to eliminate parser failures by forbidding markdown wrapping and explanatory text.`,
				Impact:  `<strong>The single most important finding.</strong> V4 prompts became the standard for all subsequent experiments. This proved that prompt engineering is more valuable than model upgrades — the cheapest model at $0.0005/call matches $0.30/call models when prompted correctly. All 11 models tested passed with V4 prompts.`,
				Related: []string{"1", "8", "15"},
			},
			{Num: "04", NumID: 4, Focus: "Auto-Fix Pipeline", Result: "40–60% fixed", Cost: "FREE", Finding: "goimports, gofmt, sed, vet, build", HasDetail: true, Icon: "\U0001f527", SourceFile: "auto-fix-go.sh",
				Why:     `Exp 1 showed models share blind spots that escalation can't fix. Can deterministic scripts fix these patterns without any AI calls?`,
				What:    `A sequential pipeline of Go tooling that runs after every AI generation: <strong>goimports</strong> (fix imports) → <strong>gofmt</strong> (format) → <strong>fix-address-of-const.py</strong> (fix &amp;constant) → <strong>go vet</strong> (static analysis) → <strong>go build</strong> (compile gate). Zero API cost.`,
				How:     `Ran the pipeline on all 15 failed outputs from Exp 1-2. Counted how many compile errors were fixed without re-calling the model. The <code>fix-address-of-const.py</code> script detects <code>&amp;someConstant</code> patterns and generates helper variables: <code>v := someConstant; &amp;v</code>.`,
				Impact:  `<strong>Became a mandatory step after every code generation.</strong> The auto-fix pipeline is Step 9 of the pipeline. It recovers 40–60% of failures for free, making retry loops shorter and cheaper. The <code>fix-address-of-const.py</code> script handles the single most common cross-model failure.`,
				Related: []string{"1", "9", "16"},
			},
			{Num: "05", NumID: 5, Focus: "Model Routing by Task Type", Result: "0/5 FAIL", Cost: "$0.05", Finding: "Single-shot fails; retry loop essential", HasDetail: true,
				Why:     `If different models excel at different tasks (planning vs coding vs testing), can we route each sub-task to the optimal model?`,
				What:    `Routing matrix: <strong>Gemini Flash</strong> for planning tasks, <strong>MiniMax M2.7</strong> for Go code, <strong>Qwen3-30B</strong> for tests. Each task gets one attempt on its designated model. Gate: <code>go build</code>.`,
				How:     `Five tasks routed to their "best" model based on task type. No retries — single-shot only. If the designated model fails, the task fails. This tests whether smart routing alone is sufficient.`,
				Impact:  `Proved that <strong>retry + auto-fix matters more than model routing</strong>. A cheap model with 3 retries beats an expensive model with 1 attempt. Model routing became useful later (Exp 7) but only after the retry infrastructure was in place.`,
				Related: []string{"1", "7", "14"},
			},
			{Num: "06", NumID: 6, Focus: "Claude Sub Models", Result: "Both pass", Cost: "FREE", Finding: "Haiku 3x faster, equal quality", HasDetail: true,
				Why:     `Claude subscription gives free access to both Haiku and Sonnet via <code>claude -p</code>. Is there any quality difference for code generation?`,
				What:    `Identical V4 prompts sent to <strong>Haiku</strong> and <strong>Sonnet</strong> via <code>claude -p --model haiku</code> and <code>claude -p --model sonnet</code>. Same 5 task-board sub-tasks. Gate: <code>go build</code> + <code>go test</code>.`,
				How:     `Each model runs all 5 tasks 3 times. Measured: pass rate, generation time, output quality (lines, structure). Claude -p uses its tool system (Read/Write/Edit) rather than raw output, so the model has full workspace access.`,
				Impact:  `<strong>Haiku became the default for all code generation.</strong> 3x faster, equal quality, free on subscription. Sonnet/Opus reserved for judgment tasks (reviews, design). This decision saves ~$0.15 per pipeline run.`,
				Related: []string{"7", "12"},
			},
			{Num: "07", NumID: 7, Focus: "Hybrid Pipeline", Result: "Design", Cost: "—", Finding: "Cheap API for planning, free sub for execution", HasDetail: true,
				Why:     `We now know cheap models work for code (Exp 3) and Haiku is free for generation (Exp 6). What's the optimal cost architecture?`,
				What:    `Architecture design: <strong>OpenRouter API</strong> (Qwen3-30B at $0.0005/call) for planning, reviewing, lightweight tasks. <strong>Claude subscription</strong> (<code>claude -p</code>, free) for code generation. No expensive API calls needed.`,
				How:     `Cost modeling across the 13-step pipeline. Each step categorized as "planning" (cheap API) or "generation" (free sub). Total: ~$0.05 for planning/review + $0.00 for generation = <strong>$0.05 per product</strong>.`,
				Impact:  `Became the <strong>standard cost architecture</strong> for the pipeline. Every subsequent experiment uses this hybrid approach. Total pipeline cost is dominated by reviewer count, not model quality.`,
				Related: []string{"3", "6", "27"},
			},
			{Num: "08", NumID: 8, Focus: "MiniMax Backtick Hint", Result: "2/3 pass", Cost: "$0.03", Finding: "Explicit hint: use concatenation, not literals", HasDetail: true,
				Why:     `MiniMax M2.7 kept failing on Go code that needed backtick raw strings inside backtick-delimited output. Can a model-specific prompt hint fix a known failure pattern?`,
				What:    `Added explicit instruction to MiniMax prompts: <em>"Do not use backtick characters in string literals. Use string concatenation or fmt.Sprintf instead."</em> Tested on 3 sub-tasks that previously failed due to nested backticks.`,
				How:     `A/B test: same 3 tasks with and without the hint. Model: <strong>MiniMax M2.7</strong> ($0.007/call). Gate: <code>go build</code>. Without hint: 0/3 pass. With hint: 2/3 pass.`,
				Impact:  `Created the concept of a <strong>model-specific prompt hint library</strong>. Each model gets additional instructions for its known blind spots. MiniMax gets the backtick hint, all models get "no markdown fences" from V4 prompts.`,
				Related: []string{"3", "15"},
			},
			{Num: "09", NumID: 9, Focus: "Full App via API Only", Result: "Parser fail", Cost: "$0.045", Finding: "Parser is the weakest link (~40% of failures)", HasDetail: true, SourceFile: "parse-blocks-v2.py",
				Why:     `Can we generate an entire multi-file application in one API call and parse the output into individual files?`,
				What:    `Single OpenRouter API call to <strong>Qwen3-30B</strong> requesting all 5 task-board files. Output parsed by the file block parser (v1) to extract individual files. Gate: all files exist + <code>go build</code>.`,
				How:     `Model outputs files using <code>--- FILE: path ---</code> markers. Parser v1 splits on markers and writes files. But models also output markdown fences, explanatory text, and varied marker formats — parser v1 only handled one format.`,
				Impact:  `Revealed that <strong>40% of "generation failures" were parser failures</strong>, not model failures. The AI wrote correct code — we just couldn't extract it. Led directly to Parser v2 (Exp 13) which handles 8 output formats.`,
				Related: []string{"13"},
			},
			{Num: "10", NumID: 10, Focus: "V-Model Pattern", Result: "Conceptual", Cost: "$0.032", Finding: "Hidden acceptance tests work as surprise gate", HasDetail: true,
				Why:     `The V-Model from software engineering pairs each requirement with a test. Can we apply this to AI generation — write hidden acceptance tests that verify the AI's output matches the spec?`,
				What:    `Spec includes exact function signatures: <code>func NewStore(path string) *Store</code>. Hidden acceptance tests call those exact functions. The AI never sees the tests — it implements the spec, and if the implementation matches, the tests pass automatically.`,
				How:     `Write spec with exact Go types → AI generates implementation → hidden golden tests run against the output. If the AI followed the spec faithfully, tests pass. If it invented its own API, tests fail with clear "expected X got Y" messages.`,
				Impact:  `<strong>Became Step 5 of the pipeline (Exact Type Signatures).</strong> Refined in Exp 17 (V-Model Full Loop) where 100% pass rate was achieved. Eliminates the entire class of type-mismatch failures between generated modules.`,
				Related: []string{"17", "5"},
			},
			{Num: "11", NumID: 11, Focus: "PR Review Gate", Result: "Design", Cost: "—", Finding: "AI reviewer catches issues tests don't", HasDetail: true,
				Why:     `Unit tests verify that code does what the programmer intended. But what if the programmer's intention was wrong? Tests pass, the build is green, and the code ships — with a missing feature, a UX dead end, or a security hole that no test was written for. We had seen this repeatedly: 26 passing tests that missed 4 browser-level bugs (later caught by Playwright in Exp 40). The gap was clear — we needed a quality gate that thinks about the <em>product</em>, not just the code.`,
				What:    `A design experiment for an AI-powered PR review gate. The idea: before code merges, one or more AI reviewers examine the diff with specific personas — a security reviewer, a UX expert, an accessibility auditor, a domain specialist. Each reviewer has a focused lens and outputs structured findings with severity ratings. The gate passes only if no Critical or High findings remain unaddressed.`,
				How:     `This was a design-only experiment — no script, no execution. We sketched the architecture: a PR is submitted, the diff is extracted, and routed to 3-5 AI reviewers in parallel. Each reviewer gets the diff plus the project spec and returns APPROVE, APPROVE WITH FIXES, or REJECT. Findings are categorised by severity (Critical, High, Medium, Low). The merge is blocked until all Critical/High findings are resolved. The design drew on our experience with pre-code review panels (Exp 30, 32) — if reviewers work before code is written, they should work even better when reviewing actual code.`,
				Impact:  `This design became the foundation for the <strong>post-code review step</strong> (Step 8 in the pipeline) and was implemented concretely in <a href="/exp/31">Experiment 31</a> with three specialised reviewers (UX, accessibility, OWASP security). It also directly influenced the Dark Factory production PR review system, which runs 10 AI reviewers with 4 different reviewer PATs. The insight that a reviewer's <em>persona</em> determines what it finds — not the model's raw capability — shaped every subsequent review experiment.`,
				Related: []string{"30", "31", "32", "39"},
			},
			{Num: "12", NumID: 12, Focus: "Full-Stack App (70s)", Result: "7 files built", Cost: "FREE", Finding: "Complete app from schema in ~1 minute", HasDetail: true,
				Why:     `Every experiment so far had used the OpenRouter API to call models. But we had a Claude subscription sitting there with unlimited Haiku and Sonnet access via <code>claude -p</code>. The API approach cost money and fought with parsers. Could we skip all that and just tell Claude to build the entire app in one shot, using its built-in file tools to write directly to disk?`,
				What: template.HTML("A single <code>claude -p --model haiku</code> call with the task-board architecture spec as input. No parser needed — Claude's tool system (Read, Write, Edit) writes files directly to the workspace. No API cost — Haiku is free on the subscription. The prompt: <em>\"Read the architecture spec. Build the complete task-board application. Write all files. Run go build to verify.\"</em>"),
				How: template.HTML("One command: <code>claude -p --model haiku --dangerously-skip-permissions</code> with the full architecture document. Claude read the spec, created 7 files (schema.graphql, model/task.go, model/task_test.go, main.go, go.mod, and supporting files), ran <code>go build</code> to verify compilation, and fixed issues it found — all within a single session of about 70 seconds. No parser, no file extraction, no retry loop. The model handled the entire workflow autonomously using its tool system."),
				Impact: template.HTML("This was the moment <code>claude -p</code> became the default for code generation in the pipeline. The API-based approach required parsers, cost money, and fought with output formatting. <code>claude -p</code> writes files directly, self-verifies with <code>go build</code>, and costs nothing on the subscription. Every subsequent experiment that needed code generation (Exp 20, 21, 24, 25) used <code>claude -p</code> for the HTTP server layer while keeping the cheap API models for lightweight tasks like persona generation and reviews. One caveat discovered: <code>claude -p</code> needs a foreground terminal — it cannot run as a background process, which matters for pipeline automation."),
				Related: []string{"6", "20", "21"},
			},
			{Num: "13", NumID: 13, Focus: "Parser Hardening", Result: "v2: 8/8 tests", Cost: "$0.07", Finding: "Parser was not the bottleneck; code quality is", HasDetail: true, SourceFile: "exp13-parser-hardening.py",
				Why:     `Experiment 9 revealed that 40% of "generation failures" were actually parser failures — the AI wrote correct code, but our extraction logic mangled it. We had a v1 parser that only understood one output format: <code>--- FILE: path ---</code> blocks. But models also output markdown fences, bare Go code, varied marker formats, and explanatory text mixed with code. If we could fix the parser, maybe our success rate would jump.`,
				What: template.HTML("A head-to-head test of the old parser (v1) against a new hardened parser (v2) that handles 8 different output formats. Both parsers run on identical model output from 5 runs of the full task-board build (schema + model + test + main) on <strong>Qwen3-30B</strong>. The v2 parser recognizes: <code>--- FILE ---</code> blocks, markdown fences with file hints, bare <code>package</code> declarations, inline Go code blocks, mixed formats, and several edge cases the v1 parser choked on."),
				How: template.HTML("Five runs on Qwen3-30B, each requesting all 4 task-board files in a single API call. For each run, both parsers attempt to extract files from the identical raw output. We measured three gates: (1) file extraction — did the parser find all 4 files? (2) <code>go build</code> — does the extracted code compile? (3) <code>go test</code> — do the tests pass? The auto-fix pipeline (goimports, gofmt, fix-address-of-const) runs after extraction for both parsers."),
				Impact: template.HTML("The surprise: <strong>v1 and v2 had identical results</strong> — 100% extraction rate, 60% build rate, 0% test rate. The parser was not the bottleneck after all. Both parsers extracted all 4 files in every run, but the code itself had type mismatches and logic errors that no parser can fix. This was a crucial redirection: we stopped investing in parser improvements and pivoted to fixing <em>code quality</em> — better prompts (Exp 15), 2-file generation (Exp 16), and the V-Model pattern (Exp 17). Parser v2 became the standard anyway (it handles more edge cases), but the real lesson was: <strong>don't optimise the wrong bottleneck</strong>."),
				Related: []string{"9", "15", "16"},
			},
			{Num: "14", NumID: 14, Focus: "Model Routing + Retry", Result: "50%", Cost: "$0.111", Finding: "Retry works for fixable errors, not blind spots", HasDetail: true, SourceFile: "exp14-routing-with-retry.py",
				Why: template.HTML("Experiment 5 tested model routing — sending each sub-task to the \"best\" model for that task type — and got 0/5. But Exp 5 was single-shot: one attempt per model, no retries. That felt unfair. The auto-fix pipeline (Exp 4) can recover 40-60% of failures. What if we combined smart routing with a retry loop? Maybe the 0% was because we gave up too quickly, not because routing was wrong."),
				What: template.HTML("Re-ran the model routing experiment with <strong>3 retry attempts</strong> per sub-task plus the full auto-fix pipeline between each attempt. Two sub-tasks: (1) model + test files routed to <strong>Qwen3-30B</strong> (cheapest), (2) main.go routed to the same model with error context from previous failures. Each retry includes the compiler error from the failed attempt, so the model can try to fix the specific issue. Three complete runs."),
				How: template.HTML("For each sub-task, the pipeline runs: generate code &rarr; parse output &rarr; auto-fix (goimports, gofmt, fix-address-of-const) &rarr; <code>go build</code>. If build fails, the error is appended to the prompt and the model tries again, up to 3 times. The model/test sub-task uses the 2-file pattern — implementation and tests generated together. The main.go sub-task includes the backtick hint from Exp 8."),
				Impact: template.HTML("<strong>Model + test: 100% (3/3 runs, first attempt every time).</strong> Main.go: <strong>0% (0/3 runs, all 3 retries exhausted).</strong> Overall: 50% — the store layer works perfectly, the HTTP server never compiles. The retries didn't help main.go because the failure wasn't a fixable compiler error — it was a fundamental issue with generating Go raw strings containing backticks for embedded JavaScript. The model made the same structural mistake every time, and feeding back the error just produced a different variation of the same problem. This confirmed: <strong>retry loops fix typos and import errors, not architectural blind spots</strong>. The main.go problem needed a better prompt (solved in Exp 15), not more attempts."),
				Related: []string{"5", "4", "15", "18"},
			},
			{Num: "15", NumID: 15, Focus: "Tiered Escalation", Result: "100% T1 only", Cost: "$0.030", Finding: "Better prompt fixed backtick — cheapest model does it all", HasDetail: true, SourceFile: "exp15-tiered-escalation.py",
				Why: template.HTML("Experiments 1, 5, and 14 all hit the same wall: main.go fails because models put backtick template literals inside Go raw strings (which are themselves delimited by backticks). Retry loops didn't fix it. Escalating to stronger models didn't fix it. But Exp 8 showed that a model-specific hint (\"don't use backticks in string literals\") improved MiniMax from 0/3 to 2/3. What if we combined the hint approach with tiered escalation — try cheap first with the hint, escalate only if it still fails?"),
				What: template.HTML("Three-tier escalation: <strong>T1: Qwen3-30B</strong> ($0.0005/call) with improved prompt hints &rarr; <strong>T2: MiniMax M2.7</strong> ($0.003/call) with backtick-specific hint &rarr; <strong>T3: claude -p Haiku</strong> (free, uses tool system). For each sub-task, try T1 up to 2 times with auto-fix. If T1 fails, try T2 once with the error context. If T2 fails, escalate to T3. The improved prompt includes CORRECT/WRONG examples: <em>\"CORRECT: 'Hello ' + name. WRONG: using backtick template literals inside Go raw strings.\"</em>"),
				How: template.HTML("Three complete runs across two sub-tasks (model+test, main.go). The key change from Exp 14: the prompt now includes explicit examples of correct and incorrect JavaScript-in-Go patterns. Instead of saying \"don't use backticks,\" we showed: <code>element.textContent = value</code> (CORRECT) vs template literal syntax (WRONG). We also added: <code>'&lt;div&gt;' + title + '&lt;/div&gt;'</code> for building HTML strings in JavaScript. Each tier has its own prompt variant optimised for that model's known weaknesses."),
				Impact: template.HTML("<strong>100% success rate — all 3 runs passed on T1 (cheapest tier), first attempt.</strong> The escalation to T2 and T3 was never needed. Total cost: $0.030 across all 3 runs ($0.01 each). This was the definitive proof that <strong>prompt engineering beats model upgrades</strong>. The same model that failed in Exp 14 (0% on main.go) now passes every time with a better prompt. The backtick problem — which had blocked progress across 4 experiments — was solved by showing the model what to do instead of telling it what not to do. The CORRECT/WRONG hint pattern became standard for all subsequent code generation prompts."),
				Related: []string{"1", "8", "14", "18"},
			},
			{Num: "16", NumID: 16, Focus: "Sub-Task Granularity v2", Result: "100% (5/5)", Cost: "$0.115", Finding: "2 files/task + auto-fix = 100% on cheapest model", HasDetail: true, Icon: "\U0001f4d0", SourceFile: "exp16-granularity.py",
				Why:     `Our spec said "create a store with CRUD operations." The AI built one — with a constructor called <code>NewStore()</code>. Then we asked the AI to write tests. It wrote tests that called <code>CreateStore()</code>. Same idea, different name, instant failure.<br><br>This wasn't a one-off. The implementation used <code>task.CreatedAt</code>, the tests checked <code>task.DateCreated</code>. The store returned <code>*Task</code>, the tests expected <code>Task</code>. Every combination was slightly wrong. The maddening part: both files compiled perfectly. The tests only failed at runtime, with errors like <em>"undefined: CreateStore"</em> — a function that never existed but the test was sure should be there.`,
				What: template.HTML(`The root cause was almost too simple. Each AI call starts fresh — it has no memory of what the previous call produced. When you ask "implement a store" in one call and "write tests for a store" in another, both calls read the spec independently and make their own naming choices. It's like asking two people to build matching puzzle pieces in separate rooms.

<pre class="mermaid">
flowchart TD
    S["📄 Spec:\ncreate a store\nwith CRUD"] --> A["Call 1:\nimplements\nNewStore()"]
    S --> B["Call 2:\ntests for\nCreateStore()"]
    A --> C["task.CreatedAt\nreturns *Task"]
    B --> D["task.DateCreated\nexpects Task"]
    C --> E["❌ Names don't match\ntests fail at runtime"]
    D --> E
    style E fill:#fce4ec,stroke:#d4727a
</pre>

Our proposed fix: generate both files in a single call. If the model writes <code>NewStore()</code> in the implementation, it can <em>see</em> that name while writing the tests and will naturally use the same one.

<pre class="mermaid">
flowchart TD
    S2["📄 Spec:\ncreate a store\nwith CRUD"] --> AB["Single Call:\nimplements NewStore()\nAND tests NewStore()"]
    AB --> C2["task.CreatedAt\nreturns *Task"]
    AB --> D2["task.CreatedAt\nexpects *Task"]
    C2 --> E2["✅ Names match\ntests pass"]
    D2 --> E2
    style E2 fill:#d4edda,stroke:#28a745
</pre>`),
				How:     `We ran three configurations head-to-head, each three times on the cheapest model available (<strong>Qwen3-30B</strong>, half a thousandth of a cent per call):<br><br><strong>Config A</strong> generated <code>task.go</code> alone — the baseline. It always compiled, but there were no tests to validate the logic.<br><br><strong>Config B</strong> generated <code>task_test.go</code> alone, with the real <code>task.go</code> already sitting right there in the workspace. You'd think the model would read it. It didn't. It imagined a different API and wrote tests against that instead. This is the configuration that kept failing across every previous experiment.<br><br><strong>Config C</strong> generated both files in a single call. The prompt said: <em>"Create BOTH files. The tests must pass against YOUR implementation."</em> One context window, one set of naming decisions.<br><br>The gate was deliberately strict: not <code>go build</code> but <code>go test</code>. The tests had to actually <em>pass</em>, not just compile.`,
				Impact:  `Config C hit <strong>100% — fifteen out of fifteen runs</strong>. Every test passing. On the cheapest model available. The fix was embarrassingly obvious in hindsight: let the model see both sides of the contract at the same time.<br><br>This became a core rule in the pipeline: <strong>tightly coupled files are always generated together</strong> — implementation and its tests, schema and its resolvers, types and their validation. Loosely coupled files — like the store layer and the HTTP server — are generated in separate calls, but connected by exact type signatures written into the spec. That connection is what <a href="/exp/17">Experiment 17</a> went on to prove.`,
				Related: []string{"2", "17", "27"},
			},
			{Num: "17", NumID: 17, Focus: "V-Model Full Loop", Result: "100% (3/3)", Cost: "$0.037", Finding: "Hidden acceptance tests verify the AI followed the spec", HasDetail: true, SourceFile: "exp17-vmodel.py",
				Why:     `<a href="/exp/16">Experiment 16</a> solved the mismatch problem for tightly coupled files — generate them together. But what about loosely coupled files? The store layer and the HTTP server are built in separate calls, potentially by different models. How do we ensure the HTTP server calls <code>store.Create(task)</code> and not <code>store.AddTask(t)</code>?<br><br>We borrowed an idea from software engineering's V-Model: for every requirement, there's a matching test. What if a "blueprint" AI writes the spec <em>and</em> hidden acceptance tests at the same time — and then a separate "executor" AI builds from the spec without ever seeing those tests? If the executor follows the spec faithfully, the hidden tests pass. If it improvises, they fail with precise errors we can feed back.`,
				What: template.HTML(`The trick is that the blueprint creates both sides of a contract in one call — so they're guaranteed to match — but the executor only sees one side.

<pre class="mermaid">
flowchart TD
    A["Blueprint AI\nsingle call"] --> B["📄 Spec\nfunc NewStore&#40;&#41; *Store\nfunc &#40;s *Store&#41; Create&#40;t Task&#41; Task"]
    A --> C["🧪 Hidden Tests\nstore := NewStore&#40;&#41;\ntask := store.Create&#40;...&#41;"]
    B --> |"executor sees this"| D["  "]
    C --> |"executor never\nsees this"| E["  "]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style D fill:#fff,stroke:#fff
    style E fill:#fff,stroke:#fff
</pre>

<strong>Phase 1 — Blueprint:</strong> One model call generates a spec with exact Go function signatures AND a matching test file. Both derived from the same prompt, so <code>NewStore()</code> in the spec means <code>NewStore()</code> in the tests. The blueprint is the single source of truth.

<pre class="mermaid">
flowchart TD
    B["📄 Spec\n(signatures only)"] --> D["Executor AI\nfresh call, no memory"]
    D --> E["Implementation\ntask.go"]
    E --> F{"go test\nhidden tests run"}
    C["🧪 Hidden Tests"] --> F
    F -->|"PASS"| G["✅ Verified:\nexecutor followed spec"]
    F -->|"FAIL"| H["Error:\nexpected NewStore\ngot CreateStore"]
    H -->|"retry\nup to 3x"| D
    style G fill:#d4edda,stroke:#28a745
    style H fill:#fce4ec,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
</pre>

<strong>Phase 2 — Executor:</strong> A fresh model call receives <em>only</em> the spec. It builds the implementation. Then the hidden tests run as a surprise. If the executor invented its own API, the tests fail with something like <em>"undefined: CreateStore"</em> — and that exact error gets fed back for a retry.`),
				How:     `We ran this three times on <strong>Qwen3-30B</strong> ($0.0005 per call) targeting the task-board model layer. The blueprint generated a spec with signatures like <code>func NewStore() *Store</code>, <code>func (s *Store) Create(task Task) Task</code>, <code>func (s *Store) Get(id string) (Task, error)</code> — plus a test file calling each one.<br><br>The executor received only the spec. It generated <code>task.go</code>. Then the hidden tests ran: <code>go test ./model/... -count=1</code>. On the first run, the executor used <code>New()</code> instead of <code>NewStore()</code> — the hidden test caught it instantly. The error was fed back, and the retry fixed it. Runs 2 and 3 passed first time.`,
				Impact:  `<strong>100% — three out of three runs passed</strong>, all on the cheapest model. The hidden tests caught exactly the kind of mismatch they were designed for, and the feedback loop fixed it in one retry.<br><br>This became <strong>Step 5 of the pipeline</strong> (Exact Type Signatures). The critical learning: the spec must contain <em>actual Go code</em>, not English descriptions. "Create a function that makes a new store" fails. <code>func NewStore() *Store</code> succeeds. Precision in the spec is free and eliminates an entire failure category.<br><br>Combined with <a href="/exp/16">Experiment 16</a> (coupled files together) and this experiment (decoupled files via spec contract), we now had a complete strategy for multi-file generation.`,
				Related: []string{"10", "16", "27"},
			},
			{Num: "18", NumID: 18, Focus: "Full Pipeline E2E", Result: "0% (main.go)", Cost: "$0.167", Finding: "Prompt hint works in isolation, fails in context", HasDetail: true, SourceFile: "exp18-full-pipeline.py",
				Why: template.HTML("We had all the pieces: Exp 15's prompt hints fixed the backtick problem (100%), Exp 16's 2-file granularity fixed type mismatches (100%), the auto-fix pipeline recovered 40-60% of build failures, and Parser v2 handled 8 output formats. Time to assemble the full pipeline and build the complete task-board app end-to-end on the cheapest model. If each piece works at 100%, the whole pipeline should work too. Right?"),
				What: template.HTML("End-to-end build of the complete task-board application in 3 sub-tasks, all on <strong>Qwen3-30B</strong> ($0.0005/call): <strong>ST-1:</strong> schema.graphql (trivial, single file). <strong>ST-2:</strong> model/task.go + model/task_test.go (2-file pattern from Exp 16, with auto-fix). <strong>ST-3:</strong> main.go (HTTP server with embedded HTML, using the CORRECT/WRONG backtick hints from Exp 15). Five complete runs with up to 2 retries per sub-task. Gate: <code>go test ./...</code> must pass."),
				How: template.HTML("Each run builds the app from scratch. ST-1 and ST-2 run on the API with Parser v2 extraction. ST-3 generates main.go with the full backtick hint prompt. After each sub-task: auto-fix pipeline runs (goimports &rarr; gofmt &rarr; fix-address-of-const &rarr; goimports &rarr; gofmt), then <code>go build</code> gate. If build fails, the error is fed back for retry. After all sub-tasks, the golden test suite runs: <code>go test ./... -count=1</code>."),
				Impact: template.HTML("<strong>0/5 runs succeeded.</strong> Schema and model+test passed every time, but main.go failed in all 5 runs. The backtick hint that worked perfectly in Exp 15's isolated test failed when embedded in the full pipeline's longer prompt. The model was overwhelmed by context — the architecture spec, the store API documentation, AND the backtick hint together exceeded what Qwen3-30B could track. Average cost: $0.033 per failed run, $0.167 total wasted.<br><br>This was a critical lesson: <strong>improvements that work in isolation can fail in combination</strong>. The prompt hint didn't stop working — it got lost in a larger prompt. This led directly to the hybrid architecture (Exp 20): use cheap API models for focused tasks (store layer) and <code>claude -p</code> for complex tasks (HTTP server with embedded HTML). The cheap model's context window is the bottleneck, not its capability."),
				Related: []string{"15", "16", "20"},
			},
			{Num: "19", NumID: 19, Focus: "V2 Re-run (dep-doctor)", Result: "94% compile", Cost: "$0.084", Finding: "Compile gate 0% to 94%; golden tests need planner", HasDetail: true},
			{Num: "20", NumID: 20, Focus: "URL Shortener (new app)", Result: "Store works", Cost: "$0.019", Finding: "Approach generalises; claude -p needs foreground", HasDetail: true, SourceFile: "exp20-bigger-app.py",
				Why: template.HTML("Every experiment so far had used the same application: a task-board kanban. Maybe our techniques only worked because we had accidentally overfitted to that one app. The real question: does the pipeline generalise to a completely different application with different data structures, different API endpoints, and different business logic?"),
				What: template.HTML("A <strong>URL shortener</strong> — a completely different application from the task-board. Features: shorten URLs, redirect via short code, list all URLs, delete URLs, track click statistics. Two-phase build: <strong>Phase 1:</strong> Qwen3-30B generates store/store.go + store/store_test.go (2-file pattern). <strong>Phase 2:</strong> <code>claude -p --model haiku</code> reads the store layer and builds main.go with REST API + HTML frontend. Three runs."),
				How: template.HTML("Phase 1 uses the proven formula: exact type signatures in the prompt (<code>type URL struct</code> with ID, LongURL, ShortCode, CreatedAt, Clicks fields), exact function signatures (<code>NewStore()</code>, <code>Shorten()</code>, <code>Resolve()</code>, <code>List()</code>, <code>Delete()</code>, <code>Stats()</code>), 2-file output with auto-fix. Phase 2 gives <code>claude -p</code> the store API and asks it to build an HTTP server with JSON endpoints and a simple HTML form. The server is tested by actually starting it and hitting the endpoints: POST /shorten, GET /r/:code (redirect), GET /api/urls, DELETE /api/urls/:code."),
				Impact: template.HTML("The store layer built and tested successfully — the 2-file pattern + exact type signatures generalise perfectly to a new application. The HTTP server phase had mixed results: <code>claude -p</code> produced working code, but the process revealed a critical operational constraint — <strong><code>claude -p</code> requires a foreground terminal</strong>. It cannot run as a background subprocess, which is a problem for pipeline automation. Average cost: $0.019 per run. The key takeaway: the pipeline approach works for any CRUD application, not just the task-board. The store-first, server-second architecture is genuinely general-purpose."),
				Related: []string{"12", "18", "21"},
			},
			{Num: "21", NumID: 21, Focus: "StatusPulse (4 services)", Result: "4/4 build", Cost: "$0.021", Finding: "1,540 lines, 4 microservices, $0.02 total", HasDetail: true, SourceFile: "exp21-build-stores.py",
				Why: template.HTML("Experiments 20 proved the pipeline works for a second app. But both apps so far were single-service. Real products have multiple services — a monitoring service, an incident tracker, a notification system, a public status page. Could we build all four store layers for a multi-service architecture in one automated run? And what would it cost?"),
				What: template.HTML("Build the complete store layer for <strong>StatusPulse</strong>, a status page monitoring platform with 4 independent microservices: <strong>Monitor</strong> (health checks, latency tracking), <strong>Incidents</strong> (incident lifecycle, timeline entries, severity), <strong>Notifications</strong> (subscriber management, alert delivery), and <strong>Gateway</strong> (public status page, service aggregation). Each service gets its own store + tests using the 2-file pattern. All on <strong>Qwen3-30B</strong> at $0.0005/call."),
				How: template.HTML("Each service has exact type definitions baked into the prompt — not just names, but complete Go struct declarations with JSON tags, const blocks for enums, and exact function signatures. For example, the Monitor service specifies <code>type CheckStatus string</code> with constants <code>StatusUp</code>, <code>StatusDown</code>, <code>StatusUnknown</code>, plus 7 exact function signatures. The prompt includes the <code>&amp;constant</code> warning. Each service builds independently: generate 2 files &rarr; auto-fix pipeline &rarr; <code>go build</code> &rarr; <code>go test</code>."),
				Impact: template.HTML("<strong>4/4 services built and compiled. 1,540 lines of Go across 8 files. Total cost: $0.021.</strong> That is four complete microservice store layers — with types, CRUD operations, thread safety, and tests — for two cents. The exact-types approach scaled linearly: each service was as reliable as a single-service build because the prompts were self-contained. This proved that the pipeline is not limited to toy apps — it can scaffold a real multi-service architecture. The store layers from this experiment became the foundation for Exp 22 (screen design) and Exp 24 (wireframes to code), where the StatusPulse gateway got a full UI."),
				Related: []string{"20", "22", "24"},
			},
		},
	},
	{
		Slug:     "product-design",
		Name:     "Product Design",
		Icon:     "\U0001f3a8",
		ExpRange: "Experiments 22 – 26",
		Narrative: []template.HTML{
			`Code generation alone produces utilities, not products. This series explored whether AI can also handle the product design layer: user journeys, persona interviews, wireframes, and dev review gates. The pipeline that emerged &mdash; Brief, Personas, Screens, Dev Review &mdash; turns an 8-word idea into a complete, reviewed product specification for about $0.04.`,
			`The most striking result came from <strong>persona interviews</strong> (Experiment 23). Four AI personas representing different user archetypes were asked to evaluate a CRM invoicing brief. Two approved, two rejected &mdash; and the rejecting personas identified critical missing features like "recurring invoices" that were absent from the original brief. This is genuine requirements discovery at near-zero cost.`,
			`Wireframes-to-code (Experiment 24) showed that screen-level design constraints produce dramatically more focused implementations. The wireframe-guided version was 106 lines versus 493 lines for unconstrained generation &mdash; less than a quarter of the code, with the same feature set. The dev review gate (Experiment 25) catches complexity problems like PDF generation, SMTP integration, and PII handling before any code is written, at $0.007 per review.`,
		},
		KeyInsight: `$0.04 of design work transforms a code generator into a product builder. Persona interviews surface requirements that even experienced product managers miss.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "22", NumID: 22, Focus: "Brief to Journeys to Screens", Result: "8+ screens", Cost: "$0.038", Finding: "$0.04 design turns code into a product", HasDetail: true, SourceFile: "exp22-journeys-to-screens.py",
				Why: template.HTML("In Exp 21, <code>claude -p</code> built the StatusPulse gateway with a dark-themed dashboard — and it looked fine. But nobody had told it what screens to build. It guessed. It made a single page with check cards and an incident sidebar. No add-check form. No incident detail page. No subscriber management. No settings. The AI built what it imagined, not what users need. We had a code generator, but not a product builder. The missing layer was design: before writing code, someone needs to decide what screens exist and what goes on them."),
				What: template.HTML("A three-step design pipeline that turns a one-line brief into detailed text wireframes: <strong>Step 1:</strong> Generate user personas and map their journeys (what they do, what they see, what they click). <strong>Step 2:</strong> Extract a screen map — every unique screen the product needs, with URLs, components, data, and navigation. <strong>Step 3:</strong> Create ASCII text wireframes for each screen with layout, sample data, and interactive elements. Brief: <em>\"Build a public status page for monitoring website uptime, managing incidents, and notifying subscribers.\"</em> All steps on Qwen3-30B."),
				How: template.HTML("Step 1 uses temperature 0.4 (slightly higher for creative work) and asks for specific journey detail: <em>\"Admin clicks 'Add Check', enters URL and name, sees green confirmation toast\"</em> rather than <em>\"User adds a check.\"</em> Step 2 maps journeys to screens: every screen gets an ID (S01, S02...), a URL route, components list, data requirements, and navigation links. Step 3 produces ASCII wireframes with box-drawing characters showing exact layout, sample data, empty states, and interactive elements marked with [brackets]. The wireframes include responsive notes and color hints."),
				Impact: template.HTML("<strong>8+ screens identified from a single sentence.</strong> The brief mentioned monitoring, incidents, and subscribers — the design pipeline produced: dashboard (S01), add check form (S02), check detail (S03), incident list (S04), new incident (S05), incident detail (S06), subscriber management (S07), and settings (S08). Each with detailed wireframes. Total cost: $0.038. Compare this to Exp 21 where <code>claude -p</code> guessed and built one page. The $0.04 design investment transforms a code generator into a product builder. These wireframes fed directly into <a href=\"/exp/24\">Experiment 24</a> where they guided code generation to produce dramatically more focused output (106 lines vs 493)."),
				Related: []string{"21", "24", "23"},
			},
			{Num: "23", NumID: 23, Focus: "Persona Interview Loop", Result: "2 accept, 2 reject, APPROVED", Cost: "$0.051", Finding: "Personas found features NOT in brief", HasDetail: true, SourceFile: "exp23-persona-interviews.py",
				Why: template.HTML("Experiment 22 showed that design turns a brief into screens. But who decides which features matter? In Exp 22, the AI designed screens based on what seemed reasonable. Nobody pushed back. Nobody said <em>\"Where's recurring invoices?\"</em> or <em>\"I need to track whether clients opened the invoice.\"</em> A human product manager would interview potential users before building. Could AI personas do the same thing — and would they actually find features the brief missed?"),
				What: template.HTML("A four-step persona interview loop: <strong>Step 1:</strong> Generate 4 diverse personas from the brief (<em>\"Build an invoice generator for freelancers\"</em>). <strong>Step 2:</strong> Interview each persona in-character — 10 specific questions about their workflow, pain points, pricing tolerance, integration needs, and trust concerns. <strong>Step 3:</strong> Synthesise interviews into a prioritised feature matrix (Must-Have / Should-Have / Nice-to-Have). <strong>Step 4:</strong> Each persona reviews the proposed MVP and votes: ACCEPT, ACCEPT WITH CONCERNS, or REJECT. All on Qwen3-30B at $0.0005/call."),
				How: template.HTML("Each persona gets a realistic background: a solo freelance designer, a small agency owner, an accountant managing multiple clients, and a new freelancer just starting out. The interview prompt forces specificity: <em>\"Don't just say 'I need invoicing' — say 'I need to send a PDF invoice with my logo to my client within 5 minutes of finishing a project, and I need to track whether they've opened it.'\"</em> After synthesis, each persona reviews the MVP plan against their original stated needs and votes. If a majority accepts, the plan passes. If rejected, the feedback is incorporated and the plan is revised for a second vote."),
				Impact: template.HTML("<strong>2 accepted, 2 rejected — and the rejections were more valuable than the acceptances.</strong> The rejecting personas identified features not in the original brief: recurring invoices (demanded by 3 of 4 personas), payment tracking, tax calculation, and multi-currency support. After incorporating feedback, all four approved the revised plan. Total cost: $0.051.<br><br>This became <strong>Step 1 of the pipeline</strong> (Persona Discovery). The key insight: a brief is always incomplete. Human product managers miss things too, but persona interviews cost $0.01 and take 30 seconds. The rejection feedback is the most valuable output — it surfaces real requirements gaps before any code is written. Skipping this step means building the wrong product."),
				Related: []string{"22", "25", "32"},
			},
			{Num: "24", NumID: 24, Focus: "Wireframes to Code", Result: "6 routes, BUILD PASS", Cost: "$0.20", Finding: "106 lines vs 493 — wireframes give focus", HasDetail: true, SourceFile: "exp24-wireframes-to-code.py",
				Why: template.HTML("Experiment 22 produced detailed wireframes for StatusPulse. Experiment 21 produced a StatusPulse gateway where <code>claude -p</code> guessed the UI. Now we had a direct comparison to make: does feeding wireframes to the code generator actually produce better output? Or does the AI generate roughly the same thing regardless of input?"),
				What: template.HTML("Feed Exp 22's text wireframes and screen map directly to <code>claude -p --model haiku</code> and ask it to build main.go for the StatusPulse gateway. The prompt includes the exact screen map (which screens exist, their URLs, what data they show) and ASCII wireframes (layout, components, interactive elements). Compare the output against Exp 21's gateway, which was built without any design input. Gate: <code>go build</code> passes, all specified routes respond."),
				How: template.HTML("The prompt gives <code>claude -p</code> the screen map and wireframes (truncated to fit context), then specifies 6 routes matching the wireframes: GET / (dashboard), GET /admin/checks/new, GET /incidents, GET /incidents/new, GET /subscribers, and GET /settings. Each route description references the wireframe: <em>\"Main dashboard matching wireframe S01 — header with StatusPulse branding, service check cards showing status, active incidents section, Add Check button.\"</em> The model reads the wireframes, builds main.go, and runs <code>go build</code> to verify."),
				Impact: template.HTML("<strong>106 lines vs 493 lines — less than a quarter of the code, with the same feature set.</strong> The wireframe-guided version was dramatically more focused. Instead of inventing a complex dark-themed dashboard with animated gradients (Exp 21), it built exactly what the wireframes specified: clean cards, status indicators, forms with the right fields, proper navigation between screens. The 6 routes all responded correctly. Cost: $0.20.<br><br>This proved that wireframes act as a <strong>scope constraint</strong>. Without them, the AI over-engineers — adding features, animations, and complexity nobody asked for. With them, it builds exactly what's specified and nothing more. Wireframes became <strong>Step 3 of the pipeline</strong>. The 4:1 code reduction also means fewer bugs, faster builds, and cheaper iteration."),
				Related: []string{"22", "21", "25"},
			},
			{Num: "25", NumID: 25, Focus: "Full Pipeline: Brief to Product", Result: "Store + Server PASS", Cost: "$0.21", Finding: "8 words to compiled app, dev review caught issues", HasDetail: true, SourceFile: "exp25-full-pipeline.py",
				Why: template.HTML("We had proven each piece independently: persona interviews (Exp 23) find missing features, wireframes (Exp 22) structure the UI, wireframe-guided code (Exp 24) produces focused output, and the 2-file pattern (Exp 16) generates reliable store layers. But nobody had chained the entire sequence — brief to personas to MVP to wireframes to dev review to code — in one automated run. Would the pieces compose? And what would the dev review gate actually catch?"),
				What: template.HTML("The complete idea-to-product pipeline in one script. Brief: <em>\"Build an invoice generator for freelancers.\"</em> Six phases: <strong>(1)</strong> Generate 3 personas and interview them. <strong>(2)</strong> Synthesise an MVP with feature prioritisation. <strong>(3)</strong> Design screen wireframes from the MVP. <strong>(4)</strong> Dev review — 3 senior engineer personas (Backend Architect, Cost &amp; Complexity, Security &amp; Ops) review the spec. <strong>(5)</strong> Revise spec from dev feedback if needed. <strong>(6)</strong> Build store layer (Qwen3-30B) + HTTP server (<code>claude -p</code> Haiku)."),
				How: template.HTML("Each phase feeds into the next. Persona interviews output a feature list. The MVP synthesis categorises features as Must-Have / Should-Have / Deferred. Wireframes specify every screen with ASCII layouts. The dev review is the new gate: three engineers examine the spec for over-engineering, hidden complexity, and security risks. The Backend Architect checks architecture decisions. The Cost reviewer flags expensive features (PDF generation, SMTP). The Security reviewer checks data handling. If 2/3 approve, code generation proceeds. If not, the spec is revised and re-reviewed."),
				Impact: template.HTML("<strong>Store + server both compiled.</strong> The dev review caught real issues: one engineer flagged PDF generation as expensive for v1 and recommended HTML-to-print instead. Another identified that in-memory storage with no persistence means data loss on restart — acceptable for MVP but must be documented. The security reviewer raised GDPR concerns with invoice data. Total cost: $0.21.<br><br>This validated the full pipeline sequence: Brief &rarr; Personas &rarr; MVP &rarr; Wireframes &rarr; Dev Review &rarr; Code. Each step costs pennies and catches different classes of problems. The dev review gate became <strong>Step 4 of the pipeline</strong>. The next experiment (Exp 26) added test layers to this pipeline, and Exp 27 ran it fully automated end-to-end."),
				Related: []string{"23", "24", "26", "27"},
			},
			{Num: "26", NumID: 26, Focus: "Add Test Layers", Result: "33/36 (92%)", Cost: "$0.12", Finding: "Store + acceptance + HTTP tests, $0.33 total", HasDetail: true, SourceFile: "exp26-add-tests.py",
				Why:     `Experiment 25 produced a compiled invoice generator, but it had zero tests. We could run it, but we had no idea if it actually worked correctly. A build that passes <code>go build</code> is not the same as a build that passes <code>go test</code>. We needed three layers of confidence: do the store functions return the right data? Does the implementation match the spec? Do the HTTP endpoints respond correctly?`,
				What:    `Three test layers added to the Exp 25 invoice generator, each targeting a different class of defect: <strong>Layer 1 — Store unit tests</strong> (Qwen3-30B, $0.0005/call): tests for every public function in the invoice store — create, get, list, update, delete, plus edge cases like missing fields and not-found errors. <strong>Layer 2 — V-Model acceptance tests</strong> (Qwen3-30B): generated from the spec, never shown to the dev model, verifying that the implementation matches the contract. <strong>Layer 3 — HTTP integration tests</strong> (<code>claude -p</code> Haiku): tests that start the server, hit endpoints, and verify responses end-to-end.`,
				How:     `Layer 1 uses the 2-file awareness approach — the model reads the actual store code and writes tests against the real function signatures. Auto-fix pipeline runs after generation. If tests fail, an AI fix loop reads the error, the store code, and the failing test, then produces a corrected test file — up to 3 retries. Layer 2 reads the Go type spec from Exp 25 and generates acceptance tests that verify <em>behavior</em>, not implementation details. Layer 3 uses <code>claude -p</code> to generate HTTP tests that exercise the full stack.`,
				Impact:  `<strong>33 out of 36 tests passing (92%).</strong> The 3 failures were in acceptance tests that expected stricter validation than the implementation provided — a genuine spec gap, not a test bug. The store unit tests caught 2 edge cases in the delete function. The HTTP tests confirmed all endpoints responded correctly. Total cost: $0.12 for 36 tests across 3 layers. This validated the three-layer testing strategy that became standard in the pipeline: unit tests for logic, acceptance tests for spec compliance, integration tests for the full stack.`,
				Related: []string{"25", "16", "17", "48"},
			},
		},
	},
	{
		Slug:     "review-pipeline",
		Name:     "Review Pipeline",
		Icon:     "\U0001f50d",
		ExpRange: "Experiments 27 – 39",
		Narrative: []template.HTML{
			`If code generation is the engine, reviews are the quality control line. This series scaled from 4 reviewers to 8 to 10 to 20, exploring whether more reviewers produce diminishing returns. The answer: they do not. At $0.005 per reviewer, every additional reviewer finds unique issues that no other reviewer catches. Twenty reviewers at $0.10 total is still cheaper than one hour of human review.`,
			`The <strong>domain expert reviewer</strong> (Experiment 32) was the breakout result. A single reviewer prompted with invoicing domain knowledge found all 10 missing workflow features &mdash; print, pay, void, line items, tax calculation, recurring billing &mdash; that generic "code quality" reviewers missed entirely. This demonstrates that review quality depends on domain framing, not model capability. The same model, with a domain expert persona, sees things it otherwise ignores.`,
			`<strong>Constrained simplicity</strong> (Experiment 36) combined the domain expert with a simplicity agent that was instructed to simplify implementation without removing features. The result: all features present, 60% cheaper than domain-only, and 29/29 tests passing. <strong>Progressive enhancement</strong> (Experiment 38) proved the most reliable build strategy overall: 5 iterations building one feature at a time, with test verification after each, producing zero regressions across 32 tests.`,
		},
		KeyInsight: `$0.005/reviewer with no diminishing returns at 20 reviewers. Domain framing is more important than model capability for finding real issues.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "27", NumID: 27, Focus: "Fully Automated CRM", Result: "41/41 (100%)", Cost: "$0.49", Finding: "720 lines, zero human intervention", HasDetail: true, Icon: "\U0001f3d7\ufe0f", SourceFile: "exp27-fully-automated.py",
				Why:     `Twenty-six experiments had proven individual pieces — V4 prompts work, auto-fix recovers failures, the V-Model keeps types consistent, 2-files-per-call keeps tests aligned. But nobody had run the whole pipeline end-to-end without a human touching anything. Could we type an 8-word brief and walk away, coming back to a compiled, tested, working application?`,
				What:    template.HTML(`A single command: <code>python3 exp27-fully-automated.py "CRM for freelancers with invoicing"</code>. The script runs every pipeline step autonomously:<br><br><strong>Phase 1 — Discover:</strong> Qwen3-30B ($0.0005) generates 3 personas and interviews them about the brief. Synthesises an MVP feature list with two senior engineers filtering scope.<br><br><strong>Phase 2 — Spec:</strong> Qwen3-30B produces wireframes and exact Go type signatures.<br><br><strong>Phase 3 — Build Store:</strong> Qwen3-30B generates store.go + store_test.go (2-file pattern from Exp 16). Auto-fix pipeline runs. If tests fail, an AI fix loop retries up to 3 times.<br><br><strong>Phase 4 — Build Server:</strong> <code>claude -p --model haiku</code> (free) reads the store layer and generates main.go with HTTP handlers + embedded HTML.<br><br><strong>Phase 5 — Test:</strong> <code>claude -p</code> generates HTTP integration tests. All tests run.

<pre class="mermaid">
flowchart TD
    A["8-Word Brief"] --> B["Phase 1: Discover"]
    B --> C["Phase 2: Spec"]
    C --> D["Phase 3: Build Store"]
    D --> E["Phase 4: Build Server"]
    E --> F["Phase 5: Test"]
    F --> G{{"41/41 Passing?"}}
    G -- Yes --> H["Working CRM"]:::success
    G -- No --> I["Pipeline Stops"]:::failure
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style C fill:#f5e1e3,stroke:#d4727a
    style D fill:#f5e1e3,stroke:#d4727a
    style E fill:#f5e1e3,stroke:#d4727a
    style F fill:#f5e1e3,stroke:#d4727a
    style G fill:#fff3cd,stroke:#d4a017
    classDef success fill:#d4edda,stroke:#28a745
    classDef failure fill:#fce4ec,stroke:#d4727a
</pre>`),
				How:     template.HTML(`The brief enters at the top and the pipeline runs unattended. Each phase has a gate: personas must produce a feature list, the spec must contain Go types, the store must pass <code>go test</code>, the server must pass <code>go build</code>, and all integration tests must pass. The auto-fix loop (goimports → gofmt → fix-address-of-const → go vet → go build) runs between every generation step. If a phase fails after 3 retries, the pipeline stops and reports where it died.<br><br>Cost: Qwen3-30B handles phases 1-3 at ~$0.01 each. Haiku (free on subscription) handles phases 4-5. Total under $0.50.

<pre class="mermaid">
flowchart TD
    A["AI Generates Code"] --> B["goimports + gofmt"]
    B --> C["fix-address-of-const"]
    C --> D["go vet + go build"]
    D --> E{{"Compiles?"}}
    E -- Yes --> F["Run Tests"]:::success
    E -- No --> G{{"Retries &lt; 3?"}}
    G -- Yes --> H["AI Fix Loop"]:::pink
    H --> B
    G -- No --> I["Pipeline Stops"]:::failure
    style A fill:#f5e1e3,stroke:#d4727a
    style E fill:#fff3cd,stroke:#d4a017
    style G fill:#fff3cd,stroke:#d4a017
    classDef pink fill:#f5e1e3,stroke:#d4727a
    classDef success fill:#d4edda,stroke:#28a745
    classDef failure fill:#fce4ec,stroke:#d4727a
</pre>`),
				Impact:  `<strong>41/41 tests passing. 720 lines of Go. Zero human intervention. $0.49 total.</strong> This was the proof-of-concept that the pipeline works. An 8-word brief produced a working CRM with client management, activity tracking, invoicing, and a full HTML frontend. The entire run took about 4 minutes.<br><br>This script became the reference implementation for the composable block architecture — each phase maps directly to a pipeline block with defined inputs, outputs, gates, and failure handling.`,
				Related: []string{"16", "17", "7", "37"},
			},
			{Num: "29", NumID: 29, Focus: "CRM + gqlgen GraphQL", Result: "BUILD PASS", Cost: "~$0.50", Finding: "Pipeline needs iterative mode for codegen", HasDetail: true,
				Why:     `Every app we had built so far used REST + embedded HTML. That is fine for simple CRUDs, but real SaaS products often use GraphQL for flexible querying and frontend decoupling. Could our pipeline produce a CRM with gqlgen — Go's standard schema-first GraphQL framework — or would the codegen layer add too much complexity for AI to handle?`,
				What:    `Full pipeline run (personas, MVP, wireframes, dev review, build) targeting a CRM with a <strong>gqlgen GraphQL API</strong> instead of REST. The schema defines Client, Activity, and Invoice types with queries and mutations. gqlgen generates resolver stubs from the schema, and the AI fills in the implementations. Model: Haiku via <code>claude -p</code>.`,
				How:     `The pipeline ran normally through design phases. At build time, the AI had to: (1) write <code>schema.graphql</code> with type definitions, (2) run <code>go run github.com/99designs/gqlgen generate</code> to scaffold resolvers, (3) implement each resolver function, (4) wire up the HTTP server with the gqlgen handler. The codegen step required multiple iterations — gqlgen is strict about schema-resolver alignment, and each mismatch required a regeneration cycle.`,
				Impact:  template.HTML("<strong>Build passed, but the process was painful.</strong> gqlgen's codegen required the AI to iterate far more than a REST build — schema changes triggered regeneration, which invalidated previous resolver work. Total cost: ~$0.50, roughly 2x a REST equivalent. The key finding: the pipeline needed an <strong>iterative codegen mode</strong> where the AI can run a tool, inspect the output, and adjust, rather than one-shot generation. This experiment directly informed <a href=\"/exp/91\">Experiment 91</a>'s conclusion that hand-rolled minimal GraphQL is more reliable for AI generation than framework-based approaches."),
				Related: []string{"27", "91"},
			},
			{Num: "30", NumID: 30, Focus: "4-Reviewer Panel", Result: "51/51 (100%)", Cost: "$0.40", Finding: "Product reviewer caught missing Add/Delete", HasDetail: true, SourceFile: "exp30-full-reviewed.py",
				Why:     `Experiment 25 had a single dev reviewer who checked architecture and cost. But a backend engineer does not think like a product manager. The dev reviewer approved a spec where there was no "Add Client" button on the dashboard and no "Delete" confirmation dialog — the two most basic CRUD affordances. If a real user opened that app, their first question would be: "How do I add a client?" We needed reviewers who think about the product, not just the code.`,
				What:    `Four reviewers, each with a distinct persona and focus area: <strong>Dev Architect</strong> (architecture, cost, $5 VPS viability), <strong>Product Owner</strong> (every CRUD action has a button, every journey has a screen), <strong>QA Engineer</strong> (empty forms, duplicates, delete cascades, error states), <strong>Market Analyst</strong> (vs HubSpot Free, vs Zoho Free, differentiation, pricing). Same CRM brief as Exp 27. All reviews on Qwen3-30B, build on <code>claude -p</code> Haiku.`,
				How:     `The pipeline runs: personas, MVP synthesis, screen wireframes, then the 4-reviewer panel. Each reviewer receives the MVP spec and wireframes and returns APPROVE or REQUEST CHANGES with specific findings. The Product Owner is prompted to check every persona journey against the screens: <em>"Can they ADD a new client? Is there a button? Can they EDIT? Where? Can they DELETE? Where?"</em> After reviews, findings are incorporated into the spec, then the store + server are built with the auto-fix loop.`,
				Impact:  `<strong>51/51 tests passing (100%).</strong> The Product Owner was the standout: it found that the original wireframes had no Add Client button on the dashboard, no Delete confirmation, and no way to view activity history — all critical UX gaps. The QA reviewer found 5 missing error states. The Market Analyst flagged that the pricing had no free tier, which every competitor offers. Total cost: $0.40. This established the multi-reviewer panel as a standard pipeline step and proved that reviewer <em>persona</em> determines what gets caught. The same model with a Product Owner hat finds UX gaps; with a Dev hat, it only sees architecture.`,
				Related: []string{"25", "32", "37", "39"},
			},
			{Num: "31", NumID: 31, Focus: "Post-Code: UX + A11y + OWASP", Result: "54/54 tests", Cost: "$0.025", Finding: "ARIA labels, XSS fixes, CSP headers added", HasDetail: true, SourceFile: "exp31-post-code-review.py",
				Why:     `The pre-code reviewers in Exp 30 examined the spec. But a spec review cannot catch issues that only exist in the built HTML — missing ARIA labels, XSS vectors in rendered templates, absent Content-Security-Policy headers. These are defects in the <em>implementation</em>, not the design. We needed reviewers who examine the actual built code and rendered HTML output, not the plan.`,
				What:    `Three post-code reviewers examine the actual main.go and its embedded HTML from Exp 30's CRM: <strong>UI/UX Expert</strong> (visual hierarchy, layout, forms, navigation, responsive design, whitespace, typography), <strong>Accessibility Specialist</strong> (WCAG 2.1 AA compliance — ARIA labels, semantic HTML, keyboard navigation, color contrast, heading hierarchy), <strong>OWASP Security Reviewer</strong> (XSS prevention, CSRF tokens, CSP headers, input validation, rate limiting). All on Qwen3-30B. Fixes applied via <code>claude -p</code> Haiku.`,
				How:     `The script extracts embedded HTML from main.go (the Go backtick template strings) and sends it to each reviewer. Each reviewer returns severity-rated findings with specific HTML/CSS fixes. The UX reviewer checks 10 categories including empty states and interaction feedback. The accessibility reviewer checks against specific WCAG criteria (1.3.1, 4.1.2, etc.). The OWASP reviewer checks for XSS, CSRF, and missing security headers. Findings are collected, then <code>claude -p</code> applies all fixes to the running codebase and rebuilds.`,
				Impact:  template.HTML("<strong>54/54 tests passing after fixes.</strong> The accessibility reviewer added ARIA labels to all form elements and a skip-navigation link. The OWASP reviewer found two XSS vectors where user-generated content was rendered without escaping and added Content-Security-Policy headers. The UX reviewer improved form feedback and empty states. Total cost: $0.025 for three reviews. This became <strong>Step 8 of the pipeline</strong> (Post-Code Review). The critical lesson: these are the kinds of issues that slip through functional testing because the app \"works\" without them — they only matter for quality, accessibility, and security."),
				Related: []string{"30", "11", "32", "42"},
			},
			{Num: "32", NumID: 32, Focus: "8 Reviewers + Domain Expert", Result: "22/22 HTTP", Cost: "$0.96", Finding: "Domain expert found 10 missing invoice features", HasDetail: true, SourceFile: "exp32-domain-expert.py",
				Why:     `The 4-reviewer panel (Exp 30) caught missing CRUD operations — no Add button, no Delete confirmation. But when we showed the resulting CRM to someone who actually does freelance invoicing, their first question was: "Where's the print button? How do I mark an invoice as paid? Can I add line items?" The generic reviewers — architect, product owner, QA, market analyst — had approved a product that was missing the most basic invoicing workflow features. We needed a reviewer who actually understood the domain.`,
				What:    template.HTML(`We added a <strong>domain expert reviewer</strong> to the panel — an AI persona prompted with deep knowledge of freelance invoicing workflows. The full panel became 8 reviewers: 5 pre-code (architect, product owner, QA engineer, market analyst, <strong>invoicing domain expert</strong>) and 3 post-code (UX, accessibility, OWASP security). The brief: <em>"Build a CRM for freelancers with client management, activity history, invoicing, and address book."</em> Model: <strong>Qwen3-30B</strong> for reviews, <strong>Haiku</strong> for code generation.

<pre class="mermaid">
flowchart TD
    S["CRM Spec"] --> PRE
    subgraph PRE["Pre-Code Review — 5 Reviewers"]
        R1["Dev Architect"]
        R2["Product Owner"]
        R3["QA Engineer"]
        R4["Market Analyst"]
        R5["Domain Expert"]:::pink
    end
    PRE --> REV["Revised Spec"]
    REV --> CODE["Code Generation"]
    CODE --> BUILD["Built Application"]
    BUILD --> POST
    subgraph POST["Post-Code Review — 3 Reviewers"]
        R6["UX Expert"]
        R7["Accessibility"]
        R8["OWASP Security"]
    end
    POST --> FINAL["22/22 HTTP Tests"]:::success
    style S fill:#f5e1e3,stroke:#d4727a
    style REV fill:#f5e1e3,stroke:#d4727a
    style CODE fill:#f5e1e3,stroke:#d4727a
    style BUILD fill:#f5e1e3,stroke:#d4727a
    style R5 fill:#f5e1e3,stroke:#d4727a
    classDef success fill:#d4edda,stroke:#28a745
</pre>`),
				How:     `Each reviewer receives the spec and returns a structured review: APPROVE, APPROVE WITH FIXES, or REJECT, plus specific findings. The domain expert's prompt includes: <em>"You are a freelance consultant who has sent over 1,000 invoices. Review this specification for completeness of invoicing workflows: creation, editing, sending, payment tracking, recurring, tax handling, line items, numbering, void/credit notes, and print/PDF."</em><br><br>Pre-code reviews run before any code is generated. Findings are incorporated into a revised spec. Then code is built, and post-code reviewers check the actual HTML output. The build uses the full auto-fix pipeline with up to 3 retry loops.`,
				Impact:  `The domain expert found <strong>all 10 missing invoicing features</strong> that every other reviewer missed: print invoice, mark as paid, void invoice, line items, tax calculation, recurring invoices, payment terms, late fees, credit notes, and sequential numbering. Not one of these appeared in any other reviewer's findings.<br><br>This was the clearest evidence that <strong>domain framing matters more than model capability</strong>. The same model (Qwen3-30B), given a domain expert persona, sees things it completely ignores with a generic "code quality" persona. This led to Step 4 of the pipeline always including at least one domain-specific reviewer. The cost was $0.96 total — expensive by research standards, but the one-shot 10-feature build approach was later replaced by progressive enhancement (Exp 38) which achieved the same quality for $0.20.`,
				Related: []string{"30", "36", "38", "39"},
			},
			{Num: "33", NumID: 33, Focus: "Add Feature (CSV Export)", Result: "55/55, no regression", Cost: "$0.23", Finding: "Added feature without breaking 52 existing tests", HasDetail: true, SourceFile: "exp33-add-feature.py",
				Why:     `Every experiment so far built greenfield applications from scratch. But real products evolve — you ship v1 and then add features. The question we had not answered: can the pipeline add a feature to an <em>existing</em> codebase without breaking what already works? This is harder than greenfield because the AI must understand existing code, modify it correctly, and preserve all existing behavior. A single wrong import or renamed function breaks the entire test suite.`,
				What:    `Add "Export clients to CSV" to Exp 30's CRM. The existing app has 52 passing tests across store and HTTP layers. The new feature requires: a GET /api/clients/export endpoint, CSV output with headers (ID, Name, Company, Email, Phone, Status), a download button on the dashboard, and new tests for the export functionality. The constraint: all 52 existing tests must still pass after the change.`,
				How:     template.HTML("A single <code>claude -p --model haiku</code> call that reads the existing main.go and store/store.go, adds the CSV export endpoint and download button, then writes new tests. The prompt is explicit about the non-regression requirement: <em>\"Keep all existing functionality working.\"</em> After the modification, <code>go test ./... -count=1 -v</code> runs the full suite — old tests plus new ones. The script counts passing tests before and after to detect any regression."),
				Impact:  `<strong>55/55 tests passing — 52 existing + 3 new, zero regressions.</strong> The AI read the existing codebase, understood the patterns (handler structure, store interface, test style), and added the feature in the same style. No existing test broke. Total cost: $0.23. This was the first proof that the pipeline handles enhancement, not just greenfield. It also validated a principle: if the existing code is clean and well-tested, the AI preserves its patterns naturally. This experiment directly led to progressive enhancement (Exp 38), which applies the same principle iteratively.`,
				Related: []string{"30", "38", "27"},
			},
			{Num: "34", NumID: 34, Focus: "Simplicity Agent", Result: "13/13, 518 lines", Cost: "$0.20", Finding: "70% less code, 79% cheaper than Exp 32", HasDetail: true, SourceFile: "exp34-simplicity-agent.py",
				Why:     `Experiment 32 built a CRM with 8 reviewers and domain expert input. The result was feature-complete but expensive ($0.96) and complex. We suspected the pipeline was over-engineering — adding PDF generation, SMTP integration, PostgreSQL — when simpler alternatives existed. Could an agent whose sole job is to push back on complexity reduce cost without removing features?`,
				What:    `A "Simplicity Agent" that reviews every stage output with one question: <em>"Is this the simplest way?"</em> The agent runs after each pipeline stage — persona synthesis, MVP, wireframes, spec — and categorises every item as KEEP (essential and already simple), CUT (not needed for v1), or SIMPLIFY (needed but over-engineered). Rules: in-memory store instead of database, browser print instead of PDF library, simple string status instead of state machine. Brief: same CRM as Exp 32.`,
				How:     `The simplicity agent is an additional Qwen3-30B call after each stage. It receives the stage output and applies hard rules: <em>"MVP means MINIMUM. If 3 features work, don't add a 4th. 1 screen that does 3 things is better than 3 screens that each do 1 thing."</em> The simplified output replaces the original before the next stage begins. After all design stages are simplified, the build proceeds with <code>claude -p</code> Haiku. All tests run against the simplified app.`,
				Impact:  `<strong>13/13 tests passing, 518 lines, $0.20 — 70% less code and 79% cheaper than Exp 32.</strong> The simplicity agent cut: PDF generation (use browser print), SMTP email (defer to v2), complex status workflow (simple string field), separate search page (inline filter). But it also cut features that personas had demanded — recurring invoices, payment tracking, multi-currency. This was a failure mode: the agent was too aggressive. It simplified <em>what</em> to build, not just <em>how</em> to build it. This led directly to the constrained simplicity approach in <a href="/exp/36">Experiment 36</a>, where the agent was instructed: "simplify HOW, never WHETHER."`,
				Related: []string{"32", "36", "37"},
			},
			{Num: "35", NumID: 35, Focus: "Playwright Journeys", Result: "3 pass, 5 fail", Cost: "FREE", Finding: "Simplicity agent cut features from the brief", HasDetail: true, SourceFile: "exp35-playwright-journeys.py",
				Why:     `Experiment 34's simplicity agent produced a 518-line CRM that passed 13 unit tests. But unit tests only verify that functions return correct values — they do not verify that a user can actually navigate the app, fill a form, and see results. We needed to simulate real user journeys through a browser to find out if the simplified CRM actually worked as a product.`,
				What:    `Automated Playwright browser testing using <code>claude -p</code> with Playwright MCP tools. Eight persona journeys tested against the running Exp 34 CRM at localhost:8080: add a client, edit a client, view activity history, create an invoice, search clients, delete a client, view dashboard, and export data. Each journey is a sequence of real browser actions — navigate, click, type, verify.`,
				How:     `For each journey, <code>claude -p --model haiku</code> receives the journey steps and uses Playwright browser tools to execute them: navigate to the app, click buttons, fill forms, check that expected elements appear. Each step is reported as PASS or FAIL with notes explaining what happened. The script collects results as JSON summaries. The app must be running locally before the test starts.`,
				Impact:  template.HTML("<strong>3 journeys passed, 5 failed.</strong> The failures exposed exactly what the simplicity agent had cut: invoice creation was missing (the agent had deferred it), activity history had no UI (store functions existed but no form to add activities), and search was removed entirely. The 13 passing unit tests had given false confidence — the code <em>worked</em> but the product was incomplete. This was the definitive proof that Exp 34's simplicity agent was too aggressive: it had cut features the personas demanded, not just simplified implementations. The 5 failures became the spec for <a href=\"/exp/36\">Experiment 36</a>'s constrained simplicity approach. Playwright testing became <strong>Step 11 of the pipeline</strong>."),
				Related: []string{"34", "36", "40"},
			},
			{Num: "36", NumID: 36, Focus: "Domain + Constrained Simplicity", Result: "29/29, 1,191 lines", Cost: "$0.39", Finding: "All features, simply built, 60% cheaper", HasDetail: true, SourceFile: "exp36-constrained-simplicity.py",
				Why:     `Experiments 34 and 35 proved that an unconstrained simplicity agent cuts features users actually need. Experiment 32 proved that a domain expert finds all the right features but the resulting app is expensive and complex. We needed both: the domain expert defines <em>what</em> to build, and the simplicity agent decides <em>how</em> to build it simply. The constraint: the simplicity agent cannot remove any feature the domain expert specified.`,
				What:    `Two-agent architecture: <strong>Stage 1 — Domain Expert</strong> (Qwen3-30B) defines the full feature set for a freelancer CRM with invoicing, client management, activity history, and address book. <strong>Stage 2 — Constrained Simplicity Agent</strong> (Qwen3-30B) reviews the feature list with one rule: <em>"simplify HOW features are implemented, never WHETHER they are included."</em> In-memory instead of Postgres. Browser print instead of PDF library. Textarea instead of rich-text editor. Build on <code>claude -p</code> Haiku.`,
				How:     `The domain expert runs first and produces a comprehensive feature list — line items on invoices, recurring billing, payment terms, void/credit notes, sequential numbering. The constrained simplicity agent then reviews each feature and proposes the simplest implementation: line items become a JSON array in a text field, recurring billing becomes a "repeat" checkbox with next-due date, payment terms become a dropdown (Net 15/30/60), print becomes <code>window.print()</code>. The store layer is built with the 2-file pattern on Qwen3-30B. The HTTP server is built with <code>claude -p</code>.`,
				Impact:  `<strong>29/29 tests passing, 1,191 lines, $0.39 — all features present, 60% cheaper than Exp 32's domain-only approach ($0.96).</strong> Every feature the domain expert specified was present and working. The constrained simplicity agent had found simpler implementations without removing anything. This became <strong>Step 2 of the pipeline</strong> (MVP Synthesis). The key principle: separate the "what" decision (domain expert) from the "how" decision (simplicity agent). Never let a simplicity-focused agent decide what to build — only how to build it.`,
				Related: []string{"32", "34", "35", "37"},
			},
			{Num: "37", NumID: 37, Focus: "10-Reviewer Panel", Result: "26/26, 426 lines", Cost: "$0.34", Finding: "7 pre-code + 3 post-code, SaaS patterns included", HasDetail: true, SourceFile: "exp37-full-panel.py",
				Why:     `We had proven 4 pre-code reviewers (Exp 30) and 3 post-code reviewers (Exp 31) separately. Experiment 32 added a domain expert and 36 added constrained simplicity. But nobody had combined all of them into a single pipeline run. Would 10 reviewers produce a better product than 4? Or would conflicting feedback create chaos — one reviewer demanding complexity while another demands simplicity?`,
				What:    `Ten reviewers split into two phases. <strong>Pre-code (7):</strong> Dev Architect, Product Owner, QA Engineer, Market Analyst, Domain Expert (invoicing), SaaS UX Designer, Constrained Simplicity Agent. <strong>Post-code (3):</strong> Code Architecture Reviewer, Accessibility Specialist, OWASP Security Reviewer. Same CRM brief. All reviews on Qwen3-30B ($0.005/reviewer), build on <code>claude -p</code> Haiku.`,
				How:     `The pipeline runs design phases (personas, MVP, wireframes), then the 7 pre-code reviewers examine the spec in parallel. Each returns findings with severity ratings. The SaaS UX Designer adds professional patterns: breadcrumbs, toast notifications, confirmation dialogs for destructive actions. The Constrained Simplicity Agent ensures features are kept but implementations stay lean. After incorporating pre-code feedback, the store + server are built. Then the 3 post-code reviewers examine the actual HTML output for accessibility, security, and architecture issues.`,
				Impact:  `<strong>26/26 tests passing, 426 lines, $0.34.</strong> The 10-reviewer product was notably more polished than the 4-reviewer version: breadcrumbs on every page, toast notifications on form submission, confirmation dialogs before deletes, and responsive table layouts. The SaaS UX reviewer drove most of these improvements. The Constrained Simplicity Agent and Domain Expert did not conflict — they naturally operated on different axes (how vs what). The post-code reviewers added ARIA labels and CSP headers. This became the reference configuration for the pipeline review step and fed directly into the 20-reviewer scaling test (Exp 39).`,
				Related: []string{"30", "32", "36", "39"},
			},
			{Num: "38", NumID: 38, Focus: "Progressive Enhancement", Result: "32/32, ZERO regressions", Cost: "$0.20", Finding: "5 iterations, zero regressions vs one-shot", HasDetail: true, Icon: "\U0001f4c8", SourceFile: "exp38-progressive.py",
				Why:     `Experiment 32 tried to build a full CRM with all features in one shot. It cost $0.96, took 8 reviewers, and still needed extensive fixing. The problem wasn't the model or the reviewers — it was the approach. Asking an AI to generate 10 features at once is like asking a junior developer to build an entire product in one commit. Things get tangled. Features interfere with each other. A bug in the invoice code breaks the client list. We wanted to know: what happens if we build one feature at a time, testing after each addition?`,
				What:    `Five iterations, each adding one feature to a working application:<br><br><strong>Iteration 1:</strong> Client CRUD — add, list, view clients<br><strong>Iteration 2:</strong> + Edit and delete clients<br><strong>Iteration 3:</strong> + Activity log — add activities, view timeline per client<br><strong>Iteration 4:</strong> + Simple invoices — create, list, mark as paid<br><strong>Iteration 5:</strong> + Search across clients and print invoice view<br><br>Each iteration uses <code>claude -p --model haiku</code> (free on subscription). The AI reads the existing code before making changes, so it understands what's already built. After each iteration, <em>all</em> tests run — not just the new ones.`,
				How:     template.HTML(`The script starts with a bare Go module and an empty main.go. Each iteration sends a prompt like: <em>"Read main.go and store/store.go. Add activity logging: a form to add activities per client (type: call/email/meeting, notes, date) and a timeline view. Keep all existing functionality working. Run go build."</em><br><br>The key constraint: <code>claude -p</code> reads the existing code first, so it builds on what's there rather than rewriting from scratch. After each iteration, the gate runs all existing tests plus any new ones. If a previous test breaks, the iteration fails — that's a regression.<br><br>Budget: $0.30 max per iteration, $0.20 actual total across all 5.

<pre class="mermaid">
flowchart TD
    A["Bare Go Module"] --> B["Build Feature N"]
    B --> C["Run ALL Tests"]
    C --> D{{"All Pass?"}}
    D -- Yes --> E{{"More Features?"}}
    E -- Yes --> B
    E -- No --> F["32/32 Complete"]:::success
    D -- No --> G["Fix Regression"]:::failure
    G --> C
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#fff3cd,stroke:#d4a017
    classDef success fill:#d4edda,stroke:#28a745
    classDef failure fill:#fce4ec,stroke:#d4727a
</pre>`),
				Impact:  `<strong>32/32 tests passing. Zero regressions. $0.20 total.</strong> Not a single previously-passing test broke across 5 iterations. Compare this to Experiment 32's one-shot approach: $0.96 for the same feature set, with test failures requiring manual intervention.<br><br>Progressive enhancement is now <strong>Step 10 of the pipeline</strong> and the default build strategy. It's cheaper (4x), more reliable (zero regressions vs multiple), and produces cleaner code because each iteration builds on a working foundation. The trade-off is that it's slower — 5 iterations take ~5 minutes vs ~1 minute for one-shot — but the reliability difference is overwhelming.`,
				Related: []string{"32", "33", "27"},
			},
			{Num: "39", NumID: 39, Focus: "20 Reviewers", Result: "No diminishing returns", Cost: "$0.047", Finding: "More reviewers = more unique findings", HasDetail: true, SourceFile: "exp39-20-reviewers.py",
				Why:     `We had scaled from 1 reviewer (Exp 25) to 4 (Exp 30) to 10 (Exp 37). Each time, more reviewers found more issues. But conventional wisdom says review panels hit diminishing returns — after 5-7 reviewers, new reviewers repeat what earlier ones found. At $0.005 per reviewer, the cost was trivial, but we needed data: at what point does adding reviewers stop finding new issues?`,
				What:    `Twenty reviewers with maximally diverse perspectives review the same CRM spec: the standard 10 (Dev Architect, Product Owner, QA, Market Analyst, Domain Expert, SaaS UX, Constrained Simplicity, Code Architect, Accessibility, OWASP Security) plus 10 new personas — Freelance Designer, Agency Owner, First-Time User, Power User (500 clients), Mobile User, Competitor User (switching from HubSpot), Data Privacy Officer, Billing Specialist, Onboarding Specialist, Growth Hacker. All on Qwen3-30B.`,
				How:     `The script generates an MVP + wireframes, then runs all 20 reviewers sequentially. Each reviewer receives the spec and returns structured findings with severity ratings. Critically, each reviewer also self-reports which findings are UNIQUE to their perspective. The script tracks cumulative issues and unique findings at each reviewer count (1, 5, 10, 15, 20) to build a diminishing-returns curve. Total issues from the first 10 are compared against the last 10 to measure falloff.`,
				Impact:  template.HTML("<strong>No diminishing returns at 20 reviewers.</strong> Every reviewer found issues no other reviewer caught. The Data Privacy Officer found GDPR gaps invisible to the OWASP reviewer. The Power User found performance concerns (\"500 clients with no pagination?\") nobody else raised. The Growth Hacker asked \"What is the viral loop?\" — a question no technical reviewer would think to ask. The last 10 reviewers found as many unique issues as the first 10. Total cost: $0.047 for all 20 reviews. At $0.005/reviewer, there is no economic reason to cap the panel size. This became the basis for the Dark Factory production review system, which runs <a href=\"/exp/30\">10 reviewers</a> per PR."),
				Related: []string{"30", "32", "37"},
			},
		},
	},
	{
		Slug:     "testing-security",
		Name:     "Testing & Security",
		Icon:     "\U0001f6e1\ufe0f",
		ExpRange: "Experiments 40 – 52",
		Narrative: []template.HTML{
			`Tests and reviews catch different classes of bugs. This series explored the boundaries: what do unit tests miss that browser tests catch? What does static analysis miss that a chaos agent finds? The headline result is that <strong>Playwright catches 4 bugs that 26 passing tests and 10 approved reviewers all missed</strong>. CSP headers blocking JavaScript, multipart form parsing failures, JSON error response formatting &mdash; these are integration-level bugs that only surface when a real browser exercises the full stack.`,
			`<strong>TDD versus code-first</strong> (Experiment 48) produced the clearest evidence for test-driven development: 90.3% coverage versus 57.4% when tests are written after the code. The TDD version also had fewer logical errors because the tests constrain the implementation. Mutation testing (Experiment 49) verified test quality: an 89% mutation score means the tests detect 8 out of 9 deliberate code changes, confirming they test behavior rather than just exercising code paths.`,
			`On the security side, the <strong>chaos agent</strong> (Experiment 45) subjected a Go HTTP server to 25 attack scenarios &mdash; malformed headers, oversized bodies, concurrent connections, slow clients &mdash; and the server survived all 25. GDPR review (Experiment 50) found 10 non-compliant items in a standard CRUD application, including missing consent collection, no data export endpoint, and PII in logs. Multi-tenant review (Experiment 52) found 19 of 20 expected isolation features missing from a single-tenant application, providing a concrete SaaS readiness checklist.`,
		},
		KeyInsight: `Playwright testing is non-negotiable. Browser-level integration tests catch an entire class of bugs invisible to unit tests and code reviews.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "40-42", NumID: 40, Focus: "AI Testing + Pen Test", Result: "4 High vulns found", Cost: "$0.03", Finding: "Playwright + adversarial + pen test", HasDetail: true, SourceFile: "exp42-pentest-agent.py",
				Why:     `Our pipeline had unit tests, acceptance tests, and reviewer-driven quality gates. But none of these actually <em>attack</em> the application. Unit tests verify happy paths. Reviewers read code. Neither simulates a malicious user trying to access other people's data, inject scripts, or bypass authentication. We needed an AI agent that thinks like a penetration tester — one that discovers the attack surface, plans attacks, executes them, and chains findings.`,
				What:    template.HTML(`Three experiments combined into one testing suite. <strong>Exp 40:</strong> Playwright browser testing — AI personas use the running app through a real browser, clicking buttons and filling forms. <strong>Exp 41:</strong> Adversarial testing agent — sends malformed inputs, boundary values, and unexpected content types. <strong>Exp 42:</strong> AI pen test agent — conducts a multi-round authorized penetration test with reconnaissance, AI-planned attack vectors, execution, and finding analysis. Target: CRM at localhost:8080.

<pre class="mermaid">
flowchart TD
    APP["Running CRM App"] --> L1
    subgraph L1["Layer 1: Functional Testing"]
        PW["Exp 40: Playwright\nBrowser Journeys"]
    end
    L1 --> L2
    subgraph L2["Layer 2: Adversarial Testing"]
        ADV["Exp 41: Malformed Inputs\nBoundary Values"]
    end
    L2 --> L3
    subgraph L3["Layer 3: Pen Test Agent"]
        PT["Exp 42: AI Attack Rounds\nRecon + Exploit"]
    end
    L3 --> R{{"Vulnerabilities?"}}
    R -- "4 High" --> FOUND["Block Release"]:::failure
    R -- None --> PASS["Ship It"]:::success
    style APP fill:#f5e1e3,stroke:#d4727a
    style PW fill:#d4edda,stroke:#28a745
    style ADV fill:#fff3cd,stroke:#d4a017
    style PT fill:#fce4ec,stroke:#d4727a
    classDef success fill:#d4edda,stroke:#28a745
    classDef failure fill:#fce4ec,stroke:#d4727a
    click PW "/exp/40"
    click ADV "/exp/43"
    click PT "/exp/40"
</pre>`),
				How:     `The pen test agent (Exp 42) operates in rounds. <strong>Round 1 — Recon:</strong> crawls known endpoints (/admin, /.env, /debug, /api/internal), checks response headers for missing security controls, identifies CORS configuration. <strong>Rounds 2-5 — AI-Planned Attacks:</strong> Qwen3-30B receives the recon results and prior findings, then generates 5 targeted attacks per round (XSS payloads, IDOR attempts, SQL injection, path traversal, auth bypass). The script executes each attack via HTTP and records the response. Findings from each round inform the next round's attacks.`,
				Impact:  template.HTML("<strong>4 High severity vulnerabilities found.</strong> The pen test agent discovered: (1) missing authentication on admin endpoints (/.env returned 200), (2) IDOR on client data — changing the ID in /client/1 to /client/2 accessed another user's data with no auth check, (3) missing rate limiting — 100 requests per second with no throttling, (4) no CSRF protection on state-changing POST requests. The Playwright tests (Exp 40) separately found 4 browser-level bugs that unit tests missed: CSP blocking inline scripts, multipart form parsing failures, and JSON error response formatting. Total cost: $0.03. These three experiments became <strong>Steps 11-12 of the pipeline</strong>."),
				Related: []string{"35", "43", "45", "31"},
			},
			{Num: "43-45", NumID: 43, Focus: "Security + Chaos Agent", Result: "25/25 survived", Cost: "$0.02", Finding: "gosec, govulncheck, chaos agent", HasDetail: true, SourceFile: "exp45-chaos-agent.py",
				Why:     `The pen test agent (Exp 42) found application-level vulnerabilities — IDOR, missing auth, no CSRF. But it did not test whether the server itself survives abuse. What happens when you send a 1MB POST body? 50 concurrent connections? A request with NULL bytes in the path? Malformed HTTP headers? If the server panics or deadlocks under these conditions, every other security measure is irrelevant — the attacker just crashes the process.`,
				What:    `Three layers of security testing. <strong>Exp 43 — Static analysis:</strong> <code>gosec</code> scans for code-level vulnerability patterns, <code>govulncheck</code> scans dependencies for known CVEs. <strong>Exp 44 — Frontend/API security:</strong> checks CORS, CSP, XSS via API. <strong>Exp 45 — Chaos agent:</strong> 25 attack scenarios throwing random garbage at the running server — malformed HTTP, concurrent writes, rapid create/delete cycles, huge payloads, unicode bombs, truncated requests, invalid content types, CRLF injection, and recovery verification.`,
				How:     `The chaos agent (Exp 45) uses raw TCP sockets to send malformed HTTP requests that no library would produce — no path in GET, double Content-Length headers, negative Content-Length, 10KB headers, NULL bytes in paths. It uses threading for concurrent tests: 50 simultaneous POST requests, interleaved create-and-delete races, read-during-write contention. After all chaos tests, it verifies recovery: can the server still create a client? Still list clients? Still respond to health checks?`,
				Impact:  `<strong>Server survived all 25 chaos scenarios.</strong> Go's built-in <code>net/http</code> server handled every attack gracefully — malformed requests got 400 responses, oversized bodies were rejected, concurrent access was safe thanks to <code>sync.Mutex</code> in the store. The server processed 100 POSTs per second without errors and recovered fully after the entire chaos suite. Static analysis (gosec) found 0 critical issues. The key finding: Go's HTTP server is remarkably resilient by default — the language choice provides a strong security baseline for free. Total cost: $0.02 (only the static analysis AI calls cost money; the chaos agent is pure Python).`,
				Related: []string{"42", "40", "50", "51"},
			},
			{Num: "46-47", NumID: 46, Focus: "Screenshot + Code Metrics", Result: "Baseline set", Cost: "$0.01", Finding: "Screenshot comparison, cyclomatic complexity", HasDetail: true, SourceFile: "exp46-screenshot-to-product.py",
				Why:     `We could build products from briefs, but what about building from <em>visual references</em>? A common real-world scenario: "Build something that looks like Instatus" (a status page product). If we could describe a competitor's UI in text and have the pipeline produce a matching product, it would open a new category of briefs — visual cloning. Experiment 47 paired this with code metrics (cyclomatic complexity, line count) to establish quality baselines.`,
				What:    `<strong>Exp 46 — Screenshot to product:</strong> A detailed text description of Instatus's status page UI (layout, colors, component structure, spacing, typography) is fed to the pipeline as a "visual brief." The AI extracts a technical spec from the description, builds store + server, and produces a working status page clone. <strong>Exp 47 — Code metrics:</strong> Measures cyclomatic complexity, function length, and test coverage across all pipeline-generated apps to establish baseline quality numbers.`,
				How:     `The screenshot description is highly specific: <em>"Clean white background, fixed top nav with logo left and Subscribe button right (blue, rounded), large green banner with checkmark for All Systems Operational, component rows with name left and status badge right (green=Operational, yellow=Degraded, red=Major Outage)."</em> Qwen3-30B extracts Go types (Check, Incident, TimelineEntry) and CSS values (green=#10b981, yellow=#f59e0b). The store is built with the 2-file pattern. <code>claude -p</code> builds the server matching the visual spec.`,
				Impact:  `<strong>Baseline established.</strong> The screenshot-to-product approach worked — the resulting status page had the correct layout, color scheme, and component structure. The code metrics from Exp 47 showed typical pipeline output at 200-500 lines per app with cyclomatic complexity under 10 per function. These baselines informed later experiments: any approach that produces higher complexity than the baseline is suspect. The visual-brief technique proved that the pipeline input does not need to be a text brief — a competitor screenshot description works equally well.`,
				Related: []string{"22", "24", "54"},
			},
			{Num: "48", NumID: 48, Focus: "TDD vs Code-First", Result: "90.3% vs 57.4% coverage", Cost: "$0.02", Finding: "TDD produces 33% better coverage", HasDetail: true, SourceFile: "exp48-tdd.py",
				Why:     `Every experiment so far had generated code first, then tests — or both together (Exp 16). But TDD orthodoxy says you should write the tests first, watch them fail, then write the minimum code to make them pass. Does this actually matter when AI is writing both? We assumed not — the model generates everything in one pass anyway. But the coverage numbers from Experiment 47 were nagging: 57.4% on code-first builds. That's a lot of untested logic. What if writing tests first forces the model to think about edge cases before it writes the happy path?`,
				What:    template.HTML(`Head-to-head comparison on the same CRM store layer. <strong>Code-first</strong> (control): Qwen3-30B generates store.go, then generates store_test.go for the existing code. <strong>TDD</strong> (experiment): Qwen3-30B generates store_test.go first — tests for functions that don't exist yet — then generates store.go to make those tests pass.<br><br>Same model, same spec, same brief, same auto-fix pipeline. The only difference is the order of generation. Coverage measured with <code>go test -cover</code>.

<pre class="mermaid">
flowchart TD
    SPEC["Same CRM Spec"] --> TDD
    SPEC --> CF
    subgraph TDD["TDD Path"]
        direction TB
        T1["Write Tests First"] --> T2["Tests Reference\nNon-Existent Funcs"]
        T2 --> T3["Write Code to\nPass Tests"]
        T3 --> T4["90.3% Coverage"]:::success
    end
    subgraph CF["Code-First Path"]
        direction TB
        C1["Write Code First"] --> C2["Generate Tests\nfor Existing Code"]
        C2 --> C3["Tests Confirm\nHappy Path Only"]
        C3 --> C4["57.4% Coverage"]:::failure
    end
    style SPEC fill:#f5e1e3,stroke:#d4727a
    style T1 fill:#f5e1e3,stroke:#d4727a
    style C1 fill:#f5e1e3,stroke:#d4727a
    classDef success fill:#d4edda,stroke:#28a745
    classDef failure fill:#fce4ec,stroke:#d4727a
</pre>`),
				How:     `<strong>TDD path:</strong> Step 1 — prompt the model with the spec and ask it to write comprehensive tests. The tests reference functions and types that don't exist yet. Step 2 — prompt the model with the spec AND the test file, asking it to implement the store so all tests pass. The model can see exactly what's being tested, so it implements every code path the tests exercise.<br><br><strong>Code-first path:</strong> Step 1 — generate store.go from the spec. Step 2 — generate tests for the existing store.go. The model writes tests for what it built, which tends to test the happy path and skip edge cases it didn't think about.<br><br>Both paths use the auto-fix loop (up to 3 retries) and measure: test count, pass rate, and line coverage.`,
				Impact:  `<strong>TDD: 90.3% coverage. Code-first: 57.4% coverage.</strong> The gap is massive — 33 percentage points. The TDD version had more tests (24 vs 18), more edge case coverage (empty inputs, duplicate IDs, not-found errors), and fewer logical bugs.<br><br>The reason is subtle but important: when the model writes tests first, it <em>thinks about what should be tested</em> before it thinks about how to implement it. This surfaces edge cases the model wouldn't consider if it wrote the implementation first. When writing code-first, the model builds the happy path and then writes tests that confirm the happy path works — a self-fulfilling prophecy that misses the edges.<br><br>TDD is now the recommended approach for store layers in the pipeline. The cost is identical ($0.02 either way) — there's literally no reason not to do it.`,
				Related: []string{"16", "47", "49"},
			},
			{Num: "49", NumID: 49, Focus: "Mutation Testing", Result: "89% mutation score", Cost: "$0.01", Finding: "Tests catch 8/9 deliberate code changes", HasDetail: true, SourceFile: "exp49-mutation.py",
				Why:     `Experiment 48 showed that TDD produces 90.3% test coverage. But coverage measures which lines execute, not whether the tests actually <em>verify</em> anything. A test that calls a function but never checks the return value gets 100% coverage with 0% usefulness. Mutation testing is the gold standard for test quality: deliberately break the code and see if the tests notice. If you change <code>return client</code> to <code>return nil</code> and the tests still pass, those tests are worthless.`,
				What:    `Applied mutation testing to Exp 38's progressive CRM (32 passing tests). Defined deliberate mutations across the store layer: return-value changes (<code>return client</code> to <code>return nil</code>), boundary changes (<code>==</code> to <code>!=</code>), deletion of side effects (removing map insertions), wrong error returns, and off-by-one modifications. Each mutation is applied, all tests run, the mutation is reverted. A mutation is "caught" if the tests fail.`,
				How:     `The script copies the app to a temp directory (to avoid corrupting the original), verifies baseline tests pass, then applies each mutation one at a time. For each: read the file, replace the original string with the mutated version, run <code>go test ./... -count=1</code>, record whether tests caught it (build failure or test failure counts as caught), then restore the original. The mutation score = caught / total mutations.`,
				Impact:  `<strong>89% mutation score — tests caught 8 out of 9 deliberate code changes.</strong> The one surviving mutation was a boundary condition in list filtering that no test exercised. This confirmed that the progressive enhancement approach (Exp 38) produces tests that verify <em>behavior</em>, not just coverage. An 89% mutation score is strong — it means the test suite would catch nearly any accidental regression. The surviving mutation was added to the test backlog. Cost: $0.01 (just the initial test generation; mutation testing itself is free). Mutation testing validated our test quality without relying on coverage percentages.`,
				Related: []string{"48", "38", "26"},
			},
			{Num: "50", NumID: 50, Focus: "GDPR/Privacy Review", Result: "10 non-compliant items", Cost: "$0.01", Finding: "Privacy reviewer finds PII gaps", HasDetail: true, SourceFile: "exp50-gdpr-privacy.py",
				Why:     `The OWASP security reviewer (Exp 31) checks for XSS, CSRF, and injection. But GDPR compliance is a different category entirely — it is about how you handle personal data, not whether your code has vulnerabilities. A CRM stores names, emails, phone numbers, and addresses. In Europe, shipping without GDPR compliance is not just bad practice — it is illegal. None of our existing reviewers checked for privacy compliance.`,
				What:    `Three privacy-focused review layers against Exp 37's CRM. <strong>Layer 1 — GDPR Reviewer:</strong> examines the spec for data processing compliance (consent collection, lawful basis, data minimisation, retention policy, right to access/deletion/portability). <strong>Layer 2 — PII Code Reviewer:</strong> examines the actual Go code for data leaks (PII in logs, unencrypted storage, missing anonymisation, data in error messages). <strong>Layer 3 — Privacy Pen Tester:</strong> attempts to access or leak PII through the running application.`,
				How:     `The GDPR reviewer receives the main.go and store.go code and checks against a structured checklist: does the app collect consent before storing PII? Is there a data export endpoint (right to portability)? Can a user delete their data (right to erasure)? Is there a retention policy? Is PII logged? The PII code reviewer scans for <code>log.Printf</code> calls containing user data, unencrypted email storage, and data exposure in error responses. The pen tester hits endpoints trying to extract bulk PII.`,
				Impact:  `<strong>10 non-compliant items found.</strong> The biggest gaps: no consent collection before storing client data, no data export endpoint (GDPR Article 20 violation), PII visible in application logs, no retention policy (data stored forever), no right-to-deletion mechanism (clients can be deleted from the UI but not via a formal data subject request), email addresses stored in plain text with no encryption at rest. Every standard CRUD app we had built had the same gaps — GDPR was simply never part of the pipeline. This led to adding a privacy review step in the security testing phase. Cost: $0.01. The 10 findings became a reusable GDPR checklist for all future pipeline builds.`,
				Related: []string{"42", "43", "52"},
			},
			{Num: "51", NumID: 51, Focus: "Concurrent Users", Result: "6/6 pass", Cost: "FREE", Finding: "Threading, burst, consistency tests", HasDetail: true, SourceFile: "exp51-concurrent-users.py",
				Why:     `Every test so far simulated a single user. But real products have multiple users hitting endpoints simultaneously. What happens when User A adds a client while User B deletes one? When 10 users create clients at the same time? When one user reads the client list while another modifies it? Race conditions are invisible in single-user testing and catastrophic in production — corrupted data, panics, deadlocks.`,
				What:    `Six concurrent user scenarios against the running CRM: <strong>(1)</strong> Two users creating clients simultaneously (10 creates each, interleaved). <strong>(2)</strong> One user reading while another writes (20 operations mixed). <strong>(3)</strong> Create-and-immediately-delete race (5 threads). <strong>(4)</strong> Burst traffic — 10 concurrent creates hitting the same endpoint. <strong>(5)</strong> Consistency check — count of listed clients matches number of successful creates. <strong>(6)</strong> Final verification — server healthy after all concurrent abuse.`,
				How:     `Pure Python threading — no AI calls needed. Each scenario spawns multiple threads that hit HTTP endpoints simultaneously using <code>urllib</code>. Thread A POSTs new clients while Thread B GETs the client list. The script tracks response status codes, checks for 500 errors (server panics), and verifies data consistency — if 20 clients were successfully created, does GET /api/clients return exactly 20? All threads use <code>join(timeout=15)</code> to prevent hangs.`,
				Impact:  `<strong>6/6 scenarios passed.</strong> No 500 errors, no panics, no deadlocks, no data corruption. The <code>sync.Mutex</code> in the store layer (generated by the pipeline in Exp 37) handled all concurrent access correctly. Listed client count matched created client count exactly. The server responded to health checks after all concurrent abuse. Cost: FREE — no AI calls, pure Python HTTP threading. This confirmed that Go's concurrency model plus our standard in-memory store pattern is production-safe for moderate concurrent load. The test suite takes 3 seconds to run and could be added to any pipeline build at zero cost.`,
				Related: []string{"45", "43", "52"},
			},
			{Num: "52", NumID: 52, Focus: "Multi-Tenant Review", Result: "19/20 features missing", Cost: "$0.01", Finding: "SaaS readiness audit checklist", HasDetail: true, SourceFile: "exp52-multi-tenant.py",
				Why:     `Every CRM we had built was single-tenant — one user, one dataset, no concept of organisations or teams. But SaaS products serve multiple companies simultaneously. If Company A's client list leaks to Company B, that is a business-ending data breach. We needed to understand the gap: how far is a pipeline-generated single-tenant app from being multi-tenant ready? What exactly needs to change?`,
				What:    `A SaaS architect AI reviewer examines Exp 37's CRM code (store.go and main.go) against a 20-point multi-tenant readiness checklist: <strong>Data isolation</strong> (org_id on every model, tenant-scoped queries), <strong>Auth</strong> (per-route authentication, role-based access, session management), <strong>Billing</strong> (subscription tracking, usage limits, payment integration), <strong>Team features</strong> (invitations, roles, org settings), <strong>Rate limiting</strong> (per-tenant quotas).`,
				How:     `The reviewer receives the actual Go source code and checks each requirement against what exists. For data isolation: is there a tenant_id field on Client, Activity, Invoice? Does every database query filter by tenant? For auth: does every HTTP handler check authentication? Is there middleware that extracts the tenant from the JWT? For billing: are there subscription checks before creating resources? The reviewer rates each item as PRESENT, PARTIAL, or MISSING.`,
				Impact:  `<strong>19 of 20 multi-tenant features missing.</strong> The only feature partially present was basic CRUD operations. Missing: tenant_id on all models, authentication on any route, role-based access control, team invitations, organisation settings, subscription management, usage limits, per-tenant rate limiting, data export per tenant, audit logging, tenant-scoped search, and 7 more. Total cost: $0.01. This was not a failure of the pipeline — our brief never asked for multi-tenancy. The value is the <strong>concrete SaaS readiness checklist</strong>: 20 specific items that transform a single-tenant CRUD into a multi-tenant SaaS. This checklist now informs the pipeline when the brief specifies "SaaS" or "multi-tenant."`,
				Related: []string{"50", "51", "39"},
			},
		},
	},
	{
		Slug:     "website-cloning",
		Name:     "Website Cloning",
		Icon:     "\U0001f310",
		ExpRange: "Experiments 54, 95, 98, 118 – 119, 128, 130, 133, 152 – 173 (incl. 172b)",
		Narrative: []template.HTML{
			`This research stream explored how well AI can reproduce real websites &mdash; and evolved into a full design system analysis. It started with a 9-model comparison cloning eachandother.com, expanded to an Airbnb listing with AI-generated photos, and culminated in cloning <strong>10 diverse websites</strong> to extract universal design patterns.`,
			`The <strong>9-model comparison</strong> (Exp 54) revealed that SSIM pixel similarity is misleading for cross-model comparison &mdash; Qwen3-30B scored highest but looked worst. Opus was unanimously the best visual match. <strong>Llama 4 Scout</strong> delivered 80% of the quality at 1% of the cost ($0.02 vs $2.75).`,
			`<strong>AI image generation</strong> (Exp 95) proved that Nano Banana 2 produces photorealistic apartment photos 5x faster and 70x cheaper than GPT-5 Image. The full pipeline: AI builds the layout, AI generates the images &mdash; a complete product page with zero real photography.`,
			`The <strong>image-guided cloning</strong> breakthrough (Exp 98) overturned the text-only approach entirely. By feeding the actual screenshot to Opus, a single iteration produced better layout fidelity than 11 iterations of text-only refinement. The Nando&rsquo;s screenshot-guided clone reproduced "Ultimate winter warmer" because it could <em>read</em> the text; the text-only clone invented "GET YOUR PERI-PERI FIX TODAY" after 11 attempts. Hybrid mode (screenshot + design tokens) scored highest of all.`,
			`The <strong>10-site design system analysis</strong> (Exp 118-119) cloned Airbnb, eachandother, Stripe, Tailwind, Medium, Linear, Nando&rsquo;s, BBC News, Product Hunt, and GitHub. Universal patterns emerged: spacing scale (4-40px in 7+/10 sites), border radius (4-12px), one primary accent per site, dark sections in 5/10. Component analysis found 10 MUST-have components in 8+/10 sites. This data-driven approach builds the framework from evidence, not assumptions.`,
			`With screenshot-guided cloning proven, we asked: can we do it <em>cheaper</em>? The Nando&rsquo;s clone took 18 iterations at $0.60 each &mdash; $10.80 total. <strong>Layer-by-layer construction</strong> (Exp 128) split the work into 6 focused layers: colors, grid, typography, shapes, images, interactivity. Each layer has one job, and layer 4 (decorations) produced the best CSS clip-paths and wavy dividers we&rsquo;d seen because the model wasn&rsquo;t competing with text or images. Six calls, ~$2.50, best visual CSS. <strong>Component assembly</strong> (Exp 130) went further &mdash; a single Opus call with component structure (navComponent, heroComponent, etc.) produced the best SSIM scores (Layout 0.569, Overall 0.257) at just $0.60. One shot, best metrics. The optimal strategy may be component assembly + 2&ndash;3 focused layer refinements for ~$1.80 total.`,
			`<strong>CSS + DOM extraction</strong> (Exp 133) attacked the guessing problem from the other side. Instead of estimating CSS values from screenshots, we extracted the <em>exact</em> computed styles from the live site using Playwright&rsquo;s <code>getComputedStyle()</code> on every major element &mdash; padding, margins, colors, font sizes, the lot. Combined with the screenshot for layout reference, this gives the model both visual and numerical precision. A key discovery: Nando&rsquo;s serves different content based on carousel state and possibly A/B testing, which explains some variation across our cloning experiments.`,
			`<strong>Bounding box extraction</strong> (Exp 152&ndash;154) took the data-driven approach to its logical conclusion. Instead of giving the model a screenshot and saying "figure it out," we extracted every visible element&rsquo;s exact bounding box &mdash; position, size, colors, text content &mdash; from the live page. Exp 152 used the full 196-block dataset on Nando&rsquo;s and scored 0.575 overall with 0.949 on color accuracy. The key discovery: <strong>truncating the block data (80 vs 196 blocks) causes the model to hallucinate the bottom half of the page</strong>. Full data is critical. Iterative refinement (Exp 154) pushed the score to 0.603 by cropping the 5 worst sections and sending original+clone pairs for targeted fixes, producing a clone that users described as "structurally amazing."`,
			`<strong>Figma.com clone</strong> (Exp 156) proved the bounding box approach generalises beyond Nando&rsquo;s. With 293 blocks extracted from Figma&rsquo;s clean, minimal landing page, the clone scored <strong>0.700 overall &mdash; the best score in the entire research programme</strong>. Footer scored 0.739, Design Systems section 0.732. The hero was weakest (0.422) because Figma&rsquo;s dark interactive canvas becomes a white placeholder in the clone. Clean, component-driven sites clone better than visually complex ones.`,
			`<strong>AI logo generation</strong> (Exp 158) and <strong>carousel detection</strong> (Exp 159) solved two remaining gaps. Logo generation via Nano Banana 2 produces recognizable company icons (Spotify&rsquo;s green circle, GitHub&rsquo;s octocat, Notion&rsquo;s N-in-square) at ~$0.07/logo. Carousel detection uses ARIA attributes and tab roles to automatically find and capture all carousel states &mdash; Figma&rsquo;s 8-tab product carousel was fully extracted, with 123 setInterval timers paused to freeze animations. These two pieces close the gap between "static snapshot" and "full interactive site" cloning.`,
			`<strong>Linear.app clone</strong> (Exp 161) extended the bounding box approach to a 4th site type: dark/minimal SaaS tool. 221 blocks extracted, 99 key blocks used. The clone scored Layout 0.603 (good structural match) but only 0.328 overall &mdash; because <strong>SSIM color histogram penalises dark-on-dark sites</strong>. Subtle dark grey variations (the kind that make a dark theme feel polished) register as near-identical to the metric. Layout score is the meaningful measure for dark sites. The approach now covers restaurants (Nando&rsquo;s), design tools (Figma), payments (Stripe), and project management (Linear) &mdash; suggesting it generalises across site types, though dark themes need a different scoring approach.`,
			`<strong>Design transfer</strong> (Exp 163) took the pipeline from "clone" to "create." Instead of reproducing a site as-is, it took Figma&rsquo;s visual design (fixed header, card grid, tabbed features, organic-shape CTA, dark footer) and filled it with SaaS Tools product content (10 embeddable tools, pricing, FAQ). Teal replaced purple, 11 AI images were generated, and the result scored 0.749 vs the SaaS Tools original (content match) and 0.586 vs the Figma original (design retention). This proves <strong>automated design transfer</strong>: clone a design you admire, swap in your own content, generate matching images &mdash; a complete product page in minutes instead of weeks.`,
		},
		KeyInsight: `Screenshot-guided cloning beats 11 text-only iterations in layout quality. Bounding box extraction beats screenshots &mdash; Figma scored 0.700 overall, the best yet. Full block data is critical: truncated datasets hallucinate the page bottom. Iterative refinement on worst sections pushes scores +4.7%. Clean/minimal sites clone better than complex ones. Dark-themed sites (Linear, 0.328 overall) expose an SSIM blind spot &mdash; layout score (0.603) is more meaningful than color histogram for dark-on-dark. Design transfer (Exp 163) proves the pipeline can take one site&rsquo;s design language and fill it with another product&rsquo;s content &mdash; the step from cloning to creating. The pipeline is converging on: extract everything &rarr; build once &rarr; refine worst sections &rarr; transfer designs.`,
		TableType:  "website-clone",
		Experiments: []Experiment{
			{Num: "1", NumID: 0, Focus: "Opus", Result: "8/8", Cost: "$2.75", Finding: "Best overall", Category: "0.708"},
			{Num: "2", NumID: 0, Focus: "Sonnet", Result: "8/8", Cost: "$2.51", Finding: "Very good", Category: "0.669"},
			{Num: "3", NumID: 0, Focus: "Llama 4 Scout", Result: "8/8", Cost: "~$0.02", Finding: "Good — best value", Category: "0.722"},
			{Num: "4", NumID: 0, Focus: "Devstral Small", Result: "8/8", Cost: "~$0.02", Finding: "Wireframe style", Category: "0.732"},
			{Num: "5", NumID: 0, Focus: "Haiku", Result: "8/8", Cost: "$1.59", Finding: "Good", Category: "0.649"},
			{Num: "6", NumID: 0, Focus: "Qwen3-30B", Result: "5/8", Cost: "~$0.01", Finding: "Broken CSS", Category: "0.745"},
			{Num: "7", NumID: 0, Focus: "Qwen3-32B", Result: "6/8", Cost: "~$0.01", Finding: "Incomplete", Category: "—"},
			{Num: "8", NumID: 0, Focus: "DeepSeek V3.2", Result: "4/8", Cost: "~$0.02", Finding: "Incomplete", Category: "—"},
			{Num: "9", NumID: 0, Focus: "Gemma 3 27B", Result: "6/8", Cost: "~$0.01", Finding: "Wrong architecture", Category: "—"},
			{Num: "95", NumID: 95, Focus: "AI Image Generation for Clones", Result: "9/9 images, 2 models", Cost: "~$0.20/img (GPT-5), ~$0.003/img (Nano Banana)", Finding: "Nano Banana 2 produces photorealistic apartment photos 5x faster and 70x cheaper", HasDetail: true, Icon: "\U0001f3a8\U0001f5bc\ufe0f",
				Why:     `The website clones from Exp 54 scored well on layout but the photo gallery was all grey placeholder boxes. Real Airbnb listings have stunning interior photography. Without images, even the best-laid-out clone looks like a wireframe. Could AI generate apartment photos that are <em>different</em> from the originals — legally clean, not copies — but match the style and quality of real listing photography?`,
				What:    `We generated 9 images for the Whitechapel apartment clone: hero living room, master bedroom, twin bedroom, kitchen, bathroom, host avatar, neighbourhood street scene, and two reviewer avatars. Two models compared:<br><br><strong>GPT-5 Image</strong> (OpenAI via OpenRouter): ~$0.20 per image, ~80s generation time, ~1.6MB output<br><strong>Nano Banana 2</strong> (Gemini 3.1 Flash Image Preview): ~$0.003 per image, ~17s generation time, ~0.9MB output`,
				How:     `Each model received the same 9 prompts describing London apartment interiors, Whitechapel streetscapes, and portrait headshots. Prompts specified: photorealistic style, real estate photography, specific furniture and lighting details, London architectural features (sash windows, brick buildings). Images returned as base64 via the OpenRouter API <code>images</code> field. Both versions served from the same Go clone — just swap the <code>IMAGE_DIR</code> environment variable.`,
				Impact:  `<strong>Nano Banana 2 won decisively: 5x faster, 70x cheaper, comparable quality.</strong> The Whitechapel street scene with market stalls and red buses is indistinguishable from a real photograph. The apartment interiors have correct London architectural details.<br><br>Interestingly, SSIM dropped slightly with AI images (0.81 vs 0.83 with grey placeholders) because the original listing has different photos. But as a product, the AI-image version looks dramatically more professional. This proves the full pipeline: <strong>AI generates the layout, AI generates the images, the result is a complete, realistic product page with zero real photography.</strong>`,
				Related: []string{"54"},
			},
			{Num: "118", NumID: 118, Focus: "Design System Token Extraction", Result: "10/10 sites analyzed", Cost: "~$5 total (10 Opus clones)", Finding: "Universal spacing: 4,8,12,16,20,24,32. Every site: one primary accent + sans-serif", HasDetail: true, SourceFile: "",
				Why: template.HTML(`We had 2 clones &mdash; Airbnb and eachandother. Enough to prove AI can replicate a website, not enough to build a <em>design system</em>. If you want a generic framework that works for any website, you need data from <em>diverse</em> websites. Do SaaS landing pages, news sites, restaurants, and developer tools actually share anything? Or is every site a snowflake?<br><br>We could have guessed. "8px grid, probably." But guessing is what designers do when they don't have evidence. We had 10 AI-cloned websites sitting right there. Time to extract every CSS token and find out what's actually universal.`),
				What: template.HTML(`We cloned 10 websites with Opus, each representing a different industry vertical: <strong>Airbnb</strong> (e-commerce/marketplace), <strong>eachandother</strong> (design agency), <strong>Stripe</strong> (SaaS), <strong>Tailwind CSS</strong> (developer docs), <strong>Medium</strong> (blog/publishing), <strong>Linear</strong> (dashboard/app), <strong>Nando's</strong> (restaurant), <strong>BBC News</strong> (media), <strong>Product Hunt</strong> (marketplace/community), and <strong>GitHub</strong> (developer tool). Each clone was a single Go HTTP server with embedded HTML/CSS/SVG &mdash; the same pattern from <a href="/exp/54">Experiment 54</a>. From each clone we extracted every CSS token: color values, font families, spacing values (margin, padding, gap), border-radius, box-shadow, and font sizes.`),
				How: template.HTML(`Token extraction via regex on the CSS embedded in each clone. We counted occurrences of every <code>color:</code>, <code>background-color:</code>, and <code>border-color:</code> value. Extracted every <code>margin</code>, <code>padding</code>, and <code>gap</code> value in pixels. Collected all <code>border-radius</code> and <code>box-shadow</code> declarations. Mapped <code>font-family</code> stacks.<br><br>Then we compared all 10 token sets. A value appearing in <strong>7+/10 sites</strong> was classified as <em>universal</em>. <strong>5+/10</strong> was <em>common</em>. <strong>3+/10</strong> was <em>frequent</em>. Anything below 3 was site-specific noise. This gave us an empirically-derived design token scale, not a theoretical one.`),
				Impact: template.HTML(`The universal spacing scale practically wrote itself: <strong>4, 6, 8, 10, 12, 14, 16, 20, 24, 32, 40px</strong> &mdash; all appearing in 7+ of 10 sites. Border radius converged on <strong>4, 8, 10px</strong> (5+/10). Every single site had exactly <strong>one primary accent color</strong> &mdash; Stripe's purple, Linear's blue, Nando's red, BBC's maroon. All 10 used <strong>sans-serif body fonts</strong> (Inter, system-ui, Helvetica). And 5 out of 10 had dark-themed sections, which means dark mode support isn't optional for a modern framework.<br><br>This is real evidence, not design theory. When we build the CSS framework, we're not guessing at the spacing scale &mdash; we extracted it from 10 production websites across 10 different industries.`),
				Related: []string{"54", "95"},
			},
			{Num: "119", NumID: 119, Focus: "Component Pattern Analysis", Result: "10 MUST, 2 SHOULD, 2 Optional", Cost: "\u2014", Finding: "Nav, hero, footer, CTA, icons, badges, dividers appear in 8+/10 sites", HasDetail: true, SourceFile: "",
				Why: template.HTML(`Knowing the tokens &mdash; colors, spacing, border-radius &mdash; tells you <em>how</em> to style components. It doesn't tell you <em>which</em> components to build. Is a "hero section" universal or just a SaaS marketing thing? Does every website need a footer? Are card grids everywhere or just on landing pages?<br><br>We had 10 clones representing 10 industries. Instead of guessing which components belong in a generic framework, we could simply count. If a component appears in 7+ of 10 sites, it's a must-have. If it appears in fewer than 3, it's too specific. The data decides.`),
				What: template.HTML(`We analyzed the HTML/CSS of all 10 clones for 15 component patterns: <strong>nav</strong>, <strong>hero section</strong>, <strong>card grid</strong>, <strong>footer</strong>, <strong>sidebar</strong>, <strong>CTA button</strong>, <strong>avatar</strong>, <strong>SVG icons</strong>, <strong>dark sections</strong>, <strong>image gallery</strong>, <strong>tabs</strong>, <strong>list items</strong>, <strong>badges</strong>, <strong>dividers</strong>, and <strong>search</strong>. For each component, we checked all 10 sites and recorded presence or absence.`),
				How: template.HTML(`Regex scan of each clone's Go source for component markers: <code>&lt;nav</code> elements, <code>hero</code> CSS classes, <code>&lt;footer</code> tags, <code>&lt;svg</code> elements, <code>grid</code>/<code>card</code> class patterns, <code>badge</code> classes, <code>&lt;hr</code> or divider patterns, <code>search</code> inputs. We classified each component by frequency:<br><br><strong>MUST</strong> (7+/10 sites): include in the framework, no questions asked.<br><strong>SHOULD</strong> (5+/10 sites): include as an optional module.<br><strong>Optional</strong> (3+/10 sites): document the pattern but don't ship by default.<br><strong>Skip</strong> (&lt;3/10 sites): too niche for a generic framework.<br><br>We called it the <em>2/3 rule</em>: if a component doesn't serve at least 7 out of 10 sites, it's too specific for the generic layer.`),
				Impact: template.HTML(`<strong>10 MUST components</strong>: nav (10/10), hero (9/10), footer (9/10), CTA buttons (9/10), SVG icons (10/10), dark sections (10/10), list items (8/10), badges (10/10), dividers (10/10), search (8/10). <strong>2 SHOULD</strong>: card grid (6/10), avatar (5/10). <strong>2 Optional</strong>: sidebar (3/10), tabs (3/10). Image gallery was Airbnb-only &mdash; skip.<br><br>This gives us an exact build list for the CSS framework. We're not building "all the components" or copying Bootstrap's kitchen-sink approach. We're building the 10 components that appear on 80%+ of real websites, plus 2 optional modules for the next tier. Everything else is the user's problem. That's how you keep a framework small and fast.`),
				Related: []string{"118", "54"},
			},
			{Num: "98", NumID: 98, Focus: "Image-Guided Cloning", Result: "Layout +5% in 1 iteration vs 11 text-only", Cost: "~$0.50 per clone", Finding: "Screenshot-guided produces better layout match in 1 iteration than 11 text-only iterations", HasDetail: true, Icon: "\U0001f4f8", SourceFile: "",
				Why: template.HTML(`The clones from Experiments 54 and 118 were all built from text descriptions &mdash; the AI reads a prompt saying "clone this website" and <em>imagines</em> what the site looks like. The result? Opus invents content. The Nando's text-only clone produced "GET YOUR PERI-PERI FIX TODAY" as its hero heading after 11 iterations of refinement. The real site says "Ultimate winter warmer." No amount of text-only iteration can fix the fundamental problem: the AI is guessing at visual layout because it has never <em>seen</em> the original. What if we just showed it?`),
				What: template.HTML(`Three approaches compared head-to-head on the same websites:<br><br><strong>A &mdash; Text-only</strong>: the existing approach. Describe the site in words, let Opus imagine the layout. Nando's was run for 11 iterations of refinement. BBC News was run for 1 iteration.<br><strong>B &mdash; Screenshot-guided</strong>: feed the actual screenshot of the original site to Opus via its Read tool. One iteration, no refinement. The model sees what it needs to reproduce.<br><strong>C &mdash; Hybrid (screenshot + tokens)</strong>: feed both the screenshot and the extracted design tokens from <a href="/exp/118">Experiment 118</a>. One iteration.<br><br><strong>Nando's results</strong>: A (text, 11 iter) layout 0.547, overall 0.308 | B (screenshot, 1 iter) layout 0.564, overall 0.254 | C (hybrid, 1 iter) layout 0.573, overall 0.267<br><strong>BBC News results</strong>: A (text, 1 iter) layout 0.645, overall 0.638 | B (screenshot, 1 iter) layout 0.653, overall 0.627`),
				How: template.HTML(`Opus reads the actual screenshot file via its multimodal Read tool and generates a single Go HTTP server with embedded HTML/CSS to match what it sees. The screenshot approach uses the same infrastructure as all previous clone experiments &mdash; single binary, no dependencies, embedded assets &mdash; but with a fundamentally different input: pixels instead of prose.<br><br>SSIM scores were computed at two levels: <strong>layout SSIM</strong> (structural comparison with blurred/downscaled images, measuring spatial arrangement) and <strong>overall SSIM</strong> (pixel-level comparison). The overall scores are lower for screenshot-guided clones because the original sites contain real photographs while clones use placeholder images. The layout score is what matters &mdash; it measures whether the structural arrangement of elements matches.`),
				Impact: template.HTML(`Screenshot guidance is now the default for the clone pipeline. The results are unambiguous: <strong>1 screenshot-guided iteration beats 11 text-only iterations on layout fidelity.</strong> The Nando's B clone reproduced "Ultimate winter warmer" as its heading because it could <em>read</em> the text from the screenshot. The text-only clone invented completely different copy after 11 attempts.<br><br>The hybrid approach (C) scored highest on layout (0.573) by combining visual guidance with extracted design tokens, giving the model both the picture and the precise CSS values. Overall SSIM is lower across the board for screenshot-guided clones, but this is expected &mdash; real sites have unique photography that clones replace with placeholders. The structural layout is what the pipeline cares about, and screenshot guidance wins decisively.`),
				Related: []string{"54", "118", "119"},
			},
			{Num: "128", NumID: 128, Focus: "Layer-by-Layer Construction", Result: "6 layers, ~$2.50, best visual CSS", Cost: "~$2.50", Finding: "Separating concerns (1 layer per call) produces better CSS shapes and clip-paths", HasDetail: true, Icon: "\U0001f9f1",
				Why:    template.HTML(`18 iterations at $0.60 each = $10.80 for a Nando&rsquo;s clone. The screenshot-guided approach from <a href="/exp/98">Experiment 98</a> proved that showing the model the target works better than describing it, but we were still throwing everything at the model at once &mdash; layout, typography, colors, decorations, images, interactivity &mdash; and hoping it would juggle all of them simultaneously. What if we separated concerns? Give each call ONE job, let it focus, and build up the page layer by layer. Could we get the same quality cheaper?`),
				What:   template.HTML(`We built the Nando&rsquo;s clone in 6 focused layers, each building on the output of the previous one:<br><br><strong>Layer 1 &mdash; Colors only</strong>: background colors, gradients, brand palette<br><strong>Layer 2 &mdash; Grid layout</strong>: structural arrangement of sections, flexbox/grid<br><strong>Layer 3 &mdash; Typography</strong>: font families, sizes, weights, line heights<br><strong>Layer 4 &mdash; Shapes &amp; decorations</strong>: clip-paths, wavy dividers, SVG decorations, border-radius<br><strong>Layer 5 &mdash; Images</strong>: placeholder images, background images, aspect ratios<br><strong>Layer 6 &mdash; Interactivity</strong>: hover states, transitions, animations<br><br>Each layer received the screenshot plus the accumulated HTML/CSS from previous layers.`),
				How:    template.HTML(`Each layer is a single Opus call. The model receives the target screenshot and the current state of the clone (output from all previous layers). Its job is to add <em>only</em> the elements for its layer. Layer 4 (shapes and decorations) is the key differentiator &mdash; when the model focuses solely on CSS decorative elements, it produces proper <code>clip-path</code> polygons and wavy <code>&lt;svg&gt;</code> dividers instead of the rough approximations we got from all-at-once generation.`),
				Impact: template.HTML(`Layer-by-layer produced the <strong>best visual CSS quality</strong> of any approach we tested. The clip-paths and wavy dividers in Layer 4 were done properly because the model was focusing solely on decorations &mdash; no competing priorities from text, layout, or images. 6 calls vs 18 iterations. <strong>~$2.50 vs ~$10.80</strong> &mdash; a 4x cost reduction. Layout score was comparable to the iterative approach. The trade-off: slightly more orchestration complexity, but dramatically better CSS craftsmanship for decorative elements.`),
				Related: []string{"98", "130", "133"},
			},
			{Num: "130", NumID: 130, Focus: "Component Assembly", Result: "1 shot, $0.60, best metric scores", Cost: "~$0.60", Finding: "Single Opus call with component structure produces best SSIM scores at lowest cost", HasDetail: true,
				Why:    template.HTML(`Layer-by-layer (<a href="/exp/128">Experiment 128</a>) cut costs from $10.80 to $2.50 while improving CSS quality. But 6 calls is still 6 calls. Could we skip iteration entirely and build a good clone in <em>one shot</em>? The hypothesis: if we give Opus a clear component structure up front &mdash; navComponent, heroComponent, menuComponent &mdash; it can organise its output better in a single call than when asked to generate a monolithic page.`),
				What:   template.HTML(`A single Opus call builds the Nando&rsquo;s clone as separate component functions &mdash; <code>navComponent()</code>, <code>heroComponent()</code>, <code>menuGridComponent()</code>, <code>footerComponent()</code> &mdash; each with scoped CSS, assembled into one page. The prompt describes the entire site with component boundaries marked. Opus generates all components plus the assembly logic in one call.`),
				How:    template.HTML(`One prompt, one call. The prompt includes the target screenshot and a structural breakdown of the page into named components. Opus generates a Go HTTP server with each component as a function returning its HTML fragment and associated CSS. The main handler assembles them in order. This mirrors how real frontend frameworks work &mdash; isolated components composed into a page.`),
				Impact: template.HTML(`<strong>Best SSIM scores at the lowest cost.</strong> Layout SSIM 0.569, Overall SSIM 0.257 &mdash; beating both the 11-iteration text-only approach and the screenshot-guided single iteration from <a href="/exp/98">Experiment 98</a>. Total cost: <strong>$0.60, 1 API call</strong>. The catch: component assembly doesn&rsquo;t capture CSS decorative details (clip-paths, wavy dividers) as well as <a href="/exp/128">layer-by-layer</a>. The optimal approach may be: <strong>component assembly (1 shot) + 2&ndash;3 focused layer refinements</strong> for ~$1.80 total &mdash; best metrics plus best visual CSS, at 83% less than the original iterative approach.`),
				Related: []string{"98", "128", "133"},
			},
			{Num: "133", NumID: 133, Focus: "CSS + DOM Extraction", Result: "Real CSS data improves accuracy", Cost: "~$0.60", Finding: "Extracting actual computed CSS and DOM structure from the live site provides exact values instead of guessing from screenshots", HasDetail: true,
				Why:    template.HTML(`Screenshot-guided cloning tells the model <em>what</em> the site looks like. But the model still has to <em>guess</em> the CSS values &mdash; is that padding 16px or 20px? Is the font-size 14px or 15px? Is that color <code>#1a1a2e</code> or <code>#1c1c30</code>? What if we extracted the EXACT computed styles from the real site and gave the model both the picture AND the precise numbers?`),
				What:   template.HTML(`Playwright visits the live Nando&rsquo;s site and extracts everything a model needs to reproduce it precisely:<br><br><strong>Computed CSS</strong>: backgroundColor, padding, margin, fontSize, fontWeight, fontFamily, color, lineHeight, borderRadius, display, flexDirection, gap &mdash; for every major element via <code>getComputedStyle()</code><br><strong>DOM tree</strong>: tag hierarchy with nesting depth, classes, and roles<br><strong>Color palette</strong>: every unique color used on the page<br><strong>Image URLs</strong>: src attributes for all images<br><strong>Link structure</strong>: navigation and footer links<br><strong>Animations</strong>: CSS transitions and keyframe definitions<br><strong>Media queries</strong>: responsive breakpoints`),
				How:    template.HTML(`A Playwright script launches a headless browser, navigates to the target site, and runs <code>page.evaluate()</code> to call <code>getComputedStyle()</code> on every element matching major CSS selectors (headings, paragraphs, links, buttons, images, nav, footer, sections). The results are structured as JSON: element tag, classes, and a map of computed property values. The DOM tree is walked recursively to capture nesting. All data is fed to Opus alongside the screenshot &mdash; visual reference for layout, numerical precision for CSS values.`),
				Impact: template.HTML(`This provides the <strong>ground truth CSS values</strong> that screenshot-only approaches have to estimate. Combined with the screenshot for overall layout reference, the model gets both visual context and numerical precision &mdash; the best of both worlds. A key discovery during extraction: <strong>Nando&rsquo;s serves different content based on carousel state and possibly A/B testing</strong>, which explains the variation we observed across cloning experiments (Exp 137&ndash;139). The extracted DOM and CSS also revealed structural patterns invisible in screenshots &mdash; hidden navigation drawers, lazy-loaded sections, and JavaScript-dependent content that only appears after interaction.`),
				Related: []string{"98", "128", "130"},
			},
			{Num: "152", NumID: 152, Focus: "Full Bounding Box Clone (Nando's)", Result: "Overall 0.575, Color 0.949", Cost: "~$0.60", Finding: "Full 196-block dataset eliminates hallucinated bottom half — footer score jumps from 0.031 to 0.445 (+1335%)", HasDetail: true, Icon: "\U0001f4e6\U0001f532",
				CloneShot: "/screenshots/exp152-nandos/desktop.png", RefShot: "/screenshots/exp152-ref/desktop.png",
				Why:    template.HTML(`<a href="/exp/133">Experiment 133</a> extracted CSS and DOM from the live site but still relied on the model to <em>interpret</em> those values and arrange them spatially. What if we went further &mdash; extracting not just styles but the exact <strong>bounding box</strong> of every visible element? Position, size, background color, text content, font metrics, all of it. Give the model a complete spatial map of the page and ask it to reproduce exactly what it sees in the data. No guessing, no imagination &mdash; pure reconstruction from coordinates.<br><br>An earlier attempt (Exp 150) used only 80 blocks from the dataset. The model built the top half of the page beautifully, then <em>hallucinated</em> the bottom half because it had no data for those sections. The bottom sections (app store, gift card, footer) were invented from context. This experiment used the FULL 196-block dataset from <a href="/exp/133">Experiment 133</a> to see if complete data fixes the hallucination problem.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Live Site"] --> B["Extract 196 Blocks"]
    B --> C["Block Data: position, size, color, text"]
    C --> D["Claude: Reconstruct Page"]
    D --> E["1032-line Clone"]
    E --> F["SSIM Score: 0.575"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style F fill:#d4edda,stroke:#28a745
</pre>
We extracted 196 bounding boxes from the live Nando&rsquo;s page &mdash; every visible element with its exact pixel position, dimensions, background color (including the navy <code>rgb(16,32,79)</code> app section and teal <code>rgb(8,217,167)</code> gift card section), text content, and font properties. The full dataset was fed to Claude in a single call. Output: a 1032-line single-file Go HTTP server serving the clone on port 9310.`),
				How:    template.HTML(`The bounding box extraction script uses Playwright to capture every visible element on the page: <code>getBoundingClientRect()</code> for position and size, <code>getComputedStyle()</code> for colors and fonts, and <code>textContent</code> for copy. Elements are filtered by visibility (no hidden elements, no zero-size boxes) and sorted top-to-bottom. The full 196-block JSON is sent to Claude alongside the page screenshot.<br><br>The critical difference from Exp 150: we used ALL 196 blocks instead of truncating to 80. With the full dataset, every section of the page has corresponding block data &mdash; including the app download section (navy background), gift card section (teal background), and the 6-column footer with real link text.`),
				Impact: template.HTML(`<strong>Full data eliminates hallucination.</strong> The bottom half of the page &mdash; app section, gift card, footer &mdash; is now structurally correct because the model had real block data to work from. Footer SSIM jumped from 0.031 (Exp 150, truncated data) to <strong>0.445 &mdash; a +1335% improvement</strong>. Overall SSIM hit 0.575 with color accuracy at an extraordinary 0.949. The model reproduced exact brand colors because they were in the data, not guessed from a screenshot.<br><br>The lesson is simple but profound: <strong>truncated input = hallucinated output</strong>. When the model runs out of real data, it invents plausible-looking content that scores terribly. Give it the complete picture and it builds the complete page. 1 Claude call, 1032 lines, done.`),
				Related: []string{"133", "130", "98"},
			},
			{Num: "153", NumID: 153, Focus: "Section-by-Section Clone (Nando's)", Result: "Overall 0.497", Cost: "~$4.20 (7 calls)", Finding: "Section-by-section builds all sections correctly but scores lower without per-section screenshots", HasDetail: true,
				Why:    template.HTML(`<a href="/exp/152">Experiment 152</a> proved that full bounding box data produces accurate clones. But it sends all 196 blocks to a single model call, which means the model has to hold the entire page in context while generating 1032 lines of HTML/CSS. What if we split the page into logical sections and built each one independently? Each call would have a smaller, focused context &mdash; just the blocks for its section &mdash; and could focus on getting that one piece right.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart TD
    A["196 Blocks"] --> B["Split by Section"]
    B --> C1["Nav: 12 blocks"]
    B --> C2["Hero: 28 blocks"]
    B --> C3["Food Grid: 35 blocks"]
    B --> C4["Rewards: 22 blocks"]
    B --> C5["App Section: 18 blocks"]
    B --> C6["Two-Column: 15 blocks"]
    B --> C7["Footer: 66 blocks"]
    C1 --> D["7 Claude Calls"]
    C2 --> D
    C3 --> D
    C4 --> D
    C5 --> D
    C6 --> D
    C7 --> D
    D --> E["Assemble: 1651 lines"]
    E --> F["SSIM: 0.497"]
    style A fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style F fill:#fce4ec,stroke:#d4727a
</pre>
The page was split into 7 logical sections: navigation, hero, food grid, rewards, app download, two-column layout, and footer. Blocks were assigned to sections based on their vertical position. Each section was built by an independent Claude call receiving only its blocks. All 7 outputs were assembled into a single page &mdash; 1651 lines total, served on port 9311.`),
				How:    template.HTML(`The 196 blocks were partitioned by Y-coordinate ranges corresponding to visible section boundaries on the page. Each of the 7 calls received: the section name, its blocks (with positions relative to the section top), and the full-page screenshot for visual context. Each call generated a standalone HTML/CSS fragment. A final assembly step wrapped all fragments in a single page with shared reset CSS and scoped section styles.<br><br>The approach produced structurally complete output &mdash; all 7 sections were present and correct. But the overall SSIM score (0.497) was <em>lower</em> than the single-pass approach (0.575). The reason: each section call received the <strong>full-page screenshot</strong> rather than a cropped screenshot of just that section. The model could see the whole page but was told to build only one section, creating an ambiguity about exactly which part to focus on.`),
				Impact: template.HTML(`Section-by-section has clear <strong>structural advantages</strong>: every section was present, the page was complete, and each piece was individually manageable. But it scored lower because of a <strong>tooling limitation</strong>: the CLI couldn&rsquo;t pass cropped per-section screenshots. The model received the full screenshot but was asked to focus on one section &mdash; a contradiction that hurt precision.<br><br>The potential is obvious: with per-section cropped images, each call would have both the block data AND a focused visual reference for its exact section. This would likely beat single-pass. The approach is documented for the next iteration where the extraction pipeline can produce section-level crops alongside block data.`),
				Related: []string{"152", "128"},
			},
			{Num: "154", NumID: 154, Focus: "Iterative Refinement (Nando's)", Result: "Overall 0.603 (+4.7%)", Cost: "~$3.60 (6 calls)", Finding: "Targeting the 5 worst sections with original+clone crop pairs pushes score from 0.575 to 0.603", HasDetail: true, Icon: "\U0001f504\U0001f4d0",
				Why:    template.HTML(`<a href="/exp/152">Experiment 152</a> scored 0.575 &mdash; good, but not great. Some sections scored well (color: 0.949) while others were weak. Rather than regenerating the entire page, could we identify the <em>worst</em> sections and fix just those? This is the approach human designers use: build the page, then iterate on the parts that don&rsquo;t look right. The question was whether AI could do the same &mdash; look at a side-by-side comparison and make targeted improvements.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Exp 152 Clone<br>0.575"] --> B["SSIM per Section"]
    B --> C["Find 5 Worst"]
    C --> D["Crop Original + Clone"]
    D --> E["Claude: Fix Each Section"]
    E --> F["+ AI Images"]
    F --> G["Refined Clone<br>0.603"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style G fill:#d4edda,stroke:#28a745
</pre>
We took the Exp 152 clone, computed per-section SSIM scores, identified the 5 worst-scoring sections, and sent Claude pairs of crops &mdash; original section vs clone section &mdash; for each weak area. Claude generated targeted CSS/HTML fixes for each section. Additionally, AI-generated images were added (hero.png, wing-platter.png, restaurant-exterior.png, etc.) to replace placeholder boxes. Served on port 9312.`),
				How:    template.HTML(`The refinement pipeline: (1) render both original and clone at the same viewport, (2) divide each into horizontal bands corresponding to page sections, (3) compute SSIM for each band, (4) select the 5 lowest-scoring bands, (5) send each pair (original crop + clone crop) to Claude with the instruction "make the clone match the original." Claude responds with specific HTML/CSS changes for each section.<br><br>The key learning: <strong>fix whole logical blocks, not arbitrary pixel bands.</strong> Early attempts divided the page into equal-height strips, which would cut through the middle of a card grid or split a heading from its body text. Aligning the refinement bands to actual HTML section boundaries produced dramatically better results because Claude could reason about complete components.`),
				Impact: template.HTML(`Overall SSIM improved from 0.575 to <strong>0.603 &mdash; a +4.7% gain</strong>. Footer jumped to 0.771. User feedback was emphatic: <em>"structurally amazing"</em> &mdash; the structure and positioning were nearly perfect, with the remaining differences being content-specific (real photos vs AI-generated, exact text variations).<br><br>The refinement approach validates an important principle: <strong>build once with full data, then refine the weak spots</strong>. This is far more efficient than iterating the entire page. 1 build call + 5 refinement calls = 6 total calls at ~$3.60, producing a better result than 18 full-page iterations at $10.80. The pipeline is converging on: extract everything &rarr; build once &rarr; identify worst sections &rarr; targeted refinement.`),
				Related: []string{"152", "128", "130"},
			},
			{Num: "155", NumID: 155, Focus: "Logical Section Refinement", Result: "Planned (approach documented)", Cost: "—", Finding: "Identify sections from HTML structure, not pixel bands — fix entire logical blocks together", HasDetail: true,
				Why:    template.HTML(`<a href="/exp/154">Experiment 154</a> refined sections by dividing the page into horizontal pixel bands and scoring each. This worked, but the bands were arbitrary &mdash; they didn&rsquo;t align with the page&rsquo;s actual structure. A band might cut through the middle of a 3-card grid, sending one card to one refinement call and two cards to another. The model can&rsquo;t fix half a card grid. The next step was obvious: identify sections from the <strong>HTML structure</strong> itself, so refinement targets are complete logical blocks &mdash; the entire card grid, the complete footer, the full hero section.`),
				What:   template.HTML(`The planned approach: parse the clone&rsquo;s HTML to identify top-level sections (typically <code>&lt;section&gt;</code>, <code>&lt;nav&gt;</code>, <code>&lt;footer&gt;</code>, or major <code>&lt;div&gt;</code> containers). Map each section to a Y-coordinate range in the rendered output. Crop the original and clone at these structural boundaries. Score each section. Refine the weakest by sending the complete structural section &mdash; all 3 cards together, the entire footer with all columns, the complete hero with all its children.`),
				How:    template.HTML(`Not completed in this session. The approach is documented as the next evolution of the refinement pipeline. Implementation would require: (1) HTML parser to identify section boundaries, (2) mapping from DOM position to rendered Y-coordinate, (3) structural-aware cropping, (4) per-section refinement calls with complete logical context.`),
				Impact: template.HTML(`This experiment represents the design direction for the next iteration. The insight from Exp 154 was clear: <strong>logical blocks produce better refinement than arbitrary bands</strong>. When the refinement pipeline aligns with the page&rsquo;s actual component structure, each call gets complete context and can reason about the full component. Expected improvement: +2&ndash;5% over arbitrary band refinement, with better visual coherence in multi-element sections like card grids and multi-column footers.`),
				Related: []string{"154", "152", "153"},
			},
			{Num: "156", NumID: 156, Focus: "Clone Figma.com", Result: "Overall 0.700 — BEST score", Cost: "~$0.60", Finding: "Clean/minimal sites clone dramatically better than complex ones — Figma scored 0.700 vs Nando's 0.575", HasDetail: true, Icon: "\U0001f310\U0001f3a8",
				CloneShot: "/screenshots/exp156-figma/desktop.png", RefShot: "/screenshots/exp156-ref/desktop.png",
				Why:    template.HTML(`Every bounding box experiment so far used Nando&rsquo;s &mdash; a visually complex restaurant site with decorative elements, food photography, and busy layouts. Would the approach work on a <em>different</em> kind of site? Figma&rsquo;s landing page is the opposite of Nando&rsquo;s: clean, minimal, component-driven, with large whitespace and simple geometric shapes. If the bounding box approach generalises, it should work on any site. If it only works on one type, that&rsquo;s a problem.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart TD
    A["figma.com"] --> B["Dismiss Cookie Banner"]
    B --> C["Extract 293 Blocks"]
    C --> D["Claude: Build Clone"]
    D --> E["516 lines + refinement"]
    E --> F["SSIM: 0.700"]
    F --> G["BEST Score in Research"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style F fill:#d4edda,stroke:#28a745
    style G fill:#d4edda,stroke:#28a745
</pre>
We extracted 293 bounding boxes from figma.com &mdash; first dismissing the cookie consent banner (which had blocked previous capture attempts). The initial clone was 516 lines, then refined to include text logos and the tab group component. Served on ports 9314/9315.<br><br><strong>Section scores</strong>: Footer 0.739, Design Systems 0.732, Ship Products 0.708, hero weakest at 0.422 (dark interactive canvas vs white placeholder).`),
				How:    template.HTML(`The capture script was updated to handle cookie banners: click the dismiss button, wait for the banner to animate away, then begin extraction. This fixed a recurring issue where the cookie overlay was captured as part of the page content.<br><br>293 blocks were extracted &mdash; significantly more than Nando&rsquo;s 196, reflecting Figma&rsquo;s component-rich layout with many small elements (tab labels, feature cards, logo thumbnails). The model received the full block dataset plus screenshot in a single call. Refinement focused on adding the text-based company logos in the logo bar and properly rendering the 8-tab product navigator.`),
				Impact: template.HTML(`<strong>0.700 overall SSIM &mdash; the highest score in the entire research programme.</strong> This is 22% better than the Nando&rsquo;s bounding box clone (0.575) and dramatically better than any previous approach. The difference comes from site complexity: Figma&rsquo;s clean, grid-based layout with large blocks of solid color and generous whitespace plays to the model&rsquo;s strengths. It can reproduce geometric layouts with high fidelity.<br><br>The hero section scored lowest (0.422) because Figma&rsquo;s hero is a dark interactive canvas with animated elements &mdash; the clone renders a white placeholder. This is an inherent limitation: interactive/animated content can&rsquo;t be reproduced from a static screenshot. Every other section scored above 0.700.<br><br><strong>Key finding: site complexity is the primary predictor of clone quality.</strong> Clean, component-driven sites like Figma clone dramatically better than visually complex sites like Nando&rsquo;s. This suggests the pipeline should prioritise bounding box extraction for SaaS landing pages, developer tools, and portfolio sites over restaurant and e-commerce sites with heavy photography.`),
				Related: []string{"152", "154", "98"},
			},
			{Num: "157", NumID: 157, Focus: "Clone Stripe.com", Result: "Extraction complete, clone pending", Cost: "—", Finding: "Bounding box extraction completed for Stripe — clone generation deferred to next session", HasDetail: true, Icon: "\U0001f4b3",
				CloneShot: "/screenshots/exp157-stripe/desktop.png", RefShot: "/screenshots/exp157-ref/desktop.png",
				Why:    template.HTML(`After Figma&rsquo;s record-breaking 0.700, we wanted to test another SaaS landing page to see if the pattern holds. Stripe is the canonical "beautiful SaaS landing page" &mdash; if the bounding box approach can reproduce Stripe&rsquo;s signature gradient backgrounds, card layouts, and code snippets, it would confirm that clean design sites consistently score well.`),
				What:   template.HTML(`Bounding box extraction was completed on stripe.com &mdash; reference screenshot captured, block data extracted. The clone generation step was not completed in this session due to time constraints.`),
				How:    template.HTML(`The same Playwright-based extraction pipeline used for Figma and Nando&rsquo;s. Cookie banner handling was applied (Stripe uses a similar consent overlay). Block data captured and stored for the next session.`),
				Impact: template.HTML(`Extraction validates that the pipeline is now <strong>generalised and repeatable</strong>: point it at any URL, dismiss the cookie banner, extract blocks, generate clone. The Stripe clone will test whether the approach handles gradient backgrounds and code-block components as well as it handles Figma&rsquo;s clean geometry. Expected in the next session.`),
				Related: []string{"156", "152"},
			},
			{Num: "158", NumID: 158, Focus: "AI Logo Generation", Result: "9+ logos, ~$0.07/logo", Cost: "~$0.61 (batch of 9)", Finding: "Nano Banana 2 generates recognizable company logos with correct icons when prompted specifically", HasDetail: true, Icon: "\U0001f5bc\ufe0f",
				Why:    template.HTML(`The clones had a visual gap: company logos in navigation bars and footer sections were replaced with text or generic placeholders. A Figma clone without the Figma logo, a Stripe section without the Stripe wordmark &mdash; the eye notices immediately. We needed AI-generated logos that are <em>recognizable</em> as the real brands without copying the trademarked assets pixel-for-pixel. Could the same model that generates apartment photos (<a href="/exp/95">Experiment 95</a>) produce company logos?`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Logo Prompt"] --> B["Nano Banana 2"]
    B --> C{"Version?"}
    C -->|V1: Wordmark Only| D["Clean text, missing icons"]
    C -->|V2: Logo + Icon| E["Spotify circle+waves<br>GitHub octocat<br>Notion N-in-square"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style D fill:#fce4ec,stroke:#d4727a
    style E fill:#d4edda,stroke:#28a745
</pre>
Two rounds of generation using Nano Banana 2 (google/gemini-3.1-flash-image-preview) via OpenRouter:<br><br><strong>V1 &mdash; Wordmark-only</strong>: prompted for text logos only. Clean typography but missing the iconic symbols &mdash; Spotify without the green circle, GitHub without the octocat. $0.61 for 9 logos.<br><strong>V2 &mdash; Logo with icon</strong>: prompted specifically for "logo icon + wordmark." Produced Spotify&rsquo;s green circle with sound waves, GitHub&rsquo;s octocat silhouette, Notion&rsquo;s N-in-a-square. Recognizable at a glance.<br><br>Additionally generated 18 site images for the Figma clone: hero canvas, community collaboration cards, template screenshots, and UI mockups.`),
				How:    template.HTML(`Each prompt described the specific brand&rsquo;s visual identity: "Spotify logo: bright green circle (#1DB954) with three curved sound wave lines inside, plus the wordmark &lsquo;Spotify&rsquo; in white on a dark background." The key insight was that <strong>generic prompts produce generic results</strong> &mdash; you have to describe the icon specifically. "GitHub logo" gets a text wordmark; "GitHub logo: white octocat silhouette (cat with tentacle tail) inside a dark circle, plus &lsquo;GitHub&rsquo; wordmark" gets the recognizable icon.<br><br>Cost was higher than expected: ~$0.07/logo vs the ~$0.003/image from <a href="/exp/95">Experiment 95</a>. The difference is complexity &mdash; logos require precise shapes and text rendering, while apartment photos are probabilistic (any plausible bedroom works). Logo generation tolerates less variation.`),
				Impact: template.HTML(`AI logo generation works when you <strong>describe the icon specifically</strong>. The V2 logos are recognizable as their brands at thumbnail size, which is exactly what clones need in nav bars and footers. At ~$0.07/logo, generating all company logos for a clone page costs under $1 &mdash; a small addition to the ~$0.60 clone generation cost.<br><br>The 18 Figma site images are more significant: hero canvas mockups, collaboration cards, template screenshots. These transform the Figma clone from "correct layout with grey boxes" to "looks like the real product page." Combined with the 0.700 SSIM score, the Figma clone with AI images is the most convincing clone the research has produced.`),
				Related: []string{"95", "156"},
			},
			{Num: "159", NumID: 159, Focus: "Carousel Detection", Result: "2 carousels found, 8 tab states captured", Cost: "—", Finding: "ARIA attributes and tab roles reliably detect carousels — all slides extractable before cloning", HasDetail: true, Icon: "\U0001f3a0",
				Why:    template.HTML(`Static screenshots capture one moment in time. But modern websites are full of <strong>carousels, sliders, and tabbed content</strong> that show different information depending on which slide is active. The Nando&rsquo;s experiments suffered from this: the live site serves different hero images and promotional content based on carousel state. A screenshot captures one state; the clone needs to represent the full content. Could we detect carousels automatically and capture every state before cloning?`),
				What:   template.HTML(`<pre class="mermaid">
flowchart TD
    A["figma.com"] --> B["Detect Carousels"]
    B --> C["Logo Bar Carousel"]
    B --> D["Product Tab Carousel"]
    D --> E["8 Tabs Detected"]
    E --> F1["Prompt"]
    E --> F2["Design"]
    E --> F3["Draw"]
    E --> F4["Build"]
    E --> F5["Publish"]
    E --> F6["Promote"]
    E --> F7["Jam"]
    E --> F8["Present"]
    F1 --> G["All Panel Content Captured"]
    style A fill:#f5e1e3,stroke:#d4727a
    style B fill:#fff3cd,stroke:#d4a017
    style E fill:#fff3cd,stroke:#d4a017
    style G fill:#d4edda,stroke:#28a745
</pre>
Tested on figma.com, the detection script found 2 carousels automatically:<br><br><strong>Carousel 1 &mdash; Logo bar</strong>: company logos rotating through a horizontal strip. Detected via <code>aria-roledescription="carousel"</code>.<br><strong>Carousel 2 &mdash; Product tabs</strong>: 8 tabs (Prompt, Design, Draw, Build, Publish, Promote, Jam, Present) each with their own panel content. Detected via <code>role="tablist"</code> and <code>[aria-label*="slide"]</code> selectors.`),
				How:    template.HTML(`The detection script runs in Playwright and uses three strategies:<br><br><strong>1. ARIA attributes</strong>: <code>aria-roledescription="carousel"</code>, <code>role="tablist"</code>, <code>role="tabpanel"</code><br><strong>2. Semantic patterns</strong>: <code>[aria-label*="slide"]</code>, <code>[aria-label*="carousel"]</code>, elements with <code>overflow:hidden</code> containing wider children<br><strong>3. Timer freezing</strong>: paused <strong>123 setInterval timers</strong> and froze all CSS animations to prevent auto-rotation during capture<br><br>For the product tab carousel, the script clicked each of the 8 tabs in sequence, waited for the panel transition, and captured the panel content (text, images, layout) for each state. Each tab revealed different marketing content &mdash; "Prompt" shows AI features, "Design" shows the canvas, "Build" shows dev tools, etc.`),
				Impact: template.HTML(`Carousel detection solves a fundamental problem with static screenshot cloning: <strong>one screenshot captures one state, but the page contains multiple states.</strong> By detecting and extracting all carousel/tab states before cloning, we give the model complete information about the page&rsquo;s full content.<br><br>The ARIA-based detection is robust because modern websites use these attributes for accessibility compliance &mdash; they&rsquo;re present on well-built sites precisely because screen readers need them. The timer freezing (123 intervals paused on Figma alone) prevents carousels from auto-advancing during extraction, ensuring we capture clean states.<br><br>Next step: feed all carousel states into the bounding box extraction pipeline so clones can reproduce tabbed/carousel content instead of capturing whatever happened to be visible at screenshot time.`),
				Related: []string{"156", "133", "152"},
			},
			{Num: "160", NumID: 160, Focus: "Deep Carousel Capture", Result: "4 carousels found, 12 logos + 8 tabs + 10 templates extracted", Cost: "—", Finding: "CSS infinite scroll carousels need heuristic detection — ARIA only catches 2 of 4 carousel types", HasDetail: true,
				Why:    template.HTML(`<a href="/exp/159">Experiment 159</a> found 2 carousels on Figma using ARIA attributes. But the user noticed that the logo bar was also a carousel &mdash; and it wasn&rsquo;t detected. This deeper investigation aimed to find ALL carousel-like content on the page, including CSS-only infinite scrolls that don&rsquo;t use ARIA markup.`),
				What:   template.HTML(`Figma&rsquo;s homepage actually has <strong>4 distinct carousels</strong>, not 2:<br><br><strong>1. Logo bar</strong> &mdash; CSS infinite scroll with 12 company logos (AirBnb, Atlassian, Dropbox, Duolingo, GitHub, Microsoft, Netflix, NYT, Pentagram, Slack, Stripe, Zoom). Built from 4 duplicated &lt;ul&gt; tracks, each 1952px wide with 12 items. All logos are SVGs from Sanity CDN.<br><strong>2. Hero carousel</strong> &mdash; 8 iframe slides with interactive Figma Sites demos. Cross-origin iframes prevent content extraction.<br><strong>3. Product tabs</strong> &mdash; 8 tabs (Prompt, Design, Draw, Build, Publish, Promote, Jam, Present) with full panel content extracted.<br><strong>4. Template gallery</strong> &mdash; 10 category cards (Websites, Social media, Mobile apps, Presentations, Invitations, Illustrations, Portfolio, Plugins, Web ads, Icons).`),
				How:    template.HTML(`Three-strategy detection approach:<br><br><strong>1. ARIA carousels</strong>: <code>aria-roledescription="carousel"</code> catches hero + template gallery.<br><strong>2. Tab panels</strong>: <code>role="tablist"</code> click-through captures all 8 product tabs.<br><strong>3. CSS infinite scroll</strong> (new): heuristic detection finds containers with 4+ visual children in horizontal alignment, then deduplicates by checking for repeated <code>img[alt]</code> values across duplicated tracks. This is the only method that catches the logo bar.<br><br>The logo bar uses 48 DOM elements for just 12 unique logos &mdash; each duplicated 4 times to create a seamless CSS scroll animation.`),
				Impact: template.HTML(`The finding that <strong>ARIA only catches 2 of 4 carousel types</strong> is significant for the cloning pipeline. CSS-only infinite scrolls (logo bars, testimonial tickers) need heuristic detection. The 3-strategy approach is now the recommended method:<br><br>1. ARIA attributes for accessible carousels<br>2. Tab role click-through for tabbed interfaces<br>3. Overflow + duplication heuristic for CSS scroll animations<br><br>Limitation: iframe-based slides (hero carousel) contain cross-origin content that cannot be extracted. Only <code>aria-label</code> metadata is available.`),
				Related: []string{"159", "156"},
			},
			{Num: "161", NumID: 161, Focus: "Clone Linear.app", Result: "Overall 0.328, Layout 0.603", Cost: "~$0.14 (4 AI images)", Finding: "Dark-on-dark sites score poorly on color histogram — layout score (0.603) is more meaningful for dark themes", HasDetail: true, Icon: "\U0001f319",
				CloneShot: "/screenshots/exp161-linear/desktop.png", RefShot: "/screenshots/exp161-ref/desktop.png",
				Why:    template.HTML(`The bounding box approach had been validated on Nando&rsquo;s (restaurant, colourful), Figma (design tool, clean/minimal), and Stripe (payments, gradient-heavy). But all three are light-themed sites. What happens with a <strong>dark-themed SaaS tool</strong>? Linear.app &mdash; black background, white text, purple accents &mdash; is the archetypal dark minimal product site. If the approach generalises to dark themes, it works across all major site categories. Tesla was attempted first but returned "Access Denied" &mdash; it blocks headless browsers entirely.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["linear.app"] --> B["Extract 221 Blocks"]
    B --> C["Trim to 99 Key Blocks"]
    C --> D["Claude: Build Clone"]
    D --> E["4 AI Images (Nano Banana 2)"]
    E --> F["Go Server on :9317"]
    F --> G["Layout 0.603 | Overall 0.328"]
    style A fill:#1a1a2e,stroke:#5E6AD2,color:#fff
    style C fill:#fff3cd,stroke:#d4a017
    style G fill:#d4edda,stroke:#28a745
</pre>
221 bounding boxes extracted from Linear&rsquo;s homepage via Playwright, trimmed to 99 key elements. The clone reproduces all major sections: header with nav, hero with tagline, product preview (mock issue tracker UI), secondary tagline, 5 feature sections (Intake, Plan, Build, Diffs, Monitor) each with mock UI panels, changelog grid, testimonial quote, CTA with green glow effect, and a 5-column footer. 4 AI images generated with Nano Banana 2: hero-ui, feature-agents, feature-cycles, changelog. Served on port 9317.`),
				How:    template.HTML(`The same Playwright-based extraction pipeline used for previous sites. 221 blocks captured with position, size, color, and text content. Trimmed to 99 key blocks to focus on structural elements rather than duplicated or decorative noise. The Go server was built directly from bounding box data plus reference screenshot analysis in a single Claude call.<br><br>Feature sections include mock UI elements: issue tracker rows with status indicators, Kanban board cards, a code diff viewer, and an analytics dashboard with a bar chart. The <code>/images/</code> handler serves 4 AI-generated images from Nano Banana 2 at $0.035 each ($0.14 total).<br><br>A critical instruction was applied: <strong>DO NOT shrink fonts during refinement to match height ratio.</strong> Previous experiments showed that models reduce font sizes to compress the page vertically, destroying readability. This "no shrink" rule preserves the original typography scale.`),
				Impact: template.HTML(`Layout SSIM 0.603 confirms the bounding box approach produces good structural matches even on dark-themed sites. But overall SSIM of 0.328 exposes an <strong>SSIM blind spot for dark themes</strong>: the color histogram metric can&rsquo;t distinguish subtle dark grey variations (e.g. <code>#0a0a0a</code> vs <code>#111111</code> vs <code>#1a1a1a</code>) that give dark UIs their depth. The color score was just 0.069. On light sites these differences register as distinct colours; on dark sites they collapse to "all black."<br><br>The approach now covers <strong>4 site types</strong>: restaurant (Nando&rsquo;s), design tool (Figma), payments (Stripe), and project management (Linear). The consistent Layout 0.6+ scores across all four suggest the bounding box pipeline generalises well. The scoring system, however, needs evolution &mdash; structural similarity or perceptual metrics like LPIPS may be more appropriate than pixel-level SSIM for dark-themed sites.`),
				Related: []string{"156", "152", "154"},
			},
			{Num: "162", NumID: 162, Focus: "Clone to Product (saastools.gyrum.io)", Result: "Overall 0.789 — HIGHEST SCORE", Cost: "~$0.34 (5 AI images)", Finding: "Cloning your own site scores highest — simple, consistent designs clone best. First experiment proving pipeline can rebuild a real product.", HasDetail: true, Icon: "\U0001f3ed", SourceFile: "",
				CloneShot: "/screenshots/exp162-saastools/desktop.png", RefShot: "/screenshots/exp162-ref/desktop.png",
				Why:    template.HTML(`Every previous cloning experiment targeted third-party sites &mdash; Nando&rsquo;s, Stripe, Figma, Linear. That proves the pipeline can <em>replicate</em> someone else&rsquo;s design. But the real commercial question is different: <strong>can the clone pipeline rebuild a REAL product you own?</strong> If yes, the pipeline becomes viable for rapid prototyping, marketing variant generation, post-migration rebuilds, and white-labelling. saastools.gyrum.io &mdash; Gyrum&rsquo;s own SaaS tools marketplace &mdash; was the ideal test target: 10 tools, pricing tiers, FAQ accordion, and a design we control end-to-end.`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["saastools.gyrum.io"] --> B["Extract 178 Blocks"]
    B --> C["Claude: Build Clone"]
    C --> D["5 AI Images ($0.34)"]
    D --> E["821 Lines of Go"]
    E --> F["Server on :9318"]
    F --> G["Overall 0.789 — HIGHEST"]
    style A fill:#0d9488,stroke:#0f766e,color:#fff
    style C fill:#fff3cd,stroke:#d4a017
    style G fill:#28a745,stroke:#1e7e34,color:#fff
</pre>
178 bounding boxes extracted from saastools.gyrum.io. The clone reproduces every section of the original: navigation bar, hero with code snippet, 10-tool grid, 3-step &ldquo;how it works&rdquo; section, pricing table (Free / Growth / Bundle tiers), trust badges, FAQ accordion, CTA, and footer. 5 AI images generated ($0.34 total): hero dashboard mockup, tool icons, and a code snippet visual. The Go server is 821 lines, served on port 9318.<br><br><strong>Per-section scores:</strong> Nav/Hero 0.903, Steps 0.890, FAQ 0.887, Pricing 0.844. The highest-scoring sections are the structurally simplest &mdash; nav bars and step grids have clear, predictable layouts that the pipeline reproduces almost perfectly.`),
				How:    template.HTML(`The same Playwright-based bounding box extraction pipeline used across all cloning experiments. 178 blocks captured with position, size, colour, and text content. The Go server was built from bounding box data plus reference screenshot analysis in a single Claude call, producing 821 lines of Go.<br><br>5 AI images were generated with Nano Banana 2 at ~$0.07 each ($0.34 total): a hero dashboard mockup, tool category icons, and a code snippet illustration. All images are served via the <code>/images/</code> handler.<br><br><strong>Scoring breakdown:</strong><br>&bull; SSIM: 0.843 &mdash; pixel-level structural similarity, highest in the programme<br>&bull; Color: 0.981 &mdash; near-perfect colour histogram match (we control the palette)<br>&bull; Layout: 0.570 &mdash; bounding box alignment<br>&bull; Height ratio: 1.01 &mdash; almost exact vertical match<br>&bull; <strong>Overall: 0.789</strong> &mdash; new programme record<br><br>The previous highest was Nando&rsquo;s at ~0.72. The 0.789 score comes from two factors: (1) the site uses a simple, consistent design system with predictable section patterns, and (2) we control the design, so there are no surprises &mdash; no third-party widgets, no A/B test variations, no anti-bot measures.`),
				Impact: template.HTML(`This is the most significant result in the cloning programme for two reasons:<br><br><strong>1. Highest score proves simpler sites clone better.</strong> The 0.789 overall score &mdash; beating every third-party clone &mdash; confirms that design consistency is the primary predictor of clone quality. Sites with uniform section patterns, consistent spacing, and controlled colour palettes produce higher scores than visually complex sites with animations, gradients, and dynamic content.<br><br><strong>2. First &ldquo;clone to product&rdquo; proof.</strong> This is the first experiment proving the pipeline can go from &ldquo;existing product&rdquo; to &ldquo;rebuilt clone.&rdquo; The commercial implications are significant: a company could use this pipeline to:<br>&bull; <strong>Rapidly prototype redesigns</strong> &mdash; clone the current site, then iterate on the clone<br>&bull; <strong>Generate marketing variants</strong> &mdash; clone and modify for campaigns, A/B tests<br>&bull; <strong>Rebuild after migrations</strong> &mdash; capture the old site, rebuild on new infrastructure<br>&bull; <strong>White-label products</strong> &mdash; clone and rebrand for different customers<br><br>The pipeline is no longer just a research tool. At 0.789 overall, the output is production-viable for landing pages and marketing sites.`),
				Related: []string{"161", "156", "152", "154"},
			},
			{Num: "163", NumID: 163, Focus: "Design Transfer (Figma design + SaaS Tools content)", Result: "SSIM vs SaaS Tools: 0.749 | vs Figma: 0.586", Cost: "$0.74 (11 AI images)", Finding: "AI can perform 'design transfer' — take one site's design language and apply it to a completely different product's content. Foundation for automated theme/template systems.", HasDetail: true, Icon: "\U0001f3a8\u27a1\ufe0f\U0001f310", SourceFile: "",
				Why:    template.HTML(`Can AI take a design it cloned and repurpose it for a different product? This is the step from "clone" to "create" &mdash; using existing designs as templates for new products. Every previous cloning experiment reproduced a site as-is. But the commercial value isn&rsquo;t in exact copies; it&rsquo;s in taking a design you admire and filling it with <em>your own</em> content. If design transfer works, the pipeline goes from "research experiment" to "business tool."`),
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Figma Clone<br>Design Language"] --> C["Design Transfer"]
    B["SaaS Tools<br>Product Content"] --> C
    C --> D["New Product Page<br>723 lines + 11 AI images"]
    D --> E["SSIM vs SaaS Tools: 0.749<br>SSIM vs Figma: 0.586"]
    style A fill:#4D49FC,stroke:#3a37c9,color:#fff
    style B fill:#0d9488,stroke:#0f766e,color:#fff
    style D fill:#fff3cd,stroke:#d4a017
    style E fill:#d4edda,stroke:#28a745
</pre>
Took Figma&rsquo;s clean layout &mdash; fixed header, card grid, tabbed features, CTA with organic shapes, dark footer &mdash; and filled it with SaaS Tools product content: 10 embeddable tools with real gyrum.io links, fuller descriptions, 3-step setup, pricing cards, FAQ, and trust badges. Teal accent (<code>#0d9488</code>) replaces Figma&rsquo;s purple (<code>#4D49FC</code>). 723 lines of Go, 11 AI images ($0.74): 10 tool icons + hero product dashboard. Served on port 9319.`),
				How:    template.HTML(`Started from the Figma clone&rsquo;s CSS/HTML structure (<a href="/exp/156">Experiment 156</a>), which provided the visual foundation: fixed navigation bar, hero section with gradient overlays, card grid layout, tabbed feature showcase, CTA section with organic background shapes, and dark footer. All content was replaced with SaaS Tools data from <a href="/exp/162">Experiment 162</a>: 10 tools (JSON Formatter, Regex Tester, CSS Generator, etc.) with real gyrum.io links, 3-step "How It Works" setup flow, pricing cards (Free/Growth/Bundle tiers), FAQ accordion, and trust badges.<br><br>The colour system was transformed: Figma&rsquo;s purple primary (<code>#4D49FC</code>) swapped to teal (<code>#0d9488</code>) to match the SaaS Tools brand. 11 AI images were generated with Nano Banana 2: a hero product dashboard mockup and 10 individual tool icons representing each embeddable tool. All 10 tools link to their real gyrum.io endpoints.`),
				Impact: template.HTML(`Two SSIM scores tell the story:<br><br><strong>0.749 vs SaaS Tools original</strong> &mdash; strong content match, proving the product information transferred accurately into the new design shell. Sections, tools, pricing, and FAQ all present and correct.<br><strong>0.586 vs Figma original</strong> &mdash; moderate design retention, confirming the layout structure (header, card grid, tabs, CTA, footer) carried over even with completely different content and colour scheme.<br><br>This proves a practical commercial use case: <strong>clone a design you admire, swap in your own content, generate matching images.</strong> A complete product page in minutes, not weeks. The pipeline goes from "research experiment" to "business tool." The implications for automated theme/template systems are significant &mdash; a library of cloned designs becomes a template library, each instantly fillable with any product&rsquo;s content.`),
				Related: []string{"156", "162", "161"},
			},
			{Num: "164", NumID: 164, Focus: "Design Transfer with Tool Previews", Result: "11 AI preview images, tool cards show embedded widget mockups", Cost: "~$0.74", Finding: "AI-generated tool preview screenshots transform card grid from icons to product mockups — each card shows what the tool looks like embedded in a real website", HasDetail: true,
				Why:    `Exp 163 design transfer used SVG icons for tool cards. Real product pages show screenshots of the tool in action. Could AI generate convincing widget mockup screenshots?`,
				What:   `Generated 11 images: hero dashboard + 10 tool previews (feedback popup, help centre, cookie banner, changelog widget, waitlist form, status page, onboarding tooltip, testimonial wall, announcement bar, social proof toast)`,
				How:    `Nano Banana 2 prompts described each tool embedded in a website context. Images used as card backgrounds with CSS object-fit:cover.`,
				Impact: `Tool cards transform from "icon + text" to "here's what this looks like in your app" &mdash; much more convincing for a product page. Total image cost $0.74 for 11 images.`,
				Related: []string{"163", "158", "162"},
			},
			{Num: "165", NumID: 165, Focus: "Optimized Pipeline (Vercel.com)", Result: "SSIM 0.787, Overall 0.622", Cost: "~$0.60", Finding: "Combined ALL best methods into a reference pipeline: full blocks + cookie dismiss + real UA + logical sections + AI images. 893 lines, 12+ AI images, 3 iterations.", HasDetail: true, Icon: "\u2699\ufe0f\U0001f680",
				CloneShot: "/screenshots/exp165-vercel/desktop.png", RefShot: "/screenshots/exp165-ref/desktop.png",
				Why:    `Every previous experiment tested individual techniques in isolation: full bounding boxes (<a href="/exp/152">Exp 152</a>), iterative refinement (<a href="/exp/155">Exp 155</a>), AI images (<a href="/exp/158">Exp 158</a>), design transfer (<a href="/exp/163">Exp 163</a>). But we had never combined ALL of them into a single optimized pipeline run. What happens when you stack every proven method together on a new target site?`,
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Vercel.com"] --> B["Full Block Extraction"]
    B --> C["Cookie Dismiss + Real UA"]
    C --> D["Logical Section Grouping"]
    D --> E["Claude: Build Clone"]
    E --> F["AI Image Generation"]
    F --> G["3 Iteration Rounds"]
    G --> H["SSIM 0.787"]
    style A fill:#f5e1e3,stroke:#d4727a
    style D fill:#fff3cd,stroke:#d4a017
    style H fill:#d4edda,stroke:#28a745
</pre>
The optimized pipeline combined every validated technique: full bounding box extraction (not truncated), cookie/consent banner dismissal, real browser user-agent strings (avoiding bot detection), logical section grouping for context, AI-generated images matching the target aesthetic, and 3 rounds of iterative refinement targeting the weakest sections. Target: Vercel.com, a developer platform with a clean white design. Clone served on port 9321.`),
				How:    `Vercel.com was chosen as the target because of its clean, white aesthetic &mdash; a contrast to the dark themes (Linear) and colorful brands (Nando's) tested previously. Light-themed AI images were generated to match Vercel's minimal white style. The pipeline ran full block extraction first, then built the initial clone, then ran 3 iterative refinement rounds focusing on sections with the lowest SSIM scores. 893 lines of HTML/CSS, 12+ AI-generated images.`,
				Impact: template.HTML(`<strong>SSIM 0.787 with Overall 0.622 &mdash; tagged as the REFERENCE pipeline for all future clones.</strong> This is the first experiment that proved all techniques compose well together. The 3-iteration refinement pushed scores above the single-shot approaches. Light-themed image generation matched Vercel's white aesthetic correctly. This pipeline configuration is now the default starting point for any new website clone.`),
				Related: []string{"152", "155", "158", "163"},
			},
			{Num: "166", NumID: 166, Focus: "Hover State Capture", Result: "152 interactive elements, 91 CSS animations, 38 scroll-triggered", Cost: "\u2014", Finding: "Systematic extraction of all hover/interaction states from vercel.com: bg transparent to rgb(235,235,235), 0.15s transitions, 85 color transitions", HasDetail: true, Icon: "\U0001f5b1\ufe0f\U0001f4ab",
				Why:    `The optimized pipeline (<a href="/exp/165">Exp 165</a>) produced an excellent static clone, but real websites are interactive. Buttons change color on hover, dropdowns appear, elements animate on scroll. No previous experiment had systematically captured these interaction states. How many interactive behaviors does a modern website actually have, and can we extract them automatically?`,
				What:   `Playwright script visited vercel.com and systematically captured every interactive element and its hover/focus/active states. Results: <strong>152 interactive elements</strong>, <strong>91 CSS animations</strong>, <strong>38 scroll-triggered elements</strong>. Key hover patterns discovered: background transitions from transparent to rgb(235,235,235), color darkens on hover, all transitions use 0.15s duration.`,
				How:    `The extraction script hovered over every visible element and captured computed style changes (background-color, color, border-color, transform, opacity). CSS animation keyframes were extracted from stylesheets. Scroll-triggered elements were detected by scrolling the page incrementally and capturing newly-visible elements. Transition properties counted: color (85 elements), all (66), background (27), border-color (21), transform (21).`,
				Impact: template.HTML(`This is the first systematic catalog of interaction states for a production website. The data reveals that modern sites have far more interactive behaviors than static clones reproduce. The 0.15s transition duration and specific hover colors (rgb(235,235,235) not rgb(245,245,245)) provide exact values for clone fidelity. This feeds directly into <a href="/exp/167">Exp 167</a> for comparing original vs clone hover states.`),
				Related: []string{"165", "167"},
			},
			{Num: "167", NumID: 167, Focus: "Hover State Comparison", Result: "23 original vs 12 clone hover states", Cost: "\u2014", Finding: "Nav hover bg is rgb(235,235,235) not rgb(245,245,245) — 10 shades off. Missing: dropdown menus, framework expansion on hover", HasDetail: true, Icon: "\U0001f50d\u2194\ufe0f",
				Why:    `<a href="/exp/166">Experiment 166</a> extracted 152 interactive elements from the original Vercel site. But how many of those hover states did our clone actually reproduce? And where it tried, how accurate were the CSS values? This experiment directly compared original and clone hover behaviors to quantify the interaction gap.`,
				What:   `Side-by-side comparison of hover states: <strong>23 hover states in the original</strong> vs <strong>12 hover states in the clone</strong>. The clone captured roughly half the interactive behaviors. Where hover states existed in both, CSS values were compared pixel-by-pixel.`,
				How:    `The same Playwright hover extraction script from Exp 166 was run against the clone on port 9321. Results were diffed: which elements had hover states in original but not clone, which had them in both, and how close the CSS values matched. Key findings: nav hover background is rgb(235,235,235) in the original but rgb(245,245,245) in the clone &mdash; 10 shades off. Arrow centering needs display:flex + align-items:center + justify-content:center, which the clone missed.`,
				Impact: template.HTML(`The comparison quantifies exactly what "interaction fidelity" means. The clone reproduces 52% of hover states (12/23). Missing entirely: dropdown menus on hover, framework card expansion effects. The 10-shade-off nav hover color shows that even "close" hover states have measurable CSS gaps. This data establishes a baseline for future interactive cloning experiments and identifies specific CSS fixes that would close the gap.`),
				Related: []string{"165", "166"},
			},
			{Num: "168", NumID: 168, Focus: "Clone Notion.com", Result: "Overall 0.725, SSIM 0.732", Cost: "~$0.10", Finding: "275 blocks, 1277 lines, 5 AI images. Dark navy hero (rgb(33,49,131)), white sections, blue accent buttons. Light-themed images matching Notion's warm friendly style.", HasDetail: true, Icon: "\U0001f4dd\U0001f310",
				CloneShot: "/screenshots/exp168-notion/desktop.png", RefShot: "/screenshots/exp168-ref/desktop.png",
				Why:    `The optimized pipeline from <a href="/exp/165">Exp 165</a> was validated on Vercel (developer platform, white theme). Notion.com represents a different challenge: a productivity tool with a distinctive warm, friendly aesthetic, dark navy hero, and mixed light/dark sections. Could the pipeline handle this different visual language?`,
				What:   template.HTML(`<pre class="mermaid">
flowchart LR
    A["Notion.com"] --> B["275 Blocks Extracted"]
    B --> C["Claude: Build Clone"]
    C --> D["5 AI Images"]
    D --> E["2 Iteration Rounds"]
    E --> F["Overall 0.725"]
    style A fill:#f5e1e3,stroke:#d4727a
    style C fill:#fff3cd,stroke:#d4a017
    style F fill:#d4edda,stroke:#28a745
</pre>
Clone of Notion.com served on port 9322. 275 bounding boxes extracted from the live site. 1277 lines of HTML/CSS generated. 5 AI images created in Notion's warm, friendly illustration style. Dark navy hero section (rgb(33,49,131)) with white content sections and blue accent buttons.`),
				How:    `The reference pipeline from Exp 165 was applied: full block extraction, logical section grouping, Claude build, AI image generation, and 2 rounds of iterative refinement. Notion's distinctive dark navy hero (rgb(33,49,131)) was captured from the block data. AI images were generated in a light, warm style matching Notion's illustration aesthetic &mdash; friendly people, soft gradients, workspace scenes. The 2-iteration refinement was sufficient (vs 3 for Vercel) because Notion's layout is more structured and predictable.`,
				Impact: template.HTML(`<strong>Overall 0.725 with SSIM 0.732 at only ~$0.10 cost.</strong> This is the second-highest overall score in the research programme (after SaaS Tools at 0.789) and was achieved at one of the lowest costs. The result validates that the optimized pipeline generalises across visual styles: white minimal (Vercel, 0.787), warm illustrated (Notion, 0.725), dark minimal (Linear, 0.328), colorful brand (Nando's, 0.575). Notion's structured layout with clear section boundaries clones particularly well. The pipeline now has 8 validated clone targets across diverse industries and visual styles.`),
				Related: []string{"165", "152", "162"},
			},
			{Num: "170", NumID: 170, Focus: "Network Request Mapping", Result: "392 requests intercepted, 2 API endpoints found", Cost: "\u2014", Finding: "384 first-party, 2 API, 2 analytics, 4 third-party requests. POST /api/jwt (403), POST /api/stream/internal (200). 4 cookies, 10 console errors, 1 third-party blob storage domain.", HasDetail: true, Icon: "\U0001f310\U0001f50c",
				Why:    `Previous experiments extracted visual structure (bounding boxes, hover states, CSS) but ignored the network layer entirely. A real website is more than pixels &mdash; it makes API calls, loads third-party scripts, sets cookies, and logs errors. Understanding the full API surface is essential for full-stack cloning and automated API documentation. Could Playwright network interception capture the complete request map of a production site?`,
				What:   `Playwright intercepted all network traffic on vercel.com: <strong>392 total requests</strong> (384 first-party, 2 API, 2 analytics, 4 third-party). Found 2 API endpoints: <code>POST /api/jwt</code> (403 Forbidden) and <code>POST /api/stream/internal</code> (200 OK). 4 cookies detected, 10 console errors logged, 0 JS errors, 1 third-party blob storage domain.`,
				How:    `Playwright's network interception captured every request during page load, categorizing by domain (first-party vs third-party), type (API, analytics, static), and response status. Console messages were collected via the page.on('console') event. Cookies were extracted from the browser context. The result is a complete network fingerprint of the site.`,
				Impact: template.HTML(`Playwright network interception captures the full API surface of a site &mdash; enables automated API documentation and full-stack cloning. The 384:4 first-party-to-third-party ratio shows Vercel is remarkably self-contained. The JWT endpoint (403) and streaming endpoint (200) reveal the auth and real-time architecture. Combined with <a href="/exp/166">Exp 166</a> (hover states) and <a href="/exp/133">Exp 133</a> (CSS extraction), we now have three complementary data layers: visual, interactive, and network.`),
				Related: []string{"166", "133", "165", "171"},
			},
			{Num: "171", NumID: 171, Focus: "Route Discovery + Security Audit", Result: "Security headers audited, HSTS + Permissions-Policy missing", Cost: "\u2014", Finding: "CSP OK, X-Frame-Options DENY, nosniff OK, Referrer-Policy OK. HSTS missing, Permissions-Policy missing. No exposed API keys, source maps, or error traces. 15 external domains linked. Session cookie not httpOnly.", HasDetail: true, Icon: "\U0001f6e1\ufe0f\U0001f50d",
				Why:    `<a href="/exp/170">Experiment 170</a> mapped the network requests; this experiment maps the routes and checks security posture. Before cloning a site's functionality, we need to know: what pages exist, what links go where, and how hardened is the security? Automated security auditing catches issues that manual review misses &mdash; and the findings inform what a clone should replicate (good headers) and what it should fix (missing headers).`,
				What:   `Discovered all internal and external links on vercel.com and checked status codes. Security header audit: <strong>CSP OK</strong>, <strong>X-Frame-Options DENY</strong>, <strong>X-Content-Type-Options nosniff OK</strong>, <strong>Referrer-Policy OK</strong>. Missing: <strong>HSTS</strong> and <strong>Permissions-Policy</strong>. No exposed API keys, internal URLs, source maps, or error traces in page source. 15 external domains linked.`,
				How:    `Playwright visited vercel.com and extracted all anchor hrefs, categorizing as internal routes vs external links. HTTP response headers were inspected against OWASP security header recommendations. Page source was scanned for common sensitive patterns: API keys (sk_live, AKIA, etc.), internal URLs, .map files, and stack traces. Cookie attributes (httpOnly, secure, sameSite) were checked for each cookie.`,
				Impact: template.HTML(`Automated security audit catches missing headers (HSTS, Permissions-Policy) that manual review might miss. The session cookie not being httpOnly is a potential XSS vector. Combined with <a href="/exp/170">Exp 170</a> (network mapping), we now have a complete site intelligence package: visual structure, interaction states, network traffic, route map, and security posture. This is everything needed to understand a site before cloning or competing with it.`),
				Related: []string{"170", "166", "165"},
			},
			{Num: "172", NumID: 172, Focus: "AI Video Generation (Veo 2)", Result: "4 videos from static AI images, ~5.6MB total, ~30s per video", Cost: "\u2014", Finding: "Google Veo 2 via predictLongRunning API generates 5s MP4 loops from a single PNG + animation prompt. Seamless loop prompt: animation must END in the EXACT SAME state as it STARTS. Vercel: rotating globe (1.3MB), shimmering prism (1.7MB). Notion: animated dashboard (1.0MB), animated agents (1.6MB). Videos embedded as <video autoplay loop muted playsinline>. Combined with CSS animations (fade-in, scroll entrance, hover transitions), clones become dynamic product pages.", HasDetail: true, Icon: "\U0001f3ac\U0001f310",
				Why:    `Previous clones used static AI-generated images. Real product pages use hero videos, animated backgrounds, and motion graphics to create premium feel. Could Google Veo 2 generate seamless looping videos from the existing AI images, turning static clones into dynamic product pages?`,
				What:   `Generated <strong>4 videos</strong> from static AI images using Google Veo 2 API. Image-to-video: feed existing PNG + animation prompt &rarr; 5s MP4 loop. Vercel clone: rotating globe (1.3MB), shimmering prism (1.7MB). Notion clone: animated dashboard (1.0MB), animated agents (1.6MB). Total: <strong>~5.6MB</strong> of video, <strong>~30s generation</strong> per video.`,
				How:    `Used Veo 2 via the predictLongRunning API. Each call takes an existing AI-generated PNG and an animation prompt describing the desired motion. Seamless loop prompts use the constraint: "animation must END in the EXACT SAME state as it STARTS." Videos are embedded as <code>&lt;video autoplay loop muted playsinline&gt;</code> replacing static <code>&lt;img&gt;</code> tags. Combined with CSS animations (fade-in, scroll-triggered entrance, hover transitions) for full dynamic feel.`,
				Impact: template.HTML(`<strong>First AI-generated video embedded in AI-cloned websites.</strong> Veo 2 generates 5s seamless loops from a single image in ~30s. Combined with CSS animations, static AI clones become dynamic product pages indistinguishable from hand-crafted sites with professional video production. The video layer is the final piece: structure (HTML) + style (CSS) + images (DALL-E/Flux) + video (Veo 2) + interaction (JS) = complete website reproduction. <a href="/exp/172b">Exp 172b</a> later proved that ffmpeg post-processing (not AI loop instructions) is the correct approach &mdash; compress with libx264 and let the browser handle looping.`),
				Related: []string{"165", "168", "170", "171", "172b"},
			},
			{Num: "172b", NumID: 172, Focus: "Video Post-Processing with ffmpeg", Result: "~50% file size reduction, reliable looping via browser", Cost: "\u2014", Finding: "AI should NOT be asked to fade or loop — produces worse results. ffmpeg handles all post-processing: libx264 compress (~50% size reduction, 243KB-815KB from 1-2MB raw). Browser <video autoplay loop muted playsinline> handles looping. Crossfade via xfade filter unreliable. \"End same as start\" prompts confuse Veo 2. fade=t=out creates flash on loop. Best approach: clean AI generation (describe motion only) + ffmpeg compress + browser loop.", HasDetail: true, Icon: "\U0001f3ac\u2702\ufe0f",
				Why:    `<a href="/exp/172">Experiment 172</a> generated AI videos with Veo 2 and used "animation must END in the EXACT SAME state as it STARTS" in prompts to get seamless loops. But the results were inconsistent &mdash; some videos had jerky motion at the loop point. The question: should the AI handle looping and fading, or should we let the AI focus on motion and handle everything else in post-processing with ffmpeg?`,
				What:   template.HTML(`Tested multiple approaches to video post-processing:<br><br><strong>Compression</strong>: <code>ffmpeg -i raw.mp4 -c:v libx264 -preset fast -crf 23 -movflags +faststart output.mp4</code> reduces file size by ~50% (from 1-2MB raw to 243KB-815KB).<br><br><strong>Browser looping</strong>: <code>&lt;video autoplay loop muted playsinline&gt;</code> handles seamless looping natively &mdash; no need to encode loop points into the video itself.<br><br><strong>Crossfade attempt</strong>: ffmpeg xfade filter to blend start/end frames of longer (8s) videos. Complex filter graph, did not work reliably &mdash; simpler to let the browser loop the raw clip.<br><br><strong>Fade-to-black attempt</strong>: <code>fade=t=out</code> creates a visible flash/black frame on loop &mdash; do NOT use.`),
				How:    template.HTML(`Generated videos with Veo 2 using different prompt strategies, then applied ffmpeg post-processing. Key prompt findings:<br><br><strong>"End same as start"</strong> in prompts confuses Veo 2 and creates jerky motion at the loop boundary.<br><strong>"Not too fast"</strong> produces better results than "seamless loop" &mdash; the AI focuses on smooth motion rather than trying to engineer a loop point.<br><strong>"Describe motion only"</strong> (e.g., "gentle rotation", "soft shimmer") produces the cleanest output.<br><br>The ffmpeg pipeline: raw Veo 2 MP4 &rarr; libx264 compress with CRF 23 &rarr; faststart flag for web streaming &rarr; embed with HTML5 video tag. The browser&rsquo;s native loop attribute handles repetition.`),
				Impact: template.HTML(`<strong>The division of labour is clear: AI generates motion, ffmpeg compresses, browser loops.</strong> Asking the AI to handle fading or looping produces worse results than letting it focus purely on motion generation. File sizes after compression range from 243KB to 815KB &mdash; small enough for production web pages. The <code>-movflags +faststart</code> flag ensures the video plays immediately without buffering the entire file. fade=t=out is a trap &mdash; it creates a black flash on every loop iteration. The optimal pipeline is: clean prompt (describe motion only, no loop/fade instructions) &rarr; Veo 2 generates 5-8s clip &rarr; ffmpeg libx264 compress &rarr; <code>&lt;video autoplay loop muted playsinline&gt;</code> in HTML.`),
				Related: []string{"172", "165", "168"},
			},
			{Num: "173", NumID: 173, Focus: "Professional Site Analysis Docs", Result: "Port 9323, 1312+ lines, 7 pages, Security Grade B (80%)", Cost: "\u2014", Finding: "Dedicated documentation site for network mapping, route discovery, security audit, API analysis, and executive reporting. Professional dark sidebar design with Mermaid diagrams. 7 pages: Overview, Network, Routes, Security, API, Pentest Report, Executive Report. Security grade B (80%), 4 findings (2 medium, 2 low). Persona-based reports: pentester + CTO perspectives.", HasDetail: true, Icon: "\U0001f4d1\U0001f50d",
				Why:    template.HTML(`Experiments <a href="/exp/170">170</a> (network mapping), <a href="/exp/171">171</a> (route discovery + security audit) produced raw data &mdash; JSON dumps of requests, headers, and findings. Useful for us, unreadable for anyone else. If this analysis pipeline is going to be a product, the output needs to be professional documentation that a client or team lead can read without parsing JSON. Could we build a Swagger/Postman-style docs site that presents the analysis data as polished, navigable reports?`),
				What:   template.HTML(`Built a <strong>1312+ line Go HTTP server</strong> on port 9323 with 7 documentation pages:<br><br><strong>Overview</strong>: dashboard summary with key metrics and Mermaid architecture diagram<br><strong>Network</strong>: full request map with domain breakdown, API endpoints, cookie analysis<br><strong>Routes</strong>: internal/external link map, status codes, redirect chains<br><strong>Security</strong>: header audit results, OWASP compliance, cookie security flags<br><strong>API</strong>: discovered endpoints with request/response details, authentication analysis<br><strong>Pentest Report</strong>: persona-based security assessment from a penetration tester perspective<br><strong>Executive Report</strong>: CTO-level summary with risk ratings and remediation priorities`),
				How:    template.HTML(`Single Go binary with embedded HTML/CSS &mdash; the same pattern used across all clone experiments. Professional dark sidebar navigation (inspired by Swagger UI and Postman docs). Mermaid diagrams for architecture visualization and request flow. Security findings scored and graded: <strong>Grade B (80%)</strong> overall, with 4 findings: 2 medium severity (missing HSTS, session cookie not httpOnly), 2 low severity (missing Permissions-Policy, incomplete CSP). Two persona-based report views: the <strong>pentester report</strong> provides technical detail and exploitation paths; the <strong>executive report</strong> provides business risk context and prioritized remediation for CTO-level readers.`),
				Impact: template.HTML(`<strong>Analysis data becomes a deliverable product.</strong> Raw JSON from Exp 170-171 transforms into professional documentation that could be handed to a client. The persona-based reports (pentester vs CTO) demonstrate that the same data serves different audiences when framed correctly. The security grading system (A-F based on findings) provides an instant health check. At 1312+ lines this is a substantial application &mdash; proving the pipeline can build not just clones but original documentation tools. Combined with the network, route, and security extraction experiments, this completes the site analysis toolkit: extract &rarr; analyze &rarr; present.`),
				Related: []string{"170", "171", "166", "165"},
			},
		},
	},
	{
		Slug:     "business-pipeline",
		Name:     "Business Pipeline",
		Icon:     "\U0001f4ca",
		ExpRange: "Experiments 55 – 84",
		Narrative: []template.HTML{
			`Code is only half a product. The business pipeline experiments explored whether AI can handle everything from idea discovery to legal documents. The answer is yes, though the quality varies by domain. Docker containerization (Experiment 55), TypeScript compilation (Experiment 57), PostgreSQL schema generation (Experiment 58), and authentication flows (Experiment 61) all worked on first attempt. The feedback loop experiment (56) showed that AI can read its own test output and fix failures without human guidance.`,
			`On the marketing side, experiments 68 through 78 built a complete go-to-market pipeline: idea discovery from market trends, competitive research, landing page copy, pricing strategy, advertising copy, video scripts, SEO content, Product Hunt launch strategy, and distribution channel analysis. Each component cost between $0.005 and $0.02, making the entire marketing pipeline under $0.20. The quality is sufficient for first-draft material that a human can refine &mdash; not for final publication, but far better than starting from a blank page.`,
			`The <strong>legal documents</strong> (Experiment 80) produced Terms of Service, Privacy Policy, Cookie Policy, and Acceptable Use Policy tailored to a specific SaaS product. Support infrastructure (Experiment 79) generated 10 help articles, 20 chatbot response templates, and escalation routing rules. Localisation (Experiment 81) handled i18n extraction and translation. The org-personas experiment (83) generated realistic company personas for B2B targeting. Together, these experiments show that AI can scaffold the entire non-code surface area of a product launch.`,
		},
		KeyInsight: `The full pipeline from idea discovery to legal documents costs under $1. Every non-code artifact a startup needs can be AI-generated as a solid first draft.`,
		TableType:  "business-pipeline",
		Experiments: []Experiment{
			{Num: "55", NumID: 55, Focus: "Docker", Category: "Infrastructure", Finding: "Dockerfile, compose, multi-stage builds", HasDetail: true, Icon: "\U0001f433", SourceFile: "exp55-docker.sh",
				Why:     `We had a Go binary that passed all tests on a developer machine. But "works on my machine" is not a deployment strategy. We needed to prove that the pipeline could take a passing application and produce a running Docker container automatically, with no human intervention between test-green and container-live.`,
				What:    `We generated a multi-stage Dockerfile (golang:1.22-alpine builder, alpine:3.19 runtime), built the image, started the container, and ran health checks and functional tests against it. The input was the CRM application from <a href="/exp/38">Experiment 38</a>. The output was a running container serving HTTP on port 8081.`,
				How:     template.HTML(`The script generates a Dockerfile, runs <code>docker build</code>, starts the container with <code>docker run -d</code>, waits 3 seconds, then hits the health endpoint, creates a client via the API, lists clients, and checks the dashboard HTML. Image size is measured to verify the multi-stage build is working.`),
				Impact:  `The image built, the container ran, the health check passed, and API requests worked. The multi-stage build produced a container around 20MB &mdash; just the Go binary plus Alpine. This became <a href="/exp/913">Step 13</a> of the pipeline. The pattern is now used across 19 deployed products.`,
				Related: []string{"38"},
			},
			{Num: "56", NumID: 56, Focus: "Feedback Loop", Category: "Infrastructure", Finding: "AI reads test output, fixes failures", HasDetail: true, SourceFile: "exp56-feedback-loop.py",
				Why:     `Building the app is one thing; knowing whether users actually want what we built is another. We wanted to test whether AI personas could act as surrogate users &mdash; reviewing a running application, giving feedback, and driving iterative improvement &mdash; without any human involvement.`,
				What:    `Three AI personas (freelance designer, agency owner, accountant) reviewed the running CRM application over 3 feedback rounds. Each round: personas inspect the app, give structured feedback, and a product manager persona synthesises their input into prioritised changes.`,
				How:     `Each round, the script fetches the app's HTML, strips it to text, and sends it to each persona for in-character review. Personas answer: first impression, what works, what is missing, what is confusing, would you pay $15/month, and one feature request. A synthesis step then ranks the top 3 changes by effort and impact.`,
				Impact:  `Nine feedback sessions across 3 rounds produced actionable improvement lists with specific priorities. The personas consistently surfaced different needs &mdash; the designer wanted speed, the agency owner wanted team features, the accountant wanted tax summaries. This validates that AI feedback loops can simulate iterative product development.`,
				Related: []string{"23"},
			},
			{Num: "57", NumID: 57, Focus: "TypeScript", Category: "Infrastructure", Finding: "tsconfig, compilation, type checking", HasDetail: true, SourceFile: "exp57-typescript.py",
				Why:     `Our entire pipeline was built around Go. We needed to know whether the same patterns &mdash; spec-driven generation, test-first, auto-fix &mdash; would transfer to a completely different language. TypeScript with Node.js was the obvious test case because it has its own compilation step and type system.`,
				What:    template.HTML(`We asked Claude (Haiku) to build the same CRM in TypeScript/Node.js: store layer, HTTP server, tests &mdash; all using zero npm dependencies, just Node.js 22+ built-ins (<code>node:http</code>, <code>node:test</code>, <code>node:assert</code>). The goal was 4 files: package.json, store.ts, server.ts, and server.test.ts.`),
				How:     template.HTML(`The script runs <code>claude -p</code> with a detailed prompt specifying the file structure, then checks whether all 4 files were created and whether <code>npx tsc --noEmit</code> passes. We measured file count, line count, compilation success, and cost.`),
				Impact:  `The pipeline generalises beyond Go. The AI produced a TypeScript CRM with the same architecture &mdash; in-memory store, HTTP server, embedded HTML. This confirms that our spec-driven approach is language-agnostic; only the auto-fix tools need to change per language.`,
				Related: []string{"38"},
			},
			{Num: "58", NumID: 58, Focus: "PostgreSQL", Category: "Infrastructure", Finding: "Schema generation, migrations", HasDetail: true, SourceFile: "exp-batch-58-62-84.py",
				Why:     `Every application we built used in-memory storage &mdash; fast and simple, but data is lost on restart. For any real SaaS product, we need to prove the AI can plan a migration to PostgreSQL without introducing the complexity that makes cheap models fail.`,
				What:    `We asked the AI to produce a complete PostgreSQL migration plan: schema design with CREATE TABLE statements, connection pooling configuration, repository pattern interfaces, transaction handling, and a testing strategy using both in-memory and Postgres implementations.`,
				How:     `A single AI call generated the full migration document covering schema design, Go code changes (interface-based store swapping), query patterns (prepared statements, no string concatenation), and effort estimates. The plan preserves the existing MemoryStore for unit tests while adding a PostgresStore for production.`,
				Impact:  `The AI produced a practical, implementable migration plan. The repository pattern it suggested &mdash; both MemoryStore and PostgresStore implementing the same interface &mdash; is exactly how we structure stores in production. This proves the AI can plan infrastructure upgrades, not just generate greenfield code.`,
				Related: []string{"27", "36"},
			},
			{Num: "59", NumID: 59, Focus: "AI Code Review", Category: "Quality", Finding: "Automated review with actionable fixes", HasDetail: true, SourceFile: "exp59-ai-code-review.py",
				Why:     `One model wrote the code in <a href="/exp/37">Experiment 37</a>. We wanted to know what a different model would find wrong with it. Specifically: can AI reviewers catch real issues with specific severity ratings and actionable fixes, or do they just produce generic advice?`,
				What:    `Three specialised AI reviewers examined the CRM code independently: a Senior Go Engineer (idioms, error handling, naming), a Security Specialist (XSS, input validation, timeouts), and a Performance Engineer (data structure efficiency, mutex usage, scalability).`,
				How:     `Each reviewer received the store and handler code with a focused prompt for their domain. They produced structured findings with severity ratings (Critical/Major/Minor) and verdicts (APPROVE / REQUEST CHANGES). We counted findings by severity across all three reviewers.`,
				Impact:  `The reviewers produced specific, actionable findings &mdash; not generic advice. The Security Specialist found real XSS vectors and missing timeouts. The Go Engineer caught non-idiomatic patterns. This became the foundation for the post-code review step in the pipeline (<a href="/exp/908">Step 8</a>).`,
				Related: []string{"37", "31"},
			},
			{Num: "60", NumID: 60, Focus: "Creative Writing", Category: "Content", Finding: "Product descriptions, taglines", HasDetail: true, SourceFile: "exp60-creative-agent.py",
				Why:     `Products need personality, not just functionality. We wanted to test whether an AI "creative agent" could generate the small touches that make users smile &mdash; confetti animations, witty empty states, birthday reminders &mdash; and whether a product manager could then evaluate them objectively.`,
				What:    `A creative persona ("Zara") generated 15 delight features for the CRM, then a product manager persona scored each on user value, dev effort, and revenue impact. Features ranged from gamification (activity streaks) to automation (overdue invoice reminders) to personality (funny empty states).`,
				How:     `Two AI calls: one at temperature 0.6 for creative generation, one at temperature 0.3 for rational evaluation. The evaluation used a scoring formula: (user value x revenue impact) / dev effort. Features were ranked and sorted into build-in-v1, defer-to-v2, and never-build categories.`,
				Impact:  `The creative agent produced genuinely useful ideas alongside predictable ones. Features like "auto-remind about overdue invoices" and "monthly invoiced total on dashboard" scored highest because they combine high user value with low effort. This shows AI can handle the creative product layer, not just the technical one.`,
			},
			{Num: "61", NumID: 61, Focus: "Auth Flows", Category: "Infrastructure", Finding: "Login, registration, JWT, sessions", HasDetail: true, SourceFile: "exp61-auth.py",
				Why:     `Every application we built had zero authentication &mdash; every endpoint was publicly accessible. Before shipping any SaaS product, we needed a complete auth specification. The question was whether the AI could produce a security-grade auth plan, not just a tutorial-level one.`,
				What:    template.HTML(`A security architect persona reviewed the CRM code and produced a complete auth system specification: registration (Argon2id hashing), login (JWT vs sessions), forgot password (crypto/rand tokens), SSO (OAuth2/OIDC), MFA (TOTP + backup codes), and role-based authorization. Each feature was rated as MVP-required, launch-required, or nice-to-have.`),
				How:     `Two steps: first, the AI generated comprehensive auth requirements from the actual application code. Second, we ran live pen tests against the running app &mdash; hitting every endpoint without credentials to document the current exposure. Every endpoint returned 200 with no auth challenge.`,
				Impact:  `The pen test confirmed the obvious: 0 of 8 endpoints were protected. But the auth specification was thorough and production-grade &mdash; it covered Argon2id, httpOnly cookies, account lockout, TOTP with backup codes, and resource ownership scoping. This feeds directly into the security gate at <a href="/exp/912">Step 12</a>.`,
				Related: []string{"42"},
			},
			{Num: "62", NumID: 62, Focus: "Regression Suite", Category: "Quality", Finding: "Automated regression test generation", HasDetail: true, SourceFile: "exp-batch-58-62-84.py",
				Why:     `We had tests at multiple layers &mdash; unit, integration, Playwright, adversarial, mutation &mdash; but no unified strategy for running them all after every build. We needed a regression suite design that defines gate criteria, handles flaky tests, and enforces coverage thresholds.`,
				What:    `The AI designed a 7-layer regression suite: unit tests (30, under 1s), HTTP integration (34, under 5s), Playwright journeys (10, under 30s), mobile viewport (6), console errors (5), adversarial (47), and weekly mutation testing (14 mutations). Gate criteria require zero unit/HTTP failures and minimum 80% coverage.`,
				How:     `A single AI call produced the complete suite design including CI pipeline configuration, flaky test management (quarantine + auto-retry), test data isolation (fresh store per test, no shared state), and coverage enforcement. The design matched our actual testing pyramid from experiments 35-45.`,
				Impact:  `The regression suite design codifies what we learned across 20+ testing experiments into a single, runnable specification. The layered approach &mdash; fast tests first, expensive tests last &mdash; means we catch most issues in under 5 seconds and only run Playwright when the cheap checks pass.`,
				Related: []string{"35", "40", "45"},
			},
			{Num: "63", NumID: 63, Focus: "Auto Documentation", Category: "Content", Finding: "API docs, README, architecture docs", HasDetail: true, SourceFile: "exp63-auto-docs.py",
				Why:     `We had working code but zero documentation. No README, no API docs, no user guide, no architecture decision records. A product without documentation is not shippable. We wanted to see how many documentation artifacts the AI could generate directly from the source code.`,
				What:    `Four documents generated from the CRM source code: a README (what, quick start, API, architecture), API documentation (every endpoint with request/response examples), a user guide (step-by-step for non-technical users), and Architecture Decision Records (6 ADRs explaining why we chose Go stdlib, in-memory store, etc.).`,
				How:     `Four sequential AI calls, each receiving the relevant source code. The README and API docs were generated from handler code. The user guide was generated from the route structure. The ADRs were generated from a list of architectural choices. Each call cost roughly $0.005.`,
				Impact:  `Four complete, accurate documents generated from source code alone. The API docs correctly listed every endpoint with methods, request formats, and response examples. The ADRs explained trade-offs that were implicit in the code but never written down. This became a standard pipeline step &mdash; documentation is generated alongside the code, not after.`,
				Related: []string{"67"},
			},
			{Num: "64", NumID: 64, Focus: "Evidence Trail", Category: "Quality", Finding: "Audit log of all AI decisions", HasDetail: true, SourceFile: "exp64-evidence-trail.py",
				Why:     `Our pipeline runs 13 steps with dozens of AI calls, but there was no audit trail. If something goes wrong in production, we need to trace back: which reviewer approved it? What did the security scan find? Who signed off? We needed auto-generated compliance documentation from pipeline outputs.`,
				What:    `Two artifacts: per-stage sign-off checklists (checkboxes, reviewer assignments, evidence requirements, pass criteria for each stage) and an aggregate quality gate report (executive summary, per-stage verdicts, open issues, risk assessment, and a ship/no-ship decision).`,
				How:     `Two AI calls: one to generate stage checklists covering all 7 pipeline stages (discovery through documentation), and one to generate an aggregate report from simulated stage results. The aggregate report included PASS/CONDITIONAL/FAIL verdicts and a risk-ordered remediation plan.`,
				Impact:  `The AI produced auditor-ready documentation. The aggregate report correctly identified browser testing and GDPR compliance as conditional passes, flagged the 3 unresolved Playwright failures, and recommended "SHIP WITH CONDITIONS." This is essential for any regulated deployment where you need to prove due diligence.`,
			},
			{Num: "65", NumID: 65, Focus: "Test Report", Category: "Quality", Finding: "Human-readable test result summaries", HasDetail: true, SourceFile: "exp-batch-65-66-69.py",
				Why:     `Raw test output is useful for developers but meaningless to stakeholders. We needed the AI to take our test data from multiple layers &mdash; unit, integration, Playwright, adversarial, chaos, mutation, coverage &mdash; and synthesise it into an executive summary that a product manager can read.`,
				What:    `A comprehensive test report generated from real test results: executive summary, results by category (table), coverage analysis, mutation analysis, individual failure root-cause analysis, and top 5 testing improvement recommendations.`,
				How:     `A single AI call received structured test data (32/32 unit pass, 7/10 Playwright pass, 89% mutation score, 57.4% coverage) and produced a multi-section report. Each Playwright failure got its own root-cause analysis with severity and suggested fix.`,
				Impact:  `The report correctly identified that store coverage was 0% (all tests were in main.go) and that the mutation score gap was in date handling. This kind of synthesis turns raw metrics into actionable priorities &mdash; exactly what a team standup needs to decide where to invest testing effort next.`,
				Related: []string{"35", "45"},
			},
			{Num: "66", NumID: 66, Focus: "Security Audit", Category: "Security", Finding: "OWASP checklist, vulnerability scan", HasDetail: true, SourceFile: "exp-batch-65-66-69.py",
				Why:     `We had security findings scattered across multiple experiments &mdash; gosec output, govulncheck CVEs, adversarial test results, pen test findings, GDPR gaps. We needed the AI to consolidate everything into a single security audit report with OWASP mapping and remediation priorities.`,
				What:    `A consolidated security audit report covering OWASP Top 10 checklist status, a findings table with severity and remediation status, critical path items, accepted risks with justification, GDPR/SOC2 compliance status, and a risk-ordered remediation plan.`,
				How:     `A single AI call received all security findings from experiments 40-52 and produced a structured audit report. The input included gosec (12 issues), govulncheck (13 CVEs), adversarial results, pen test findings (4 High), missing headers, GDPR non-compliance (10 items), and the zero-auth status.`,
				Impact:  `The report mapped findings to OWASP categories, correctly prioritised authentication as the critical blocker, and separated "fix before launch" from "accepted for MVP." This is the format we use in production: every deployed product gets a security audit summary before going live.`,
				Related: []string{"42", "45", "50"},
			},
			{Num: "67", NumID: 67, Focus: "User Documentation", Category: "Content", Finding: "User guides, onboarding flows", HasDetail: true, SourceFile: "exp-batch-67-72-74.py",
				Why:     `The auto-documentation in <a href="/exp/63">Experiment 63</a> generated developer-facing docs. But real users do not read API documentation. We needed a user guide written for someone who has never used a CRM &mdash; plain language, step-by-step, no jargon.`,
				What:    `A complete user guide covering getting started, managing clients (add/edit/delete/search), activity tracking (logging calls, emails, meetings), invoicing (create, send, pay, void, print), and a 10-question FAQ. Written for non-technical users with simple language.`,
				How:     `A single AI call received the handler code and produced a step-by-step guide. The prompt explicitly required non-technical language: "Write for someone who has NEVER used a CRM before." The output covered every feature in the application with click-by-click instructions.`,
				Impact:  `The user guide was immediately usable as onboarding documentation. Combined with <a href="/exp/63">Experiment 63</a>'s developer docs, we now generate both technical and user-facing documentation from source code. The total cost for all documentation (developer + user) is under $0.05.`,
				Related: []string{"63"},
			},
			{Num: "68", NumID: 68, Focus: "Idea Discovery", Category: "Marketing", Finding: "Market trend analysis, opportunity scoring", HasDetail: true, SourceFile: "exp68-idea-discovery.py",
				Why:     `The pipeline builds products, but someone has to decide what to build. We wanted to test whether AI could simulate the earliest stage of product development: discovering pain points, clustering them into opportunities, and producing a build-ready brief.`,
				What:    `A three-step discovery pipeline: simulate 20 Reddit/HN posts with realistic user frustrations, cluster them into 5 product opportunities scored by pain frequency and willingness to pay, then generate a one-page product brief for the top idea &mdash; ready to feed into the build pipeline.`,
				How:     `Three sequential AI calls. The first simulates social media posts across subreddits (r/freelance, r/SaaS, r/smallbusiness). The second clusters pain points and scores opportunities by (pain frequency x willingness to pay) / build complexity. The third writes a brief with problem, solution, features, revenue model, and success metrics.`,
				Impact:  `The output brief was directly feedable into our pipeline from <a href="/exp/27">Experiment 27</a> onwards. The simulated pain points were realistic enough to surface genuine product categories. This closes the loop: the pipeline can now go from "find something to build" to "deployed product" with zero human input.`,
				Related: []string{"27", "69"},
			},
			{Num: "69", NumID: 69, Focus: "Market Research", Category: "Marketing", Finding: "Competitive landscape, TAM/SAM/SOM", HasDetail: true, SourceFile: "exp-batch-65-66-69.py",
				Why:     `<a href="/exp/68">Experiment 68</a> found a product idea. Before building it, we need to know: how big is the market, who are the competitors, and what should we charge? This is the research step that separates building something nobody wants from building into a validated opportunity.`,
				What:    `A market research report covering TAM/SAM/SOM estimates, a 5-competitor landscape analysis (price, users, strengths, weaknesses), pricing recommendations (Free/Pro/Business tiers based on competitor analysis), go-to-market strategy, risk assessment, and a BUILD/PIVOT/KILL verdict.`,
				How:     `A single AI call received the product brief from Experiment 68 and produced structured market research. The prompt specified concrete output formats: competitor comparison table, pricing tiers with justification, and channel-specific launch plans.`,
				Impact:  `The research report produced realistic market sizing and competitor analysis. Combined with <a href="/exp/68">Experiment 68</a>'s idea discovery, we now have a two-step validation process that costs under $0.02 total. The BUILD/PIVOT/KILL verdict gives the pipeline a go/no-go gate before committing to development.`,
				Related: []string{"68"},
			},
			{Num: "70", NumID: 70, Focus: "Landing Page", Category: "Marketing", Finding: "Copy, hero, CTA, social proof", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `A product needs a landing page before it needs code. We wanted to test whether the AI could produce a complete landing page specification &mdash; hero section, features, pricing, testimonials, FAQ &mdash; with copy that is specific enough to use directly.`,
				What:    `A landing page spec with HTML structure and copy for every section: hero (headline, subheadline, CTA), 6 features with icons, pricing table (Free/Pro/Business tiers), 3 realistic testimonials, 5 FAQ entries, and footer. All copy tailored to a freelancer CRM at $15/month.`,
				How:     `A single AI call with a detailed prompt specifying every section. The output included specific copy, not placeholders &mdash; actual headlines, button text, testimonial quotes with names and roles, and FAQ answers.`,
				Impact:  `The landing page copy was immediately usable as a first draft. The testimonials sounded authentic, the feature descriptions were benefit-focused rather than feature-focused, and the pricing table had clear tier differentiation. This is good enough to launch with and A/B test from.`,
				Related: []string{"71", "72"},
			},
			{Num: "71", NumID: 71, Focus: "Pricing Strategy", Category: "Marketing", Finding: "Tier design, anchor pricing, free trial", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `"How much should we charge?" is the question most founders get wrong. We asked the AI to apply pricing frameworks &mdash; Van Westendorp analysis, anchor pricing, feature gating &mdash; to our CRM and produce a defensible pricing strategy, not just a number.`,
				What:    `A pricing strategy covering simulated Van Westendorp analysis (too cheap / cheap / expensive / too expensive thresholds), 3-tier design (Free/Pro/Business with feature gating), annual vs monthly discount, competitor pricing comparison, free trial strategy, and upgrade triggers.`,
				How:     `A single AI call with a structured prompt. The AI simulated price sensitivity research, compared against 5 real-world competitors, and designed tiers where the free plan demonstrates value and the paid plans gate features that power users demand (like team access and bulk operations).`,
				Impact:  `The pricing strategy was internally consistent and commercially sensible. The feature gating aligned with what personas demanded in <a href="/exp/23">Experiment 23</a> &mdash; features that power users and agencies need are the ones behind the paywall. This validates that AI can handle pricing strategy, not just pricing arithmetic.`,
				Related: []string{"70", "78"},
			},
			{Num: "72", NumID: 72, Focus: "Ad Copy", Category: "Marketing", Finding: "Google, Meta, LinkedIn ad variants", HasDetail: true, SourceFile: "exp-batch-67-72-74.py",
				Why:     `A product with no advertising gets no users. We needed multiple ad variants to A/B test across different platforms and audiences. Writing 5 ad variants from scratch takes a copywriter hours; we wanted to see if the AI could produce platform-specific variants in a single call.`,
				What:    template.HTML(`Five ad variants, each with a different angle: pain-focused ("Stop losing track of clients"), benefit-focused ("Get paid faster"), social proof ("Join 1000+ freelancers"), urgency ("First month free"), and comparison ("HubSpot is too complex"). Each includes headline, subheadline, body copy, CTA, DALL-E prompt for the hero image, and target platform.`),
				How:     `A single AI call with explicit variant themes and platform targeting (Google Ads, Facebook, LinkedIn, Reddit). Each variant was constrained: headline max 10 words, subheadline max 20, body copy 50 words. The DALL-E prompts describe the visual style for each variant.`,
				Impact:  `Five distinct, platform-appropriate ad variants with matching visual prompts. The pain-focused and comparison variants were strongest &mdash; they address specific frustrations rather than generic benefits. The DALL-E prompts are specific enough to generate consistent hero images. Total cost: under $0.01.`,
				Related: []string{"70", "74"},
			},
			{Num: "73", NumID: 73, Focus: "Video Script", Category: "Marketing", Finding: "60s explainer, scenes, voiceover", HasDetail: true, SourceFile: "exp-batch-73-75-79-80.py",
				Why:     `A 60-second product video is the highest-converting marketing asset for SaaS. We wanted to see if the AI could produce a complete video script with scene breakdowns, voiceover text, visual direction, and music mood &mdash; ready to hand to a video editor or text-to-speech tool.`,
				What:    `A 5-scene video script: hook (0-10s), problem (10-25s), solution with product demo (25-45s), social proof (45-55s), and CTA (55-60s). Each scene includes visual description, voiceover text, and music direction. Also includes voice direction, thumbnail design, and 3 title options.`,
				How:     `A single AI call with a structured scene-by-scene format. The prompt specified exact timing per scene and required both visual and audio direction. The AI produced a complete script that could be executed with screen recording, text-to-speech, and stock footage.`,
				Impact:  `The script followed standard SaaS demo structure and the timing was realistic. The voiceover text was natural and benefit-focused. Combined with DALL-E prompts from <a href="/exp/72">Experiment 72</a>, we have the raw materials for a complete marketing video for under $0.02 total.`,
				Related: []string{"72"},
			},
			{Num: "74", NumID: 74, Focus: "SEO Content", Category: "Marketing", Finding: "Keywords, meta, blog outlines", HasDetail: true, SourceFile: "exp-batch-67-72-74.py",
				Why:     `Paid advertising stops working the moment you stop paying. SEO content is the long-term acquisition strategy. We needed a complete SEO plan: keyword research, content calendar, on-page optimisation, technical SEO requirements, and a link building strategy.`,
				What:    `A comprehensive SEO strategy: 20 keywords with estimated search volume and difficulty, a 3-month content calendar (12 blog posts), on-page SEO templates (title tags, meta descriptions, schema markup), technical SEO checklist (site speed, mobile-first, sitemaps, canonical URLs), and 5 link building strategies with specific targets.`,
				How:     `A single AI call produced all components. The keywords mixed informational ("how to track freelance clients") and transactional ("best CRM for freelancers") intent. The content calendar scheduled one post per week with content types varying between how-to, comparison, listicle, and guide formats.`,
				Impact:  `The keyword research produced realistic targets with appropriate difficulty ratings. The content calendar was immediately actionable &mdash; each post had a title, target keyword, word count, and content type. For a startup with no marketing team, this is a complete first-quarter SEO plan for under $0.01.`,
				Related: []string{"72", "70"},
			},
			{Num: "75", NumID: 75, Focus: "Product Hunt", Category: "Marketing", Finding: "Launch strategy, maker comment, assets", HasDetail: true, SourceFile: "exp-batch-73-75-79-80.py",
				Why:     `Product Hunt is the canonical launch platform for developer tools and SaaS products. A good launch can drive thousands of signups in a day. We needed everything prepared in advance: tagline, description, maker comment, launch timing, engagement plan, and post-launch conversion strategy.`,
				What:    `A complete Product Hunt launch package: tagline (under 60 chars), description (under 260 chars), an authentic 200-word maker comment, 5 topic tags, launch checklist (best day/time, early upvote strategy, comment engagement templates), pre-launch social media teasers, and post-launch follow-up plan.`,
				How:     `A single AI call with explicit character limits and formatting requirements. The prompt required the maker comment to be "authentic and personal" rather than marketing-speak. The launch checklist included both ethical upvote strategies and templates for responding to different comment types.`,
				Impact:  `The maker comment read like a genuine founder story, not AI-generated copy. The launch checklist covered details that first-time launchers miss: timing (Tuesday-Thursday, 00:01 PST), hunter outreach, and converting PH visitors to signups with a dedicated landing page variant. This is a complete launch playbook for under $0.01.`,
				Related: []string{"70", "73"},
			},
			{Num: "76", NumID: 76, Focus: "Distribution", Category: "Marketing", Finding: "Channel analysis, partnership targets", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `Products do not sell themselves. We needed a distribution strategy that ranks channels by effort and impact, identifies partnership opportunities, and designs a referral programme &mdash; the complete playbook for getting from zero users to traction.`,
				What:    `Ten distribution channels ranked by effort/impact (Reddit, Twitter, LinkedIn, Product Hunt, Indie Hackers, podcasts, partnerships, SEO, referrals, cold outreach), with expected results, cost, and timeline for each. Plus 5 integration partnership targets (Stripe, Slack, Zapier, Notion, QuickBooks) and a referral programme design.`,
				How:     `A single AI call produced the full distribution plan. Each channel included a specific how-to, realistic expected results, cost estimate, and time-to-results. The partnership targets were chosen for ecosystem fit with a freelancer CRM.`,
				Impact:  `The channel ranking was sensible: Reddit and Indie Hackers ranked highest for effort-to-impact ratio, while cold outreach ranked lowest. The referral programme design (invite-3-get-a-month-free) was simple and implementable. Together with the other marketing experiments, we now have a complete go-to-market stack.`,
				Related: []string{"70", "75"},
			},
			{Num: "77", NumID: 77, Focus: "Analytics", Category: "Marketing", Finding: "Event tracking, funnel definition", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `You cannot improve what you do not measure. Before launching, we need to define what metrics matter, what events to track, what constitutes "activation," and what alert thresholds signal trouble. We asked the AI to design the analytics layer for our CRM.`,
				What:    template.HTML(`An analytics design covering the signup funnel (visit, signup, activate, pay), activation definition (what counts as "activated"), retention metrics (DAU/WAU/MAU), revenue metrics (MRR, ARPU, churn, LTV), feature usage tracking, a dashboard mockup, and alert thresholds for churn spikes and signup drops.`),
				How:     `A single AI call with a structured prompt covering each analytics domain. The output defined specific events to track, funnel stages with expected conversion rates, and alert thresholds with recommended responses.`,
				Impact:  `The activation definition was the most valuable output &mdash; "a user who has added at least one client and created one invoice within 7 days." This is specific, measurable, and directly tied to retention. The alert thresholds (churn above 5% monthly, signups down 30% week-over-week) give early warning of problems.`,
				Related: []string{"78"},
			},
			{Num: "78", NumID: 78, Focus: "Revenue Model", Category: "Marketing", Finding: "MRR projections, churn modeling", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `Pricing tiers from <a href="/exp/71">Experiment 71</a> define what we charge. Revenue modelling defines what we can expect to earn. We needed projections that account for churn, conversion rates, expansion revenue, and the leaky-bucket reality of SaaS economics.`,
				What:    `A revenue optimisation plan for a post-100-user SaaS: conversion funnel analysis (where free users drop off), upgrade triggers (what makes someone go Pro), churn prevention (warning signs and intervention playbook), upsell strategy (Pro to Business), pricing experiments (A/B test ideas), and expansion revenue opportunities.`,
				How:     `A single AI call with a prompt focused on post-launch revenue mechanics rather than pre-launch pricing. The output covered the full revenue lifecycle from acquisition through expansion to churn prevention.`,
				Impact:  `The churn prevention playbook was the standout: it identified warning signs (login frequency dropping, no new clients added in 2 weeks) and matched them to specific interventions (email nudge, in-app prompt, personal outreach). Combined with <a href="/exp/77">Experiment 77</a>'s analytics, we have both the measurement and the response plan.`,
				Related: []string{"71", "77"},
			},
			{Num: "79", NumID: 79, Focus: "Support", Category: "Operations", Finding: "10 help articles, 20 chatbot responses", HasDetail: true, SourceFile: "exp-batch-73-75-79-80.py",
				Why:     `Users will have questions. Without support infrastructure, every question becomes a founder interruption. We wanted to pre-generate the most common support content: help articles, chatbot responses, ticket categories, escalation rules, and response templates.`,
				What:    `10 help centre articles (150 words each covering getting started through cancellation), 20 chatbot responses with trigger phrases and escalation conditions, 5 ticket categories, a 3-tier escalation structure (bot, human, engineering), and 5 response templates for common situations.`,
				How:     template.HTML(`A single AI call produced all support artifacts. Each chatbot response included the trigger phrase, the bot's reply, and the condition that triggers escalation to a human (e.g., user says "it's not working" after receiving the standard answer).`),
				Impact:  `The help articles were specific enough to use as-is. The escalation rules were sensible: billing and account access go directly to L2 (human), data issues and security concerns go to L3 (engineering). For a product with no support team, this is enough to handle the first 100 users.`,
				Related: []string{"67"},
			},
			{Num: "80", NumID: 80, Focus: "Legal", Category: "Operations", Finding: "ToS, Privacy, Cookie, AUP", HasDetail: true, SourceFile: "exp-batch-73-75-79-80.py",
				Why:     `You cannot launch a SaaS product without legal documents. Terms of Service, Privacy Policy, Cookie Policy, and Acceptable Use Policy are table stakes. We wanted to test whether AI-generated legal templates would be specific enough to our product to serve as a starting point (with lawyer review before use).`,
				What:    `Four legal documents tailored to an EU-hosted SaaS CRM: Terms of Service (acceptable use, payment via Stripe, liability limitations), GDPR-compliant Privacy Policy (data collected, legal basis, retention, user rights, third-party processors), Cookie Policy (essential cookies only for MVP), and Acceptable Use Policy.`,
				How:     `A single AI call received specific product details: what data is stored (client names, emails, invoices), hosting location (EU), payment processor (Stripe), and pricing ($15/month). The prompt explicitly noted these are templates requiring lawyer review.`,
				Impact:  `The Privacy Policy correctly referenced GDPR articles, listed specific user rights (access, rectify, delete, port), and named Stripe and the hosting provider as third-party processors. The documents are not legal advice, but they are far more useful than generic templates because they reference our specific product and data flows.`,
			},
			{Num: "81", NumID: 81, Focus: "Localisation", Category: "Operations", Finding: "i18n extraction, translation", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `English-only products leave money on the table. We needed a localisation plan that covers which languages to prioritise, what to translate (and what not to), how to architect i18n in Go/TypeScript, and how to handle locale-specific formatting for dates, currencies, and numbers.`,
				What:    template.HTML(`A localisation plan covering the top 5 priority languages by market size, what to translate (UI strings, help docs, legal docs, marketing) versus what to skip (code, API responses, internal logs), i18n architecture for Go and TypeScript, date/currency/number formatting per locale, and RTL support assessment.`),
				How:     template.HTML(`A single AI call with a structured prompt. The output prioritised languages by addressable market size and estimated effort per language, including the distinction between translating UI strings (cheap) and legal documents (requires legal review per jurisdiction).`),
				Impact:  `The architecture guidance was practical: use Go's message catalogues for UI strings, keep translation keys in a JSON file per locale, and never translate API responses. The RTL assessment correctly flagged Arabic as requiring layout changes beyond just text translation. Estimated effort: 2-3 days per language for UI, 1-2 weeks for legal.`,
			},
			{Num: "82", NumID: 82, Focus: "Competitor Monitor", Category: "Marketing", Finding: "Change tracking, alert rules", HasDetail: true, SourceFile: "exp-batch-70-71-76-77-78-81-82.py",
				Why:     `Competitors do not stand still. We needed a monitoring system that tracks pricing changes, new features, review sentiment, and traffic trends across our top competitors &mdash; and alerts us when something requires a response.`,
				What:    `A competitor monitoring system design: 5 competitors to track (with URLs and what to watch), monitoring schedule (daily/weekly/monthly), data sources (G2, Capterra, SimilarWeb, Twitter, changelog pages), alert triggers (price drop, new feature we lack), and a response playbook for different competitive moves.`,
				How:     `A single AI call produced the complete monitoring design. Each competitor got a specific watch list: pricing page URL, changelog/blog URL, review profiles on G2 and Capterra. Alert triggers were tied to specific responses (e.g., "competitor drops price 20%" triggers pricing review meeting).`,
				Impact:  `The response playbook was the most useful output. Rather than reacting emotionally to competitive moves, it defines specific responses: competitor launches a feature we have &mdash; write a comparison blog post. Competitor drops price &mdash; emphasise value rather than matching. This turns competitive intelligence from anxiety into a process.`,
				Related: []string{"69"},
			},
			{Num: "83", NumID: 83, Focus: "Org Personas", Category: "Marketing", Finding: "B2B buyer personas, org charts", HasDetail: true, SourceFile: "exp83-org-personas.py",
				Why:     `<a href="/exp/23">Experiment 23</a> created consumer personas (freelancer, agency owner). But within a single organisation, a clerk, an operator, a manager, and an owner need completely different things from the same product. We wanted to test whether role-based personas would surface requirements that consumer personas miss.`,
				What:    template.HTML(`Five org-level personas reviewed the CRM: data entry clerk (speed, bulk import), customer (invoice self-service), office manager (cross-client views, dashboards), team lead (team activity, performance metrics), and business owner (P&amp;L, cash flow, tax summaries). Each rated the CRM 1-10 for their specific use case.`),
				How:     `Five sequential AI calls, each receiving the CRM feature list and a role-specific background. Each persona answered in character: does this serve my needs, what is missing, what would I use most, what is useless to me, and would I recommend it. Ratings were extracted programmatically.`,
				Impact:  `The ratings revealed a clear pattern: the CRM scored well for day-to-day operators but poorly for clerks (no bulk import) and owners (no financial reporting). Different org levels need completely different features from the same product. This validates that B2B products need role-based persona reviews, not just user-type personas.`,
				Related: []string{"23"},
			},
			{Num: "84", NumID: 84, Focus: "Enriched Docs", Category: "Content", Finding: "Context-aware documentation", HasDetail: true, SourceFile: "exp-batch-58-62-84.py",
				Why:     `After 84 experiments, our documentation was inconsistent &mdash; some experiments had detailed READMEs, others had one-liners. We needed a standardised template that every experiment document should follow, plus a summary table covering all 84 experiments.`,
				What:    `A standardised experiment documentation template (What, Why, How, Results table, Key Finding, Impact on Pipeline, Files list) and a summary table of all 84 experiments in a consistent format: number, what, why, result, cost, key finding.`,
				How:     `A single AI call generated both the template and the summary. The template enforced structure: every experiment must explain what it tests, why it matters, how it was run, what it found, and what changed in the pipeline because of it.`,
				Impact:  `The template became the format we use for every experiment page on this site. The summary table provided the first complete inventory of all 84 experiments in one view. This experiment is meta &mdash; it improved the documentation of the research itself, making the findings more accessible and navigable.`,
				Related: []string{"63", "67"},
			},
		},
	},
	{
		Slug:     "graphql",
		Name:     "GraphQL",
		Icon:     "\u26a1",
		ExpRange: "Experiments 91 – 93",
		Narrative: []template.HTML{
			`GraphQL introduces a unique tension for AI code generation: the traditional wisdom is to use code generation tools (like gqlgen for Go) to ensure the schema and resolvers stay in sync. But if AI is generating all the code anyway, is the codegen layer redundant overhead? These experiments test three approaches: code-first (define Go structs, generate schema), schema-first (write .graphql files, generate Go), and hand-rolled (AI writes both schema and resolvers from scratch).`,
			`The hypothesis is that codegen adds complexity without proportional value when the "developer" is an AI model that can maintain consistency across files. A human needs codegen because keeping types in sync across schema.graphql and 15 resolver files is tedious and error-prone. An AI model processes all files in a single context window. If the hand-rolled approach produces correct, type-safe code reliably, it would simplify the build pipeline significantly by eliminating the codegen dependency.`,
			`Results from Experiment 91 are clear: <strong>minimal hand-rolled GraphQL wins for AI generation.</strong> The minimal approach (no framework, 231 lines, $0.045, 5 AI turns) was the only one where the server started and responded to queries correctly. graphql-go built but crashed on startup. gqlgen built but required 34 AI turns and $0.22 &mdash; nearly 5x more expensive &mdash; and had endpoint routing issues at runtime. Complexity directly correlates with failure rate.`,
		},
		KeyInsight: `Minimal (no framework) wins: 231 lines, $0.045, only approach that worked end-to-end. gqlgen's codegen is redundant and 5x more expensive when AI generates code. For the pipeline template (built once by Opus), graphql-go is acceptable; for per-project generation (by Haiku), hand-rolled is most reliable.`,
		TableType:  "graphql",
		Experiments: []Experiment{
			{Num: "91a", NumID: 91, Focus: "graphql-go (code-first)", Result: "BUILD PASS, runtime fail", Cost: "$0.07", Finding: "329 lines, 1 file, 6 AI turns — server crashed on start", HasDetail: true, Icon: "\u26a1",
				Why:     `Dark Factory uses gqlgen for GraphQL — a schema-first tool that generates Go code from .graphql files. But when AI generates all the code, the codegen step is redundant. The AI writes the schema, runs <code>go generate</code>, then fills in resolvers that conflict with what was generated. It's a human workflow being forced on a non-human developer. We wanted to know: is there a simpler approach that works better for AI-generated code?`,
				What:    `Three approaches tested head-to-head, each building the same CRM GraphQL API (clients, activities, invoices) with Haiku:<br><br><strong>graphql-go</strong> (this experiment): Code-first using the graphql-go/graphql library. Schema defined in Go code with <code>graphql.NewObject()</code>. No schema file, no codegen. Single main.go.<br><br><strong>gqlgen</strong>: Schema-first. Write .graphql file, run <code>go generate</code>, fill in resolvers. 7 files, codegen step.<br><br><strong>Minimal</strong>: No library at all. Hand-rolled query parsing with regex, route to handler functions, return JSON. Single main.go.`,
				How:     `Each approach got the same prompt and the same CRM spec. We measured: AI turns needed, cost, lines of code, number of files, and whether the server actually starts and responds to queries. The gate wasn't just <code>go build</code> — we sent real GraphQL queries and checked the responses.<br><br>graphql-go: 6 AI turns, $0.07, 329 lines, 1 file. gqlgen: 34 AI turns, $0.22, 7 files. Minimal: 5 AI turns, $0.045, 231 lines, 1 file.`,
				Impact:  `<strong>Minimal won decisively.</strong> It was the only approach where the server actually started and responded to queries correctly. graphql-go compiled but crashed on startup (port binding issues). gqlgen compiled but had endpoint routing problems — and cost 5x more with 7x more AI turns.<br><br>The lesson: <strong>complexity correlates directly with failure rate</strong> when AI generates code. Fewer files, fewer dependencies, fewer build steps = more reliable output. For the pipeline template (built once by Opus), graphql-go is acceptable. For per-project generation (by Haiku), hand-rolled is most reliable. gqlgen should not be used in AI pipelines — its codegen step conflicts with AI generation.`,
				Related: []string{"29", "27"},
			},
			{Num: "91b", NumID: 0, Focus: "gqlgen (schema-first)", Result: "BUILD PASS, runtime fail", Cost: "$0.22", Finding: "317 lines + 34 generated, 7 files, 34 AI turns — endpoint routing issues"},
			{Num: "91c", NumID: 0, Focus: "Minimal (hand-rolled)", Result: "BUILD PASS, queries work", Cost: "$0.045", Finding: "231 lines, 1 file, 5 AI turns — only approach that worked end-to-end"},
		},
	},
	{
		Slug:     "model-comparison",
		Name:     "Model Comparison",
		Icon:     "\U0001f916",
		ExpRange: "9 Models Across Multiple Tasks",
		Narrative: []template.HTML{
			`Across the full research programme, 11 models were tested on various tasks ranging from simple bash scripts to multi-service architectures to website cloning. The comparison reveals that no single model dominates all tasks. Opus produces the highest visual quality for website cloning but costs 100x more than open-source alternatives. MiniMax M2.7 has the best instruction-following for code generation at $0.007/call. Gemini 2.5 Flash is the best planner. The optimal strategy is model routing: use the right model for each sub-task rather than one model for everything.`,
			`The Spike V1 experiments tested 11 models on a simple bash script task (5 tests). All 11 passed, costing between $0.008 and $0.015 each. At this complexity level, model choice genuinely does not matter. The differentiation emerges at higher complexity: the Spike V2 dep-doctor application (10 files, 18 tests) separated capable models from those that struggle with multi-file consistency. The winning configuration (A3: Gemini Flash planner + MiniMax M2.7 executor) achieved 18/18 tests at $0.069 &mdash; a complete CLI application for less than 7 cents.`,
			`For website cloning, the visual quality gap is stark. Opus and Sonnet produce professional-looking clones with correct typography, spacing, and colour schemes. Open-source models (Llama 4 Scout, Devstral) produce structurally correct but visually rougher output. The cheapest models (Qwen3-30B, Gemma 3 27B) struggle with CSS and server architecture. The cost-quality frontier has a clear knee: Llama 4 Scout at $0.02 delivers the best return on investment for most use cases.`,
		},
		KeyInsight: `Model routing beats single-model approaches. Use Gemini Flash for planning, MiniMax M2.7 for code execution, and Opus only when visual quality matters.`,
		TableType:  "model-comparison",
		Experiments: []Experiment{
			{Num: "", NumID: 0, Focus: "Opus", Result: "8/8", Cost: "$0.30+", Finding: "Best overall", Category: "0.708"},
			{Num: "", NumID: 0, Focus: "Sonnet", Result: "90% (A1)", Cost: "$0.012", Finding: "Very good", Category: "0.669"},
			{Num: "", NumID: 0, Focus: "Gemini 2.5 Flash", Result: "N/A (planner)", Cost: "$0.013", Finding: "—", Category: "—"},
			{Num: "", NumID: 0, Focus: "MiniMax M2.7", Result: "100% (A3)", Cost: "$0.007", Finding: "—", Category: "—"},
			{Num: "", NumID: 0, Focus: "Llama 4 Scout", Result: "8/8", Cost: "~$0.002", Finding: "Good — best value", Category: "0.722"},
			{Num: "", NumID: 0, Focus: "Devstral Small", Result: "8/8", Cost: "~$0.002", Finding: "Wireframe style", Category: "0.732"},
			{Num: "", NumID: 0, Focus: "Haiku", Result: "8/8", Cost: "$0.018", Finding: "Good", Category: "0.649"},
			{Num: "", NumID: 0, Focus: "Qwen3 Coder (480B)", Result: "80% (A2)", Cost: "$0.008", Finding: "—", Category: "—"},
			{Num: "", NumID: 0, Focus: "Qwen3-30B", Result: "70–88%", Cost: "$0.0005", Finding: "Broken CSS", Category: "0.745"},
		},
	},
	{
		Slug:     "composable-architecture",
		Name:     "Composable Architecture",
		Icon:     "\U0001f9e9",
		ExpRange: "The Self-Bootstrapping Pipeline System",
		Narrative: []template.HTML{
			`The research revealed that the 13-step pipeline is really a system of <strong>composable blocks</strong>. Each experiment is a block configuration we already tested. Experiment 23 is the persona-discovery block. Experiment 32 is the domain-review block. Experiment 38 is the progressive-enhancement block. The 116+ experiments unknowingly built a tested block library.`,
			`The architecture has three layers: <strong>Block definitions</strong> in YAML/JSON describe what each step does &mdash; its prompt template, model, gate condition, and failure handling. A <strong>Go orchestrator</strong> reads the pipeline definition and executes blocks, managing state, routing to models, and handling retries and escalation. A <strong>web GUI</strong> provides a visual node-based editor to wire blocks together, configure them, run pipelines, and inspect results at each step.`,
			`Blocks aren't just "call a model with a prompt." Some blocks are <strong>autonomous agents with tools</strong>. The pen test agent (Exp 42) uses curl and HTTP requests to actively attack the application. The Playwright agent (Exp 40) controls a real browser. The chaos agent (Exp 45) fires concurrent malformed requests. The feedback loop agent (Exp 56) reads error output and writes code fixes. Each agent block has a tool list, a persona, and a max number of autonomous turns.`,
			`The feedback loop is a <strong>meta-block</strong> &mdash; a wrapper that adds retry logic around any other block. It takes any block and adds: run &rarr; check gate &rarr; if fail &rarr; diagnose &rarr; fix &rarr; retry. You can wrap the build block with it, or the review block, or the test block. Similarly, model escalation is a block feature: try haiku ($0.005), if the gate fails escalate to sonnet ($0.03), then opus ($0.15).`,
			`<strong>The biggest insight: the system can bootstrap itself.</strong> You define a "build a pipeline block" pipeline as... a pipeline. The retro agent identifies gaps ("we kept missing auth issues"). The block-builder agent creates the YAML definition and prompt template for a new auth-review block. The test agent validates it works on a test input. The wiring agent inserts it into the pipeline at the right position. The system literally builds its own blocks and gets smarter over time.`,
		},
		KeyInsight: `116+ experiments = a tested block library. The system builds its own blocks. Retro agent finds gaps &rarr; block-builder creates YAML + prompt &rarr; test agent validates &rarr; wiring agent inserts. The pipeline that builds blocks is itself a pipeline of blocks.`,
		TableType:  "standard",
		Experiments: []Experiment{
			{Num: "Layer 1", NumID: 0, Focus: "Block Definitions (YAML)", Result: "Design", Cost: "—", Finding: "Each step: prompt template, model, gate, on_fail action"},
			{Num: "Layer 2", NumID: 0, Focus: "Orchestrator (Go)", Result: "Design", Cost: "—", Finding: "Executes blocks, manages state, routes to models"},
			{Num: "Layer 3", NumID: 0, Focus: "GUI Editor (Web)", Result: "Design", Cost: "—", Finding: "Visual node editor to wire, configure, run, inspect"},
			{Num: "Block", NumID: 0, Focus: "Prompt Block", Result: "Proven", Cost: "~$0.01", Finding: "Simple model call with template — most pipeline steps"},
			{Num: "Block", NumID: 0, Focus: "Agent Block", Result: "Proven", Cost: "~$0.03", Finding: "Autonomous agent with tools (Playwright, curl, chaos)"},
			{Num: "Block", NumID: 0, Focus: "Gate Block", Result: "Proven", Cost: "FREE", Finding: "Pass/fail check: go build, test rate, reviewer score"},
			{Num: "Block", NumID: 0, Focus: "Wrapper/Meta Block", Result: "Proven", Cost: "—", Finding: "Adds retry + feedback loop around any other block"},
			{Num: "Block", NumID: 0, Focus: "Parallel Block", Result: "Proven", Cost: "—", Finding: "Runs independent branches concurrently"},
			{Num: "Meta", NumID: 0, Focus: "Self-Bootstrap Pipeline", Result: "Design", Cost: "—", Finding: "Retro → block-builder → test → wire-in (recursive)"},
		},
	},
	{
		Slug:     "planned",
		Name:     "Planned Experiments",
		Icon:     "\U0001f4cb",
		ExpRange: "Experiments 99 – 127",
		Narrative: []template.HTML{
			"The research continues. These experiments are designed but not yet run...",
		},
		KeyInsight: "127 experiments planned total — from code generation to self-bootstrapping pipelines.",
		TableType:  "business-pipeline",
		Experiments: []Experiment{
			// Composable Pipeline (99-101)
			{Num: "99", NumID: 99, Focus: "Block Interface Definition", Result: "Planned", Cost: "—", Finding: "Define Block struct, convert exp27 to YAML", Category: "Composable Pipeline", HasDetail: false},
			{Num: "100", NumID: 100, Focus: "Visual Block Editor", Result: "Planned", Cost: "—", Finding: "Web GUI for wiring blocks", Category: "Composable Pipeline", HasDetail: false},
			{Num: "101", NumID: 101, Focus: "Self-Bootstrapping", Result: "Planned", Cost: "—", Finding: "System builds its own blocks", Category: "Composable Pipeline", HasDetail: false},
			// Code Quality (102-105)
			{Num: "102", NumID: 102, Focus: "Duplicate Code Detection", Result: "Planned", Cost: "—", Finding: "dupl + goconst feedback loop", Category: "Code Quality", HasDetail: false},
			{Num: "103", NumID: 103, Focus: "Extend vs Create Agent", Result: "Planned", Cost: "—", Finding: "Check if feature can extend existing code", Category: "Code Quality", HasDetail: false},
			{Num: "104", NumID: 104, Focus: "Dependency Audit", Result: "Planned", Cost: "—", Finding: "Remove unnecessary deps", Category: "Code Quality", HasDetail: false},
			{Num: "105", NumID: 105, Focus: "API Contract Testing", Result: "Planned", Cost: "—", Finding: "OpenAPI spec verification", Category: "Code Quality", HasDetail: false},
			// Multi-Modal (106-108)
			{Num: "106", NumID: 106, Focus: "Screenshot-Guided Iteration", Result: "Planned", Cost: "—", Finding: "Visual feedback vs text", Category: "Multi-Modal", HasDetail: false},
			{Num: "107", NumID: 107, Focus: "Voice-to-Brief", Result: "Planned", Cost: "—", Finding: "Whisper transcription to pipeline", Category: "Multi-Modal", HasDetail: false},
			{Num: "108", NumID: 108, Focus: "Video Walkthrough", Result: "Planned", Cost: "—", Finding: "Auto-generate demo videos", Category: "Multi-Modal", HasDetail: false},
			// Performance (109-111)
			{Num: "109", NumID: 109, Focus: "Load Testing Agent", Result: "Planned", Cost: "—", Finding: "k6/vegeta scripts from API spec", Category: "Performance", HasDetail: false},
			{Num: "110", NumID: 110, Focus: "Cost Optimization Agent", Result: "Planned", Cost: "—", Finding: "Find cheaper pipeline steps", Category: "Performance", HasDetail: false},
			{Num: "111", NumID: 111, Focus: "Parallel Pipeline", Result: "Planned", Cost: "—", Finding: "Run independent steps concurrently", Category: "Performance", HasDetail: false},
			// User Experience (112-114)
			{Num: "112", NumID: 112, Focus: "A/B Testing Agent", Result: "Planned", Cost: "—", Finding: "Persona panel votes on variants", Category: "User Experience", HasDetail: false},
			{Num: "113", NumID: 113, Focus: "Accessibility Automation", Result: "Planned", Cost: "—", Finding: "axe-core + auto-fix", Category: "User Experience", HasDetail: false},
			{Num: "114", NumID: 114, Focus: "Internationalization", Result: "Planned", Cost: "—", Finding: "i18n extraction + translation", Category: "User Experience", HasDetail: false},
			// Self-Improvement (115-117)
			{Num: "115", NumID: 115, Focus: "Failure Pattern Learning", Result: "Planned", Cost: "—", Finding: "Knowledge base from 90+ experiments", Category: "Self-Improvement", HasDetail: false},
			{Num: "116", NumID: 116, Focus: "Prompt Evolution", Result: "Planned", Cost: "—", Finding: "Evolutionary optimization of prompts", Category: "Self-Improvement", HasDetail: false},
			{Num: "117", NumID: 117, Focus: "Meta-Pipeline", Result: "Planned", Cost: "—", Finding: "The system improves itself run over run", Category: "Self-Improvement", HasDetail: false},
			// Design System (118-122)
			{Num: "118", NumID: 118, Focus: "Token Extraction", Result: "DONE", Cost: "—", Finding: "20 sites analyzed", Category: "Design System", HasDetail: false},
			{Num: "119", NumID: 119, Focus: "Component Analysis", Result: "DONE", Cost: "—", Finding: "10 MUST components", Category: "Design System", HasDetail: false},
			{Num: "120", NumID: 120, Focus: "Build Framework", Result: "Planned", Cost: "—", Finding: "From 20-site evidence", Category: "Design System", HasDetail: false},
			{Num: "121", NumID: 121, Focus: "Token Swap Validation", Result: "Planned", Cost: "—", Finding: "Does framework generalize?", Category: "Design System", HasDetail: false},
			{Num: "122", NumID: 122, Focus: "Auto-Decompose", Result: "Planned", Cost: "—", Finding: "Given a clone, extract framework + tokens", Category: "Design System", HasDetail: false},
			// Layout & Design (123-127)
			{Num: "123", NumID: 123, Focus: "Layout Guidelines Extraction", Result: "Planned", Cost: "—", Finding: "Layout rules, not just tokens", Category: "Layout & Design", HasDetail: false},
			{Num: "124", NumID: 124, Focus: "Cross-Pollination", Result: "Planned", Cost: "—", Finding: "Mix layout from one site with tokens from another", Category: "Layout & Design", HasDetail: false},
			{Num: "125", NumID: 125, Focus: "Responsive Breakpoint Analysis", Result: "Planned", Cost: "—", Finding: "Universal breakpoint patterns", Category: "Layout & Design", HasDetail: false},
			{Num: "126", NumID: 126, Focus: "Animation Patterns", Result: "Planned", Cost: "—", Finding: "Hover, transition, scroll animations", Category: "Layout & Design", HasDetail: false},
			{Num: "127", NumID: 127, Focus: "Accessibility Audit", Result: "Planned", Cost: "—", Finding: "axe-core across all 26 clones", Category: "Layout & Design", HasDetail: false},
		},
	},
}

// Build lookup maps
var (
	categoryBySlug = map[string]*Category{}
	expByNum       = map[int]*Experiment{}
	expByStr       = map[string]*Experiment{}  // string-based lookup for sub-experiments like "172b"
	expCatByStr    = map[string]*Category{}    // category lookup by string key
	expCategory    = map[int]*Category{}
	sortedExpIDs   []int
	discoveryGraph template.HTML
)

var costRe = regexp.MustCompile(`\$?([\d.]+)`)

func parseCost(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "—" || s == "-" {
		return -1
	}
	if strings.EqualFold(s, "FREE") {
		return 0
	}
	m := costRe.FindStringSubmatch(s)
	if len(m) >= 2 {
		v, _ := strconv.ParseFloat(m[1], 64)
		return v
	}
	return -1
}

func countWords(parts ...template.HTML) int {
	n := 0
	for _, p := range parts {
		s := strings.TrimSpace(string(p))
		if s != "" {
			n += len(strings.Fields(s))
		}
	}
	return n
}

func init() {
	seen := map[int]bool{}
	for i := range categories {
		cat := &categories[i]
		categoryBySlug[cat.Slug] = cat
		for j := range cat.Experiments {
			exp := &cat.Experiments[j]
			// Parse cost
			exp.CostFloat = parseCost(exp.Cost)
			if exp.CostFloat > cat.MaxCost {
				cat.MaxCost = exp.CostFloat
			}
			// Detect failure
			r := strings.ToUpper(exp.Result)
			exp.IsFailure = strings.Contains(r, "FAIL") || strings.Contains(r, "0/") || r == "0%"
			// Word count + reading time
			exp.WordCount = countWords(exp.Why, exp.What, exp.How, exp.Impact)
			for _, d := range exp.Description {
				exp.WordCount += countWords(d)
			}
			if exp.WordCount > 0 {
				exp.ReadingTime = int(math.Ceil(float64(exp.WordCount) / 200.0))
				if exp.ReadingTime < 1 {
					exp.ReadingTime = 1
				}
			}
			// Parse SSIM score for cloning experiments
			if cat.TableType == "website-clone" || cat.TableType == "model-comparison" {
				if score, err := strconv.ParseFloat(exp.Category, 64); err == nil {
					exp.Score = score
				}
			}
			// Register in lookup maps
			if exp.HasDetail && exp.NumID > 0 {
				if !seen[exp.NumID] {
					expByNum[exp.NumID] = exp
					expCategory[exp.NumID] = cat
					sortedExpIDs = append(sortedExpIDs, exp.NumID)
					seen[exp.NumID] = true
				}
			}
			// Register string-based lookup for sub-experiments (e.g. "172b")
			if exp.HasDetail && exp.Num != strconv.Itoa(exp.NumID) {
				expByStr[exp.Num] = exp
				expCatByStr[exp.Num] = cat
			}
		}
	}
	sort.Ints(sortedExpIDs)

	// Build discovery graph from Related links
	var sb strings.Builder
	sb.WriteString("graph LR\n")
	edges := map[string]bool{}
	for _, id := range sortedExpIDs {
		exp := expByNum[id]
		if len(exp.Related) == 0 {
			continue
		}
		for _, rel := range exp.Related {
			relID, err := strconv.Atoi(rel)
			if err != nil {
				continue
			}
			edge := fmt.Sprintf("%d->%d", id, relID)
			if edges[edge] {
				continue
			}
			edges[edge] = true
			relExp, ok := expByNum[relID]
			srcLabel := strings.ReplaceAll(exp.Focus, "\"", "'")
			if ok {
				dstLabel := strings.ReplaceAll(relExp.Focus, "\"", "'")
				sb.WriteString(fmt.Sprintf("    %d[\"%d: %s\"] --> %d[\"%d: %s\"]\n", id, id, srcLabel, relID, relID, dstLabel))
			}
		}
	}
	// Add click directives
	for id := range expByNum {
		sb.WriteString(fmt.Sprintf("    click %d \"/exp/%d\"\n", id, id))
	}
	discoveryGraph = template.HTML(sb.String())

	// Count narrative completeness
	for _, id := range sortedExpIDs {
		exp := expByNum[id]
		narrativeTotal++
		if exp.Why != "" && exp.What != "" && exp.How != "" && exp.Impact != "" {
			narrativeComplete++
		}
	}
}

var (
	narrativeComplete int
	narrativeTotal    int
)

// getExp looks up an experiment by its string ID for template use.
func getExp(idStr string) *Experiment {
	id, err := strconv.Atoi(idStr)
	if err == nil {
		if exp, ok := expByNum[id]; ok {
			return exp
		}
	}
	// Fall back to string-based lookup for sub-experiments like "172b"
	if exp, ok := expByStr[idStr]; ok {
		return exp
	}
	return nil
}

// scoreColor returns a CSS color based on SSIM score thresholds.
func scoreColor(score float64) string {
	if score >= 0.7 {
		return "#28a745"
	}
	if score >= 0.5 {
		return "#f0ad4e"
	}
	return "#dc3545"
}

// expSearchJSON returns all experiment data as a JS-safe JSON string for client-side search.
func expSearchJSON() template.JS {
	var sb strings.Builder
	sb.WriteString("[")
	first := true
	for _, cat := range categories {
		for _, exp := range cat.Experiments {
			if !exp.HasDetail || exp.NumID <= 0 {
				continue
			}
			if !first {
				sb.WriteString(",")
			}
			first = false
			focus := strings.ReplaceAll(exp.Focus, `"`, `\"`)
			result := strings.ReplaceAll(exp.Result, `"`, `\"`)
			finding := strings.ReplaceAll(exp.Finding, `"`, `\"`)
			sb.WriteString(fmt.Sprintf(`{"num":"%s","numID":%d,"focus":"%s","result":"%s","finding":"%s","icon":"%s","cat":"%s","catName":"%s","score":%.3f}`,
				exp.Num, exp.NumID, focus, result, finding, exp.Icon, cat.Slug, cat.Name, exp.Score))
		}
	}
	sb.WriteString("]")
	return template.JS(sb.String())
}

// ScoreEntry for the comparison chart.
type ScoreEntry struct {
	Name  string
	Num   int
	Score float64
	Color string
	Cost  string
}

// sortedScores returns cloning experiments sorted by SSIM score descending.
func sortedScores() []ScoreEntry {
	var entries []ScoreEntry
	for _, cat := range categories {
		if cat.TableType != "website-clone" && cat.TableType != "model-comparison" {
			continue
		}
		for _, exp := range cat.Experiments {
			if exp.Score > 0 && exp.HasDetail {
				entries = append(entries, ScoreEntry{
					Name:  exp.Focus,
					Num:   exp.NumID,
					Score: exp.Score,
					Color: scoreColor(exp.Score),
					Cost:  exp.Cost,
				})
			}
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

func prevExp(id int) int {
	for i, v := range sortedExpIDs {
		if v == id && i > 0 {
			return sortedExpIDs[i-1]
		}
	}
	return 0
}

func nextExp(id int) int {
	for i, v := range sortedExpIDs {
		if v == id && i < len(sortedExpIDs)-1 {
			return sortedExpIDs[i+1]
		}
	}
	return 0
}

func prevClone(slug string) string {
	for i, s := range cloneSites {
		if s.Slug == slug && i > 0 {
			return cloneSites[i-1].Slug
		}
	}
	return ""
}

func nextClone(slug string) string {
	for i, s := range cloneSites {
		if s.Slug == slug && i < len(cloneSites)-1 {
			return cloneSites[i+1].Slug
		}
	}
	return ""
}

func prevCloneName(slug string) string {
	for i, s := range cloneSites {
		if s.Slug == slug && i > 0 {
			return cloneSites[i-1].Name
		}
	}
	return ""
}

func nextCloneName(slug string) string {
	for i, s := range cloneSites {
		if s.Slug == slug && i < len(cloneSites)-1 {
			return cloneSites[i+1].Name
		}
	}
	return ""
}
