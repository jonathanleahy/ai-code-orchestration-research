// @ts-nocheck
import { randomUUID } from 'node:crypto';

export interface Client {
  id: string;
  name: string;
  email: string;
  phone: string;
  createdAt: Date;
}

export interface Activity {
  id: string;
  clientId: string;
  description: string;
  createdAt: Date;
}

export interface Invoice {
  id: string;
  clientId: string;
  amount: number;
  status: 'draft' | 'sent' | 'paid';
  createdAt: Date;
}

export class Store {
  private clients: Map<string, Client> = new Map();
  private activities: Map<string, Activity> = new Map();
  private invoices: Map<string, Invoice> = new Map();

  createClient(name: string, email: string, phone: string): Client {
    const id = randomUUID();
    const client: Client = { id, name, email, phone, createdAt: new Date() };
    this.clients.set(id, client);
    return client;
  }

  getClients(): Client[] {
    return Array.from(this.clients.values());
  }

  getClient(id: string): Client | undefined {
    return this.clients.get(id);
  }

  createActivity(clientId: string, description: string): Activity {
    const id = randomUUID();
    const activity: Activity = { id, clientId, description, createdAt: new Date() };
    this.activities.set(id, activity);
    return activity;
  }

  getActivities(clientId: string): Activity[] {
    return Array.from(this.activities.values()).filter(a => a.clientId === clientId);
  }

  createInvoice(clientId: string, amount: number): Invoice {
    const id = randomUUID();
    const invoice: Invoice = {
      id,
      clientId,
      amount,
      status: 'draft',
      createdAt: new Date(),
    };
    this.invoices.set(id, invoice);
    return invoice;
  }

  getInvoices(clientId: string): Invoice[] {
    return Array.from(this.invoices.values()).filter(i => i.clientId === clientId);
  }
}
