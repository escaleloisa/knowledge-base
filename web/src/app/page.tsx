'use client';

import { useEffect, useState } from 'react';
import { api, Note, Collection, Tag } from '@/lib/api';
import { Sidebar } from '@/components/Sidebar';
import { NoteCard } from '@/components/NoteCard';

export default function Home() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [collections, setCollections] = useState<Collection[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [notesData, collectionsData, tagsData] = await Promise.allSettled([
          api.notes.list(),
          api.collections.list(),
          api.tags.list(),
        ]);
        if (notesData.status === 'fulfilled') setNotes(notesData.value || []);
        if (collectionsData.status === 'fulfilled') setCollections(collectionsData.value || []);
        if (tagsData.status === 'fulfilled') setTags(tagsData.value || []);
      } catch (err) {
        console.error('Failed to load data:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  return (
    <>
      <Sidebar collections={collections} tags={tags} />
      <main className="flex-1 overflow-y-auto p-6">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">All Notes</h1>
        {loading ? (
          <p className="text-gray-500">Loading...</p>
        ) : notes.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <p className="text-lg mb-2">No notes yet</p>
            <p className="text-sm">Create your first note or import markdown files.</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {notes.map((note) => (
              <NoteCard key={note.id} note={note} />
            ))}
          </div>
        )}
      </main>
    </>
  );
}
