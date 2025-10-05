// src/routes/api/[...path]/+server.ts
import { apiProxy } from '$lib/server/apiProxy';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = (event) => apiProxy(event, { useAuth: true });
export const POST: RequestHandler = (event) => apiProxy(event, { useAuth: true });
export const PUT: RequestHandler = (event) => apiProxy(event, { useAuth: true });
export const DELETE: RequestHandler = (event) => apiProxy(event, { useAuth: true });
