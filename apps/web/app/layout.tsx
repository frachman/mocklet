import './globals.css';
import type { ReactNode } from 'react';
import type { Metadata } from 'next';
import { Inter, Space_Grotesk } from 'next/font/google';

const inter = Inter({ subsets: ['latin'], variable: '--font-inter', display: 'swap' });
const spaceGrotesk = Space_Grotesk({ subsets: ['latin'], variable: '--font-space-grotesk', display: 'swap' });

export const metadata: Metadata = {
  title: 'Mocklet by Mikrolyt — Build and test APIs before the backend is ready',
  description: 'Create disposable mock API endpoints in seconds. Build against a stable contract, unblock frontend work, and share realistic responses with your team. A free tool by Mikrolyt.',
  keywords: ['mock API', 'API mocking', 'frontend development', 'REST API testing', 'disposable API endpoint'],
  metadataBase: new URL('https://mocklet.mikrolyt.com'),
  alternates: { canonical: '/' },
  openGraph: {
    title: 'Mocklet by Mikrolyt — Build against the API contract first',
    description: 'Create a disposable mock API and keep product work moving while the real service is still being built.',
    url: 'https://mocklet.mikrolyt.com',
    siteName: 'Mikrolyt',
    type: 'website',
  },
  twitter: { card: 'summary', title: 'Mocklet by Mikrolyt — Build against the API contract first', description: 'Create disposable mock API endpoints in seconds.' },
  robots: { index: true, follow: true },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return <html lang="en" className={`${inter.variable} ${spaceGrotesk.variable}`}><body>{children}</body></html>;
}
