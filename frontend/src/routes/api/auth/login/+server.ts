// src/routes/api/auth/login/+server.ts
import { json, type RequestHandler } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { ACCESS_TOKEN } from '$lib/const';

export const POST: RequestHandler = async ({ request, cookies, fetch }) => {
  const credentials = await request.json();
  const response = await fetch(`${PUBLIC_API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(credentials)
  });

  const data = await response.json();
  if (!response.ok) return json({ error: data.error }, { status: response.status });

  cookies.set(ACCESS_TOKEN, data.data.token, {
    path: '/',
    httpOnly: true,
    sameSite: 'strict',
    secure: true,
    maxAge: 24 * 60 * 60
  });

  return json(data);
};
