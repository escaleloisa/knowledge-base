import Link from 'next/link';
import { Note } from '@/lib/api';

interface NoteCardProps {
  note: Note;
}

export function NoteCard({ note }: NoteCardProps) {
  const preview = note.content.slice(0, 150).replace(/[#*_`]/g, '');
  const date = new Date(note.created_at).toLocaleDateString();

  return (
    <Link
      href={`/notes/${note.id}`}
      className="block p-4 border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-sm transition-all"
    >
      <h3 className="font-semibold text-gray-900 mb-1">{note.title}</h3>
      <p className="text-sm text-gray-500 mb-2 line-clamp-2">{preview}...</p>
      <div className="flex items-center gap-2">
        {note.tags?.map((tag) => (
          <span key={tag} className="inline-flex px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">
            {tag}
          </span>
        ))}
        <span className="text-xs text-gray-400 ml-auto">{date}</span>
      </div>
    </Link>
  );
}
