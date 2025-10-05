import { PUBLIC_API_URL } from '$env/static/public';
import { ACCESS_TOKEN } from '$lib/const';
import type { RequestEvent } from '@sveltejs/kit';

export async function apiProxy(
  event: RequestEvent, 
  options?: { useAuth?: boolean } 
) {
  const { request, cookies, fetch } = event;
  const token = options?.useAuth ? cookies.get(ACCESS_TOKEN) : undefined;
  console.log('token', token)

  const url = new URL(request.url);
  const targetUrl = `${PUBLIC_API_URL}${url.pathname.replace(/^\/api/, '')}${url.search}`;

  const headers = new Headers(request.headers);
  headers.set('host', new URL(PUBLIC_API_URL).host); // tránh leak host local
  if (token) headers.set('Authorization', `Bearer ${token}`);

  const res = await fetch(targetUrl, {
    method: request.method,
    headers,
    body:
      request.method !== 'GET' && request.method !== 'HEAD'
        ? await request.text()
        : undefined,
  });

  return new Response(await res.text(), {
    status: res.status,
    headers: res.headers,
  });
}
