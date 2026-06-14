'use client';

import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import Link from 'next/link';

interface MarkdownViewerProps {
  content: string;
}

export function MarkdownViewer({ content }: MarkdownViewerProps) {
  // Replace [[wiki-links]] with clickable links
  const processedContent = content.replace(
    /\[\[([^\]]+)\]\]/g,
    '[$1](/notes?search=$1)'
  );

  return (
    <article className="prose prose-lg max-w-none dark:prose-invert">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
        components={{
          a: ({ href, children }) => {
            if (href?.startsWith('/')) {
              return <Link href={href} className="text-blue-600 hover:underline">{children}</Link>;
            }
            return <a href={href} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">{children}</a>;
          },
        }}
      >
        {processedContent}
      </ReactMarkdown>
    </article>
  );
}
