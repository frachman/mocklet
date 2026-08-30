import './globals.css';
import type { ReactNode } from 'react';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Mocklet — Build and test APIs before the backend is ready',
  description: 'Create disposable mock API endpoints in seconds. Build against a stable contract, unblock frontend work, and share realistic responses with your team.',
  keywords: ['mock API', 'API mocking', 'frontend development', 'REST API testing', 'disposable API endpoint'],
  metadataBase: new URL('https://mocklet.mikrolyt.com'),
  alternates: { canonical: '/' },
  openGraph: {
    title: 'Mocklet — Build against the API contract first',
    description: 'Create a disposable mock API and keep product work moving while the real service is still being built.',
    url: 'https://mocklet.mikrolyt.com',
    siteName: 'Mocklet',
    type: 'website',
  },
  twitter: { card: 'summary', title: 'Mocklet — Build against the API contract first', description: 'Create disposable mock API endpoints in seconds.' },
  robots: { index: true, follow: true },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return <html lang="en"><body>{children}</body></html>;
}
