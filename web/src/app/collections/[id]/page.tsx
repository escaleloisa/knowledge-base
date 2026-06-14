'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { api, Collection, Note } from '@/lib/api';
import { NoteCard } from '@/components/NoteCard';

export default function CollectionPage() {
  const params = useParams();
  const id = params.id as string;
  const [collection, setCollection] = useState<Collection | null>(null);
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [colData, notesData] = await Promise.all([
          api.collections.get(id),
          api.collections.notes(id),
        ]);
        setCollection(colData);
        setNotes(notesData || []);
      } catch (err) {
        console.error('Failed to load collection:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [id]);

  if (loading) return <main className="flex-1 p-6"><p className="text-gray-500">Loading...</p></main>;
  if (!collection) return <main className="flex-1 p-6"><p className="text-red-500">Collection not found</p></main>;

  return (
    <main className="flex-1 overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">{collection.name}</h1>
          {collection.description && (
            <p className="text-gray-500 mt-1">{collection.description}</p>
          )}
          <p className="text-sm text-gray-400 mt-1">{collection.note_count} notes</p>
        </div>

        {notes.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <p>No notes in this collection yet.</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {notes.map((note) => (
              <NoteCard key={note.id} note={note} />
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
