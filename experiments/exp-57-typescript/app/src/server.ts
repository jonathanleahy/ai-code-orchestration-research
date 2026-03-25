// @ts-nocheck
import { createServer, IncomingMessage, ServerResponse } from 'node:http';
import { Store } from './store.js';

const store = new Store();

function handleRequest(req: IncomingMessage, res: ServerResponse): void {
  const url = new URL(req.url || '/', `http://${req.headers.host || 'localhost'}`);
  const pathname = url.pathname;

  // CORS headers
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }

  // Health check
  if (pathname === '/health' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'ok' }));
    return;
  }

  // Dashboard
  if (pathname === '/' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    const clients = store.getClients();
    const html = `
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>CRM Dashboard</title>
        <style>
          * { margin: 0; padding: 0; box-sizing: border-box; }
          body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f5f5f5; padding: 20px; }
          .container { max-width: 1200px; margin: 0 auto; }
          h1 { color: #333; margin-bottom: 10px; }
          .subtitle { color: #666; margin-bottom: 30px; }
          .card { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
          form { display: flex; gap: 10px; flex-wrap: wrap; }
          input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
          button { padding: 8px 16px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; font-weight: 500; }
          button:hover { background: #0056b3; }
          table { width: 100%; border-collapse: collapse; }
          th { background: #f9f9f9; padding: 12px; text-align: left; font-weight: 600; border-bottom: 2px solid #ddd; }
          td { padding: 12px; border-bottom: 1px solid #eee; }
          tr:hover { background: #fafafa; }
          .empty { color: #999; font-style: italic; }
          .status { display: inline-block; padding: 4px 12px; border-radius: 12px; background: #e8f4f8; color: #0066cc; font-size: 12px; }
        </style>
      </head>
      <body>
        <div class="container">
          <h1>CRM Dashboard</h1>
          <p class="subtitle">Manage your clients and invoices</p>

          <div class="card">
            <h2>Add New Client</h2>
            <form id="addClientForm">
              <input type="text" id="name" placeholder="Full Name" required>
              <input type="email" id="email" placeholder="Email Address" required>
              <input type="tel" id="phone" placeholder="Phone Number" required>
              <button type="submit">Add Client</button>
            </form>
          </div>

          <div class="card">
            <h2>Clients <span class="status">${clients.length}</span></h2>
            ${clients.length === 0
              ? '<p class="empty">No clients yet. Add one above to get started.</p>'
              : `
              <table>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Phone</th>
                    <th>Created</th>
                  </tr>
                </thead>
                <tbody>
                  ${clients
                    .map(
                      c => `
                    <tr>
                      <td>${c.name}</td>
                      <td>${c.email}</td>
                      <td>${c.phone}</td>
                      <td>${new Date(c.createdAt).toLocaleDateString()}</td>
                    </tr>
                  `
                    )
                    .join('')}
                </tbody>
              </table>
            `
            }
          </div>
        </div>

        <script>
          document.getElementById('addClientForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const name = document.getElementById('name').value;
            const email = document.getElementById('email').value;
            const phone = document.getElementById('phone').value;

            const response = await fetch('/api/clients', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ name, email, phone })
            });

            if (response.ok) {
              location.reload();
            } else {
              alert('Failed to add client');
            }
          });
        </script>
      </body>
      </html>
    `;
    res.end(html);
    return;
  }

  // API: List clients
  if (pathname === '/api/clients' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    const clients = store.getClients();
    res.end(JSON.stringify(clients));
    return;
  }

  // API: Create client
  if (pathname === '/api/clients' && req.method === 'POST') {
    let body = '';

    req.on('data', (chunk: any) => {
      body += chunk.toString();
    });

    req.on('end', () => {
      try {
        const { name, email, phone } = JSON.parse(body) as {
          name: unknown;
          email: unknown;
          phone: unknown;
        };

        if (typeof name !== 'string' || typeof email !== 'string' || typeof phone !== 'string') {
          res.writeHead(400, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Invalid input: name, email, and phone must be strings' }));
          return;
        }

        const client = store.createClient(name, email, phone);
        res.writeHead(201, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(client));
      } catch (error) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Invalid JSON' }));
      }
    });
    return;
  }

  // 404
  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ error: 'Not found' }));
}

const PORT = 8080;
const server = createServer(handleRequest);

server.listen(PORT, () => {
  console.log(`CRM Server running at http://localhost:${PORT}`);
});
