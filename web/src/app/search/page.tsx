'use client';

import { Suspense, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { api, SearchResult } from '@/lib/api';

function SearchContent() {
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get('q') || '';
  const initialTags = searchParams.get('tags') || '';

  const [query, setQuery] = useState(initialQuery);
  const [tags, setTags] = useState(initialTags);
  const [results, setResults] = useState<SearchResult | null>(null);
  const [searching, setSearching] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    setSearching(true);
    try {
      const data = await api.search.query(query, tags || undefined);
      setResults(data);
    } catch (err) {
      console.error('Search failed:', err);
    } finally {
      setSearching(false);
    }
  };

  return (
    <main className="flex-1 overflow-y-auto p-6">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Search</h1>

        <form onSubmit={handleSearch} className="mb-6">
          <div className="flex gap-2">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search notes..."
              className="flex-1 border border-gray-300 rounded-md px-4 py-2 focus:outline-none focus:border-blue-500"
            />
            <button
              type="submit"
              disabled={searching}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {searching ? '...' : 'Search'}
            </button>
          </div>
          <input
            type="text"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            placeholder="Filter by tags (comma separated)"
            className="w-full mt-2 border border-gray-200 rounded-md px-4 py-2 text-sm focus:outline-none focus:border-blue-500"
          />
        </form>

        {/* Results */}
        {results && (
          <div>
            <p className="text-sm text-gray-500 mb-4">{results.total} result{results.total !== 1 ? 's' : ''} found</p>
            <div className="space-y-4">
              {results.results.map((result) => (
                <Link
                  key={result.id}
                  href={`/notes/${result.id}`}
                  className="block p-4 border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-sm"
                >
                  <div className="flex items-center justify-between mb-1">
                    <h3 className="font-semibold text-gray-900">{result.title}</h3>
                    <span className="text-xs text-gray-400">Score: {result._score.toFixed(2)}</span>
                  </div>
                  {result._highlight?.content && (
                    <p
                      className="text-sm text-gray-600 mt-1"
                      dangerouslySetInnerHTML={{ __html: result._highlight.content.join('...') }}
                    />
                  )}
                  <div className="flex gap-1 mt-2">
                    {result.tags?.map((tag) => (
                      <span key={tag} className="px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">
                        {tag}
                      </span>
                    ))}
                  </div>
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </main>
  );
}

export default function SearchPage() {
  return (
    <Suspense fallback={<main className="flex-1 p-6"><p className="text-gray-500">Loading search...</p></main>}>
      <SearchContent />
    </Suspense>
  );
}