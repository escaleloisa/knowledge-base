'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Collection, Tag } from '@/lib/api';

interface SidebarProps {
  collections: Collection[];
  tags: Tag[];
}

export function Sidebar({ collections, tags }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-gray-50 border-r border-gray-200 h-full overflow-y-auto p-4">
      {/* Navigation */}
      <nav className="mb-6">
        <Link
          href="/"
          className={`block px-3 py-2 rounded-md text-sm font-medium ${
            pathname === '/' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
          }`}
        >
          All Notes
        </Link>
        <Link
          href="/search"
          className={`block px-3 py-2 rounded-md text-sm font-medium ${
            pathname === '/search' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
          }`}
        >
          Search
        </Link>
        <Link
          href="/import"
          className={`block px-3 py-2 rounded-md text-sm font-medium ${
            pathname === '/import' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
          }`}
        >
          Import Files
        </Link>
      </nav>

      {/* Collections */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Collections</h3>
          <Link href="/collections/new" className="text-blue-600 hover:text-blue-800 text-sm">+</Link>
        </div>
        <ul className="space-y-1">
          {collections.map((col) => (
            <li key={col.id}>
              <Link
                href={`/collections/${col.id}`}
                className={`block px-3 py-1.5 rounded-md text-sm ${
                  pathname === `/collections/${col.id}` ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                <span>{col.name}</span>
                <span className="ml-2 text-xs text-gray-400">{col.note_count}</span>
              </Link>
            </li>
          ))}
          {collections.length === 0 && (
            <li className="px-3 py-1.5 text-sm text-gray-400 italic">No collections yet</li>
          )}
        </ul>
      </div>

      {/* Tags */}
      <div>
        <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Tags</h3>
        <div className="flex flex-wrap gap-1.5">
          {tags.map((tag) => (
            <Link
              key={tag.name}
              href={`/search?tags=${tag.name}`}
              className="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-gray-200 text-gray-700 hover:bg-blue-100 hover:text-blue-700"
            >
              {tag.name}
              <span className="ml-1 text-gray-400">{tag.count}</span>
            </Link>
          ))}
          {tags.length === 0 && (
            <span className="text-sm text-gray-400 italic">No tags yet</span>
          )}
        </div>
      </div>
    </aside>
  );
}
