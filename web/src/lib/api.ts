const NOTE_SERVICE = process.env.NEXT_PUBLIC_NOTE_SERVICE_URL || 'http://localhost:8080';
const SEARCH_API = process.env.NEXT_PUBLIC_SEARCH_API_URL || 'http://localhost:8081';
const TAG_API = process.env.NEXT_PUBLIC_TAG_API_URL || 'http://localhost:8082';

export interface Note {
  id: string;
  title: string;
  content: string;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface Collection {
  id: string;
  name: string;
  description: string;
  note_count: number;
  created_at: string;
  updated_at: string;
}

export interface Tag {
  name: string;
  count: number;
}

export interface SearchResult {
  total: number;
  results: (Note & { _score: number; _highlight?: Record<string, string[]> })[];
}

const headers = { 'Content-Type': 'application/json' };

export const api = {
  notes: {
    list: async (limit = 20, offset = 0): Promise<Note[]> => {
      const res = await fetch(`${NOTE_SERVICE}/api/notes?limit=${limit}&offset=${offset}`);
      return res.json();
    },
    get: async (id: string): Promise<Note> => {
      const res = await fetch(`${NOTE_SERVICE}/api/notes/${id}`);
      return res.json();
    },
    create: async (data: { title: string; content: string; tags: string[] }): Promise<Note> => {
      const res = await fetch(`${NOTE_SERVICE}/api/notes`, {
        method: 'POST', headers, body: JSON.stringify(data),
      });
      return res.json();
    },
    update: async (id: string, data: { title: string; content: string; tags: string[] }): Promise<Note> => {
      const res = await fetch(`${NOTE_SERVICE}/api/notes/${id}`, {
        method: 'PUT', headers, body: JSON.stringify(data),
      });
      return res.json();
    },
    delete: async (id: string): Promise<void> => {
      await fetch(`${NOTE_SERVICE}/api/notes/${id}`, { method: 'DELETE' });
    },
    backlinks: async (id: string): Promise<Note[]> => {
      const res = await fetch(`${NOTE_SERVICE}/api/notes/${id}/backlinks`);
      return res.json();
    },
  },
  collections: {
    list: async (): Promise<Collection[]> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections`);
      return res.json();
    },
    get: async (id: string): Promise<Collection> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections/${id}`);
      return res.json();
    },
    create: async (data: { name: string; description: string }): Promise<Collection> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections`, {
        method: 'POST', headers, body: JSON.stringify(data),
      });
      return res.json();
    },
    update: async (id: string, data: { name: string; description: string }): Promise<Collection> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections/${id}`, {
        method: 'PUT', headers, body: JSON.stringify(data),
      });
      return res.json();
    },
    delete: async (id: string): Promise<void> => {
      await fetch(`${NOTE_SERVICE}/api/collections/${id}`, { method: 'DELETE' });
    },
    notes: async (id: string, limit = 20, offset = 0): Promise<Note[]> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections/${id}/notes?limit=${limit}&offset=${offset}`);
      return res.json();
    },
    addNotes: async (id: string, noteIds: string[]): Promise<{ added: number }> => {
      const res = await fetch(`${NOTE_SERVICE}/api/collections/${id}/notes`, {
        method: 'POST', headers, body: JSON.stringify({ note_ids: noteIds }),
      });
      return res.json();
    },
    removeNote: async (id: string, noteId: string): Promise<void> => {
      await fetch(`${NOTE_SERVICE}/api/collections/${id}/notes/${noteId}`, { method: 'DELETE' });
    },
  },
  search: {
    query: async (q: string, tags?: string, from?: string, to?: string): Promise<SearchResult> => {
      const params = new URLSearchParams({ q });
      if (tags) params.set('tags', tags);
      if (from) params.set('from', from);
      if (to) params.set('to', to);
      const res = await fetch(`${SEARCH_API}/api/search?${params}`);
      return res.json();
    },
  },
  tags: {
    list: async (): Promise<Tag[]> => {
      const res = await fetch(`${TAG_API}/api/tags`);
      return res.json();
    },
  },
};
