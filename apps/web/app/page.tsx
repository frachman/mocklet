'use client';

import { FormEvent, useState } from 'react';

const apiOrigin = process.env.NEXT_PUBLIC_API_ORIGIN ?? 'http://localhost:8080';
type Result = { key: string; url: string; token: string; expires: string };

function endpointPayload(data: FormData) {
  return { name: data.get('name') ?? 'endpoint', method: data.get('method'), path: data.get('path'), status_code: Number(data.get('status_code')), content_type: 'application/json', body: data.get('body') ?? data.get('add-body') };
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
    <header className="site-header"><a className="brand" href="#top" aria-label="Mocklet home"><span className="brand-mark">M</span> Mocklet</a><nav aria-label="Main navigation"><a href="#how-it-works">How it works</a><a href="#case-study">Case study</a><a href="#playground">Try it</a></nav></header>

    <section className="hero" id="top"><div className="eyebrow"><span className="status-dot" /> API contract, made executable</div><h1>Build the frontend <em>before</em> the backend is ready.</h1><p className="hero-copy">Create a disposable mock API, point your app at a stable response, and keep shipping while the real service is still being built.</p><div className="hero-actions"><a className="button primary" href="#playground">Create a mock API <span>→</span></a><a className="text-link" href="#case-study">See a realistic workflow <span>↓</span></a></div><div className="hero-proof"><span>✓ No account required</span><span>✓ Expires automatically</span><span>✓ Shareable HTTP endpoints</span></div></section>

    <section className="demo-strip" aria-label="Example mock API response"><div className="demo-caption"><span className="mini-label">A response your frontend can use today</span><strong>GET /checkout/session</strong><small>Mocklet endpoint · 200 OK</small></div><pre><code>{`{
  "session_id": "sess_demo_123",
  "status": "ready",
  "total": 4999,
  "currency": "IDR"
}`}</code></pre><div className="demo-note"><span>↗</span><p>Replace it with your real API later. Your UI contract stays the same.</p></div></section>

    <section className="section" id="how-it-works"><div className="section-heading"><div><span className="eyebrow">A smaller feedback loop</span><h2>From idea to working integration in minutes.</h2></div><p>Mocklet gives product teams a simple seam between “we know what the API should do” and “the API is deployed.”</p></div><div className="steps"><article><span className="step-number">01</span><h3>Define the contract</h3><p>Choose a method, path, status code, and JSON response that match the interface your client needs.</p></article><article><span className="step-number">02</span><h3>Point your app at it</h3><p>Use the public URL in your local environment, story, prototype, or integration test.</p></article><article><span className="step-number">03</span><h3>Replace it when ready</h3><p>When the backend ships, swap the base URL. The frontend work is already exercised against a real HTTP response.</p></article></div></section>

    <section className="case-study" id="case-study"><div className="case-copy"><span className="eyebrow">Realistic workflow case study</span><h2>Checkout UI today.<br /><span>Payment service later.</span></h2><p>A frontend engineer is building a checkout flow. The payment contract is agreed, but the service is still in development. Instead of hardcoding loading states or waiting on a shared environment, the team creates a disposable Mocklet endpoint.</p><a className="text-link" href="#playground">Create this example <span>→</span></a></div><div className="case-timeline"><div className="timeline-item"><span className="timeline-dot">1</span><div><strong>Frontend owns the contract</strong><p>Checkout calls <code>GET /checkout/session</code> and renders the agreed response.</p></div></div><div className="timeline-line" /><div className="timeline-item"><span className="timeline-dot">2</span><div><strong>Mocklet unblocks the build</strong><p>The UI can be reviewed with realistic data before payment infrastructure is deployed.</p></div></div><div className="timeline-line" /><div className="timeline-item"><span className="timeline-dot muted">3</span><div><strong>Production API takes over</strong><p>Change the environment URL when the real endpoint is ready. The mock expires on its own.</p></div></div></div></section>

    <section className="playground section" id="playground"><div className="section-heading compact"><div><span className="eyebrow">Interactive playground</span><h2>Create your first disposable endpoint.</h2></div><p>Start with the checkout example, then add more routes to the same mock API.</p></div><form onSubmit={create}><div className="form-heading"><div><span className="form-kicker">Step 1</span><h3>Configure your mock</h3></div><span className="expiry-note">24h disposable by default</span></div><div className="form-grid"><div><label htmlFor="name">Mock name</label><input id="name" name="name" defaultValue="Checkout session" required /></div><div><label htmlFor="method">Method</label><select id="method" name="method" defaultValue="GET"><option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option></select></div></div><label htmlFor="path">Path</label><input id="path" name="path" defaultValue="/checkout/session" required /><div className="form-grid"><div><label htmlFor="status_code">Status code</label><input id="status_code" name="status_code" type="number" defaultValue="200" min="100" max="599" required /></div><div><label htmlFor="content-type">Content type</label><input id="content-type" value="application/json" readOnly /></div></div><label htmlFor="body">JSON response body</label><textarea id="body" name="body" defaultValue={'{"session_id":"sess_demo_123","status":"ready","total":4999,"currency":"IDR"}'} /><button className="button primary" type="submit">Create disposable mock <span>→</span></button></form>{error && <p className="error" role="alert">{error}</p>}{result && <section className="result"><div className="result-heading"><span className="success-icon">✓</span><div><strong>Mock API created</strong><p>Your endpoint is ready to use.</p></div></div><p className="result-label">Public endpoint</p><p><code>{result.url}</code></p><p className="result-label">Management token <span>save it now — it will not be shown again</span></p><p><code>{result.token}</code></p><small>Expires {result.expires}</small><form onSubmit={addEndpoint} className="add-form"><div className="form-heading"><div><span className="form-kicker">Step 2</span><h3>Add another route</h3></div></div><div className="form-grid"><div><label htmlFor="add-method">Method</label><select id="add-method" name="method" defaultValue="GET"><option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option></select></div><div><label htmlFor="add-path">Path</label><input id="add-path" name="path" defaultValue="/health" required /></div></div><label htmlFor="add-status">Status code</label><input id="add-status" name="status_code" type="number" defaultValue="200" min="100" max="599" required /><label htmlFor="add-body">Response body</label><textarea id="add-body" name="body" defaultValue={'{"ok":true}'} /><button className="button secondary" type="submit">Add endpoint</button>{message && <p className="success-message" role="status">{message}</p>}</form></section>}</section>

    <section className="faq section"><div className="section-heading compact"><div><span className="eyebrow">Before you start</span><h2>Simple by design.</h2></div></div><div className="faq-grid"><details open><summary>What is a mock API?</summary><p>A mock API is a temporary HTTP endpoint that returns a predictable response. It lets your frontend, QA flow, or integration prototype work before the real backend is available.</p></details><details><summary>How long does a mock last?</summary><p>Anonymous mocks are disposable and expire automatically. This keeps the playground useful for development without turning it into a permanent data store.</p></details><details><summary>Can I use it for production traffic?</summary><p>Mocklet is designed for development, prototyping, and controlled integration work—not as a replacement for your production backend.</p></details><details><summary>Do I need an account?</summary><p>No account is required for the current anonymous-first workflow. Keep the management token private when you create a mock.</p></details></div></section>

    <footer><div className="brand"><span className="brand-mark">M</span> Mocklet</div><p>Build against the contract. Ship with confidence.</p><span>Disposable mock APIs for teams in motion.</span></footer>
  </main>;
}
