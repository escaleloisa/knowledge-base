'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, Note } from '@/lib/api';
import { MarkdownViewer } from '@/components/MarkdownViewer';

export default function NotePage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;
  const [note, setNote] = useState<Note | null>(null);
  const [backlinks, setBacklinks] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [noteData, backlinksData] = await Promise.allSettled([
          api.notes.get(id),
          api.notes.backlinks(id),
        ]);
        if (noteData.status === 'fulfilled') setNote(noteData.value);
        if (backlinksData.status === 'fulfilled') setBacklinks(backlinksData.value || []);
      } catch (err) {
        console.error('Failed to load note:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [id]);

  if (loading) return <main className="flex-1 p-6"><p className="text-gray-500">Loading...</p></main>;
  if (!note) return <main className="flex-1 p-6"><p className="text-red-500">Note not found</p></main>;

  return (
    <main className="flex-1 overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">{note.title}</h1>
            <div className="flex items-center gap-2 mt-2">
              {note.tags?.map((tag) => (
                <span key={tag} className="px-2 py-0.5 rounded-full text-xs bg-blue-100 text-blue-700">
                  {tag}
                </span>
              ))}
              <span className="text-xs text-gray-400">
                Updated {new Date(note.updated_at).toLocaleDateString()}
              </span>
            </div>
          </div>
          <div className="flex gap-2">
            <Link
              href={`/notes/${id}/edit`}
              className="px-4 py-2 text-sm bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
            >
              Edit
            </Link>
            <button
              onClick={async () => {
                if (confirm('Delete this note?')) {
                  await api.notes.delete(id);
                  router.push('/');
                }
              }}
              className="px-4 py-2 text-sm bg-red-50 text-red-600 rounded-md hover:bg-red-100"
            >
              Delete
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="bg-white rounded-lg border border-gray-200 p-8">
          <MarkdownViewer content={note.content} />
        </div>

        {/* Backlinks */}
        {backlinks.length > 0 && (
          <div className="mt-8 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h2 className="text-sm font-semibold text-gray-600 uppercase tracking-wider mb-3">
              Linked from ({backlinks.length})
            </h2>
            <ul className="space-y-2">
              {backlinks.map((bl) => (
                <li key={bl.id}>
                  <Link href={`/notes/${bl.id}`} className="text-blue-600 hover:underline text-sm">
                    {bl.title}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </main>
  );
}
