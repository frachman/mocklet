'use client';

import { FormEvent, useState } from 'react';

const apiOrigin = process.env.NEXT_PUBLIC_API_ORIGIN ?? 'http://localhost:8080';
type Result = { key: string; url: string; token: string; expires: string };

function endpointPayload(data: FormData) {
  return { name: data.get('name') ?? 'endpoint', method: data.get('method'), path: data.get('path'), status_code: Number(data.get('status_code')), content_type: 'application/json', body: data.get('body') };
}

export default function Home() {
  const [result, setResult] = useState<Result | null>(null);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  async function create(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setError(''); setMessage('');
    const data = new FormData(event.currentTarget);
    const response = await fetch(`${apiOrigin}/api/v1/mocks`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(endpointPayload(data)) });
    if (!response.ok) { setError(await response.text()); return; }
    const value = await response.json();
    setResult({ key: value.public_key, url: `${apiOrigin}/m/${value.public_key}${data.get('path')}`, token: value.management_token, expires: value.expires_at });
  }

  async function addEndpoint(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); if (!result) return; setError(''); setMessage('');
    const form = event.currentTarget;
    const response = await fetch(`${apiOrigin}/api/v1/mocks/${result.key}/endpoints`, { method: 'POST', headers: { 'Content-Type': 'application/json', 'X-Management-Token': result.token }, body: JSON.stringify(endpointPayload(new FormData(form))) });
    if (!response.ok) { setError(await response.text()); return; }
    const value = await response.json(); setMessage(`Added ${value.method} ${value.path}`); form.reset();
  }

  return <main>
    <h1>Mocklet</h1>
    <p>Build against the API contract before the API is ready. Create a disposable endpoint and point your client at it.</p>
    <form onSubmit={create}>
      <label htmlFor="name">Mock name</label><input id="name" name="name" defaultValue="My frontend mock" required />
      <label htmlFor="method">Method</label><select id="method" name="method" defaultValue="GET"><option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option></select>
      <label htmlFor="path">Path</label><input id="path" name="path" defaultValue="/users" required />
      <label htmlFor="status_code">Status</label><input id="status_code" name="status_code" type="number" defaultValue="200" min="100" max="599" required />
      <label htmlFor="body">JSON response body</label><textarea id="body" name="body" defaultValue={'{"ok":true}'} />
      <button type="submit">Create disposable mock</button>
    </form>
    {error && <p className="error">{error}</p>}
    {result && <section className="result"><strong>Mock created</strong><p>Public endpoint: <code>{result.url}</code></p><p>Management token (save it now): <code>{result.token}</code></p><small>Expires {result.expires}</small>
      <form onSubmit={addEndpoint}>
        <label htmlFor="add-method">Add method</label><select id="add-method" name="method" defaultValue="GET"><option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option></select>
        <label htmlFor="add-path">Add path</label><input id="add-path" name="path" defaultValue="/health" required />
        <label htmlFor="add-status">Status</label><input id="add-status" name="status_code" type="number" defaultValue="200" min="100" max="599" required />
        <label htmlFor="add-body">Response body</label><textarea id="add-body" name="body" defaultValue={'{"ok":true}'} />
        <button type="submit">Add endpoint</button>
      </form>
      {message && <p>{message}</p>}
    </section>}
  </main>;
}
