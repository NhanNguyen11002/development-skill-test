import { ACCESS_TOKEN } from '$lib/const';
import type { Handle } from '@sveltejs/kit';

const publicPath = ['/login', '/api/auth/login']

export const handle: Handle = async ({ event, resolve }) => {
  const token = event.cookies.get(ACCESS_TOKEN);
  const currentPath = event.url.pathname;

  if (!token && !publicPath.includes(currentPath)) {
    return new Response(null, {
      status: 302,
      headers: { Location: '/login' }
    });
  }

  return resolve(event);
};


