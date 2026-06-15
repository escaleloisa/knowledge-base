'use client';

import { useEffect, useState } from 'react';
import { api, Collection } from '@/lib/api';
import { FileImport } from '@/components/FileImport';

export default function ImportPage() {
  const [collections, setCollections] = useState<Collection[]>([]);

  useEffect(() => {
    api.collections.list().then(setCollections).catch(() => {});
  }, []);

  return (
    <main className="flex-1 overflow-y-auto p-6">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Import Markdown Files</h1>
        <p className="text-gray-500 mb-6">Upload .md files to create notes. Titles are extracted from the first heading or filename.</p>
        <FileImport collections={collections} />
      </div>
    </main>
  );
}
