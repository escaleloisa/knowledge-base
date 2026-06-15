'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { MarkdownViewer } from '@/components/MarkdownViewer';

export default function NewNotePage() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');
  const [saving, setSaving] = useState(false);
  const [showPreview, setShowPreview] = useState(false);

  const handleSave = async () => {
    if (!title || !content) return;
    setSaving(true);
    try {
      const note = await api.notes.create({
        title,
        content,
        tags: tags.split(',').map(t => t.trim()).filter(Boolean),
      });
      router.push(`/notes/${note.id}`);
    } catch (err) {
      console.error('Failed to create note:', err);
    } finally {
      setSaving(false);
    }
  };

  return (
    <main className="flex-1 overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">New Note</h1>
          <div className="flex gap-2">
            <button
              onClick={() => setShowPreview(!showPreview)}
              className="px-3 py-1.5 text-sm border border-gray-300 rounded-md hover:bg-gray-50"
            >
              {showPreview ? 'Edit' : 'Preview'}
            </button>
            <button
              onClick={handleSave}
              disabled={saving || !title || !content}
              className="px-4 py-1.5 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {saving ? 'Creating...' : 'Create'}
            </button>
          </div>
        </div>

        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Note title"
          className="w-full text-2xl font-bold border-0 border-b border-gray-200 pb-2 mb-4 focus:outline-none focus:border-blue-500"
        />

        <input
          type="text"
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          placeholder="Tags (comma separated)"
          className="w-full text-sm border border-gray-200 rounded-md px-3 py-2 mb-4 focus:outline-none focus:border-blue-500"
        />

        {showPreview ? (
          <div className="bg-white rounded-lg border border-gray-200 p-8 min-h-[400px]">
            <MarkdownViewer content={content} />
          </div>
        ) : (
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Write your markdown here..."
            className="w-full h-[500px] font-mono text-sm border border-gray-200 rounded-md p-4 focus:outline-none focus:border-blue-500 resize-none"
          />
        )}
      </div>
    </main>
  );
}
