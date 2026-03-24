// AI Code Orchestration Research — Documentation Viewer

// File tree structure (populated from repo)
const FILES = {
    'README.md': 'README.md',
    'docs': {
        '_folder': true,
        'research-overview.md': 'docs/research-overview.md'
    },
    'experiments': {
        '_folder': true,
        'spike-v1': {
            '_folder': true,
            'REPORT.md': 'experiments/spike-v1/REPORT.md',
            'results.tsv': 'experiments/spike-v1/results.tsv'
        },
        'spike-v2': {
            '_folder': true,
            'README.md': 'experiments/spike-v2/README.md',
            'REPORT.md': 'experiments/spike-v2/REPORT.md',
            'architecture.md': 'experiments/spike-v2/architecture.md',
            'results.tsv': 'experiments/spike-v2/results.tsv',
            'call-model.py': 'experiments/spike-v2/call-model.py',
            'run-experiment.sh': 'experiments/spike-v2/run-experiment.sh',
            'run-subtasks.py': 'experiments/spike-v2/run-subtasks.py',
            'validate-gate.sh': 'experiments/spike-v2/validate-gate.sh',
            'prompts': {
                '_folder': true,
                'planner.md': 'experiments/spike-v2/prompts/planner.md',
                'executor.md': 'experiments/spike-v2/prompts/executor.md',
                'reviewer.md': 'experiments/spike-v2/prompts/reviewer.md',
                'test-writer.md': 'experiments/spike-v2/prompts/test-writer.md',
                'mutator.md': 'experiments/spike-v2/prompts/mutator.md',
                'council-review.md': 'experiments/spike-v2/prompts/council-review.md'
            },
            'approach-a': {
                '_folder': true,
                'A3-winner': {
                    '_folder': true,
                    'cli.cjs': 'experiments/spike-v2/approach-a/A3-winner/cli.cjs',
                    'lib': {
                        '_folder': true,
                        'parser.cjs': 'experiments/spike-v2/approach-a/A3-winner/lib/parser.cjs',
                        'validator.cjs': 'experiments/spike-v2/approach-a/A3-winner/lib/validator.cjs',
                        'analyzer.cjs': 'experiments/spike-v2/approach-a/A3-winner/lib/analyzer.cjs',
                        'reporter.cjs': 'experiments/spike-v2/approach-a/A3-winner/lib/reporter.cjs',
                        'config.cjs': 'experiments/spike-v2/approach-a/A3-winner/lib/config.cjs'
                    },
                    'test': {
                        '_folder': true,
                        'dep-doctor.test.cjs': 'experiments/spike-v2/approach-a/A3-winner/test/dep-doctor.test.cjs'
                    }
                }
            }
        }
    },
    'golden-master': {
        '_folder': true,
        'dep-doctor': {
            '_folder': true,
            'cli.cjs': 'golden-master/dep-doctor/cli.cjs',
            'lib': {
                '_folder': true,
                'parser.cjs': 'golden-master/dep-doctor/lib/parser.cjs',
                'validator.cjs': 'golden-master/dep-doctor/lib/validator.cjs',
                'analyzer.cjs': 'golden-master/dep-doctor/lib/analyzer.cjs',
                'reporter.cjs': 'golden-master/dep-doctor/lib/reporter.cjs',
                'config.cjs': 'golden-master/dep-doctor/lib/config.cjs'
            },
            'test': {
                '_folder': true,
                'dep-doctor.test.cjs': 'golden-master/dep-doctor/test/dep-doctor.test.cjs'
            }
        },
        'golden-outputs': {
            '_folder': true,
            'help.txt': 'golden-master/golden-outputs/help.txt',
            'scan-valid.json': 'golden-master/golden-outputs/scan-valid.json',
            'report-valid.json': 'golden-master/golden-outputs/report-valid.json',
            'test-results.txt': 'golden-master/golden-outputs/test-results.txt'
        }
    }
};

// Initialize
mermaid.initialize({ theme: 'dark', startOnLoad: false });

