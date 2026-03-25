// @ts-nocheck
import { randomUUID } from 'node:crypto';
export class Store {
    constructor() {
        this.clients = new Map();
        this.activities = new Map();
        this.invoices = new Map();
    }
    createClient(name, email, phone) {
        const id = randomUUID();
        const client = { id, name, email, phone, createdAt: new Date() };
        this.clients.set(id, client);
        return client;
    }
    getClients() {
        return Array.from(this.clients.values());
    }
    getClient(id) {
        return this.clients.get(id);
    }
    createActivity(clientId, description) {
        const id = randomUUID();
        const activity = { id, clientId, description, createdAt: new Date() };
        this.activities.set(id, activity);
        return activity;
    }
    getActivities(clientId) {
        return Array.from(this.activities.values()).filter(a => a.clientId === clientId);
    }
    createInvoice(clientId, amount) {
        const id = randomUUID();
        const invoice = {
            id,
            clientId,
            amount,
            status: 'draft',
            createdAt: new Date(),
        };
        this.invoices.set(id, invoice);
        return invoice;
    }
    getInvoices(clientId) {
        return Array.from(this.invoices.values()).filter(i => i.clientId === clientId);
    }
}
