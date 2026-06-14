'use client';

import { useState, useCallback } from 'react';
import { api } from '@/lib/api';

interface ParsedFile {
  filename: string;
  title: string;
  content: string;
  tags: string[];
}

interface FileImportProps {
  collections: { id: string; name: string }[];
}

export function FileImport({ collections }: FileImportProps) {
  const [files, setFiles] = useState<ParsedFile[]>([]);
  const [selectedCollection, setSelectedCollection] = useState('');
  const [importing, setImporting] = useState(false);
  const [results, setResults] = useState<{ success: number; failed: number } | null>(null);

  const parseFile = (file: File): Promise<ParsedFile> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        // Extract title from first # heading or filename
        const headingMatch = content.match(/^#\s+(.+)$/m);
        const title = headingMatch ? headingMatch[1] : file.name.replace(/\.md$/, '');

        // Extract tags from frontmatter if present
        const tags: string[] = [];
        const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
        if (frontmatterMatch) {
          const tagsMatch = frontmatterMatch[1].match(/tags:\s*\[([^\]]*)\]/);
          if (tagsMatch) {
            tags.push(...tagsMatch[1].split(',').map(t => t.trim().replace(/['"]/g, '')));
          }
        }

        resolve({ filename: file.name, title, content, tags });
      };
      reader.onerror = reject;
      reader.readAsText(file);
    });
  };

  const handleDrop = useCallback(async (e: React.DragEvent) => {
    e.preventDefault();
    const droppedFiles = Array.from(e.dataTransfer.files).filter(f => f.name.endsWith('.md'));
    const parsed = await Promise.all(droppedFiles.map(parseFile));
    setFiles(prev => [...prev, ...parsed]);
  }, []);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) return;
    const selectedFiles = Array.from(e.target.files).filter(f => f.name.endsWith('.md'));
    const parsed = await Promise.all(selectedFiles.map(parseFile));
    setFiles(prev => [...prev, ...parsed]);
  };

  const handleImport = async () => {
    setImporting(true);
    let success = 0;
    let failed = 0;
    const noteIds: string[] = [];

    for (const file of files) {
      try {
        const note = await api.notes.create({
          title: file.title,
          content: file.content,
          tags: file.tags,
        });
        noteIds.push(note.id);
        success++;
      } catch {
        failed++;
      }
    }

    // Add to collection if selected
    if (selectedCollection && noteIds.length > 0) {
      try {
        await api.collections.addNotes(selectedCollection, noteIds);
      } catch {
        // Collection add failed but notes were created
      }
    }

    setResults({ success, failed });
    setFiles([]);
    setImporting(false);
  };

  const removeFile = (index: number) => {
    setFiles(prev => prev.filter((_, i) => i !== index));
  };

  return (
    <div className="max-w-3xl mx-auto">
      {/* Drop zone */}
      <div
        onDrop={handleDrop}
        onDragOver={(e) => e.preventDefault()}
        className="border-2 border-dashed border-gray-300 rounded-lg p-12 text-center hover:border-blue-400 transition-colors"
      >
        <div className="text-gray-500">
          <p className="text-lg mb-2">Drag & drop .md files here</p>
          <p className="text-sm mb-4">or</p>
          <label className="cursor-pointer px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
            Browse Files
            <input
              type="file"
              multiple
              accept=".md"
              onChange={handleFileSelect}
              className="hidden"
            />
          </label>
        </div>
      </div>

      {/* File list */}
      {files.length > 0 && (
        <div className="mt-6">
          <h3 className="font-semibold mb-3">Files to import ({files.length})</h3>
          <div className="space-y-2">
            {files.map((file, i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-md">
                <div>
                  <p className="font-medium text-sm">{file.title}</p>
                  <p className="text-xs text-gray-500">{file.filename}</p>
                </div>
                <button onClick={() => removeFile(i)} className="text-red-500 hover:text-red-700 text-sm">
                  Remove
                </button>
              </div>
            ))}
          </div>

          {/* Collection selector */}
          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Add to collection (optional)
            </label>
            <select
              value={selectedCollection}
              onChange={(e) => setSelectedCollection(e.target.value)}
              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            >
              <option value="">None</option>
              {collections.map((col) => (
                <option key={col.id} value={col.id}>{col.name}</option>
              ))}
            </select>
          </div>

          {/* Import button */}
          <button
            onClick={handleImport}
            disabled={importing}
            className="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {importing ? 'Importing...' : `Import ${files.length} file${files.length > 1 ? 's' : ''}`}
          </button>
        </div>
      )}

      {/* Results */}
      {results && (
        <div className="mt-4 p-4 bg-green-50 border border-green-200 rounded-md">
          <p className="text-green-800">
            Imported {results.success} file{results.success !== 1 ? 's' : ''} successfully.
            {results.failed > 0 && ` ${results.failed} failed.`}
          </p>
        </div>
      )}
    </div>
  );
}