const treeEl = document.getElementById('tree');
const articleEl = document.getElementById('article');
const breadcrumbEl = document.getElementById('breadcrumb');

// Build tree
function buildTree(node, parentEl, path = '') {
    for (const [name, value] of Object.entries(node)) {
        if (name === '_folder') continue;

        const el = document.createElement('div');

        if (typeof value === 'object' && value._folder) {
            el.className = 'tree-item folder';
            el.textContent = name;
            el.onclick = (e) => {
                e.stopPropagation();
                const children = el.nextElementSibling;
                if (children) children.classList.toggle('collapsed');
            };
            parentEl.appendChild(el);

            const children = document.createElement('div');
            children.className = 'tree-children';
            parentEl.appendChild(children);
            buildTree(value, children, path + name + '/');
        } else {
            const ext = name.split('.').pop();
            el.className = `tree-item file ${ext}`;
            el.textContent = name;
            el.dataset.path = typeof value === 'string' ? value : path + name;
            el.onclick = (e) => {
                e.stopPropagation();
                loadFile(el.dataset.path);
                document.querySelectorAll('.tree-item.active').forEach(a => a.classList.remove('active'));
                el.classList.add('active');
            };
            parentEl.appendChild(el);
        }
    }
}

// Load and render a file
async function loadFile(path) {
    breadcrumbEl.textContent = path;
    articleEl.innerHTML = '<div class="loading">Loading...</div>';

    try {
        // Try fetching from GitHub raw (relative path)
        const basePath = window.location.pathname.includes('/site/') ? '../' : './';
        const response = await fetch(basePath + path);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const text = await response.text();

        const ext = path.split('.').pop();

        if (ext === 'md') {
            renderMarkdown(text);
        } else if (['json', 'tsv', 'txt'].includes(ext)) {
            renderCode(text, ext === 'json' ? 'json' : 'plaintext');
        } else if (['js', 'cjs'].includes(ext)) {
            renderCode(text, 'javascript');
        } else if (ext === 'py') {
            renderCode(text, 'python');
        } else if (ext === 'sh') {
            renderCode(text, 'bash');
        } else {
            renderCode(text, 'plaintext');
        }
    } catch (e) {
        articleEl.innerHTML = `<div class="loading">Failed to load ${path}: ${e.message}</div>`;
    }
}

function renderMarkdown(text) {
    // Configure marked
    marked.setOptions({
        highlight: function(code, lang) {
            if (lang && hljs.getLanguage(lang)) {
                return hljs.highlight(code, { language: lang }).value;
            }
            return hljs.highlightAuto(code).value;
        }
    });

    articleEl.className = 'markdown-body';
    articleEl.innerHTML = marked.parse(text);

    // Render Mermaid diagrams
    articleEl.querySelectorAll('code.language-mermaid').forEach((el) => {
        const pre = el.parentElement;
        const div = document.createElement('div');
        div.className = 'mermaid';
        div.textContent = el.textContent;
        pre.replaceWith(div);
    });
    mermaid.run();

    // Highlight remaining code blocks
    articleEl.querySelectorAll('pre code').forEach((el) => {
        hljs.highlightElement(el);
    });
}

function renderCode(text, language) {
    articleEl.className = 'code-viewer';
    const highlighted = language !== 'plaintext' && hljs.getLanguage(language)
        ? hljs.highlight(text, { language }).value
        : escapeHtml(text);
    articleEl.innerHTML = `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`;
}

function escapeHtml(text) {
    return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

function filterTree(query) {
    const items = document.querySelectorAll('.tree-item.file');
    query = query.toLowerCase();
    items.forEach(item => {
        const match = item.textContent.toLowerCase().includes(query) || item.dataset.path?.toLowerCase().includes(query);
        item.style.display = match || !query ? '' : 'none';
    });
}

function toggleSidebar() {
    document.getElementById('sidebar').classList.toggle('open');
}

// Init
treeEl.innerHTML = '';
buildTree(FILES, treeEl);

// Load README by default
loadFile('README.md');
