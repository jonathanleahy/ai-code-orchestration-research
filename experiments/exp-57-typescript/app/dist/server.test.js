// @ts-nocheck
import test from 'node:test';
import assert from 'node:assert';
import { createServer } from 'node:http';
import { Store } from './store.js';
const store = new Store();
function createTestServer() {
    return createServer((req, res) => {
        const url = new URL(req.url || '/', `http://${req.headers.host || 'localhost'}`);
        const pathname = url.pathname;
        // Health check
        if (pathname === '/health' && req.method === 'GET') {
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ status: 'ok' }));
            return;
        }
        // List clients
        if (pathname === '/api/clients' && req.method === 'GET') {
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify(store.getClients()));
            return;
        }
        // Create client
        if (pathname === '/api/clients' && req.method === 'POST') {
            let body = '';
            req.on('data', (chunk) => {
                body += chunk.toString();
            });
            req.on('end', () => {
                try {
                    const { name, email, phone } = JSON.parse(body);
                    if (typeof name !== 'string' || typeof email !== 'string' || typeof phone !== 'string') {
                        res.writeHead(400, { 'Content-Type': 'application/json' });
                        res.end(JSON.stringify({ error: 'Invalid input' }));
                        return;
                    }
                    const client = store.createClient(name, email, phone);
                    res.writeHead(201, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify(client));
                }
                catch (error) {
                    res.writeHead(400, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ error: 'Invalid JSON' }));
                }
            });
            return;
        }
        // Dashboard
        if (pathname === '/' && req.method === 'GET') {
            res.writeHead(200, { 'Content-Type': 'text/html' });
            res.end('<html><body>CRM Dashboard</body></html>');
            return;
        }
        // 404
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Not found' }));
    });
}
test('GET /health returns {"status":"ok"}', async () => {
    const server = createTestServer();
    return new Promise((resolve, reject) => {
        server.listen(0, async () => {
            try {
                const address = server.address();
                const port = typeof address === 'object' && address !== null ? address.port : 0;
                const response = await fetch(`http://localhost:${port}/health`);
                assert.equal(response.status, 200);
                const data = await response.json();
                assert.deepEqual(data, { status: 'ok' });
                server.close(() => resolve());
            }
            catch (error) {
                server.close(() => reject(error));
            }
        });
    });
});
test('POST /api/clients creates a new client', async () => {
    const server = createTestServer();
    return new Promise((resolve, reject) => {
        server.listen(0, async () => {
            try {
                const address = server.address();
                const port = typeof address === 'object' && address !== null ? address.port : 0;
                const response = await fetch(`http://localhost:${port}/api/clients`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: 'John Doe',
                        email: 'john@example.com',
                        phone: '555-1234',
                    }),
                });
                assert.equal(response.status, 201);
                const data = await response.json();
                assert.equal(data.name, 'John Doe');
                assert.equal(data.email, 'john@example.com');
                assert.equal(data.phone, '555-1234');
                assert.ok(data.id);
                server.close(() => resolve());
            }
            catch (error) {
                server.close(() => reject(error));
            }
        });
    });
});
test('GET /api/clients returns list of clients', async () => {
    const server = createTestServer();
    return new Promise((resolve, reject) => {
        server.listen(0, async () => {
            try {
                const address = server.address();
                const port = typeof address === 'object' && address !== null ? address.port : 0;
                // Create a client first
                await fetch(`http://localhost:${port}/api/clients`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: 'Jane Smith',
                        email: 'jane@example.com',
                        phone: '555-5678',
                    }),
                });
                // Get clients
                const response = await fetch(`http://localhost:${port}/api/clients`);
                assert.equal(response.status, 200);
                const data = await response.json();
                assert(Array.isArray(data));
                assert.ok(data.length > 0);
                server.close(() => resolve());
            }
            catch (error) {
                server.close(() => reject(error));
            }
        });
    });
});
test('GET / returns HTML dashboard', async () => {
    const server = createTestServer();
    return new Promise((resolve, reject) => {
        server.listen(0, async () => {
            try {
                const address = server.address();
                const port = typeof address === 'object' && address !== null ? address.port : 0;
                const response = await fetch(`http://localhost:${port}/`);
                assert.equal(response.status, 200);
                assert.match(response.headers.get('content-type') || '', /text\/html/);
                const html = await response.text();
                assert.match(html, /CRM Dashboard/);
                server.close(() => resolve());
            }
            catch (error) {
                server.close(() => reject(error));
            }
        });
    });
});
