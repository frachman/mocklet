'use client';

import { FormEvent, useState } from 'react';

const apiOrigin = process.env.NEXT_PUBLIC_API_ORIGIN ?? 'http://localhost:8080';

export default function Home() {
  const [result, setResult] = useState<{ url: string; token: string; expires: string } | null>(null);
  const [error, setError] = useState('');
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setError(''); setResult(null);
    const data = new FormData(event.currentTarget);
    const response = await fetch(`${apiOrigin}/api/v1/mocks`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: data.get('name'), method: data.get('method'), path: data.get('path'), status_code: Number(data.get('status_code')), content_type: 'application/json', body: data.get('body') }) });
    if (!response.ok) { setError(await response.text()); return; }
    const value = await response.json();
    setResult({ url: `${apiOrigin}/m/${value.public_key}${data.get('path')}`, token: value.management_token, expires: value.expires_at });
  }
  return <main><h1>Mocklet</h1><p>Build against the API contract before the API is ready. Create a disposable endpoint and point your client at it.</p><form onSubmit={submit}><label htmlFor="name">Mock name</label><input id="name" name="name" defaultValue="My frontend mock" required /><label htmlFor="method">Method</label><select id="method" name="method" defaultValue="GET"><option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option></select><label htmlFor="path">Path</label><input id="path" name="path" defaultValue="/users" required /><label htmlFor="status_code">Status</label><input id="status_code" name="status_code" type="number" defaultValue="200" min="100" max="599" required /><label htmlFor="body">JSON response body</label><textarea id="body" name="body" defaultValue={'{"ok":true}'} /><button type="submit">Create disposable mock</button>{error && <p className="error">{error}</p>}{result && <div className="result"><strong>Mock created</strong><p>Public endpoint: <code>{result.url}</code></p><p>Management token (save it now): <code>{result.token}</code></p><small>Expires {result.expires}</small></div>}</form></main>;
}

