import type { Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';
import * as api from '$lib/services/api';
import type { User } from '$lib/models';

export const handle: Handle = async ({ event, resolve }) => {
  const token = event.cookies.get('access_token');
  const path = event.url.pathname;
  const isAuthPage = path.startsWith('/auth');

  event.locals.token = token || null;
  event.locals.user = null;

  if (token) {
    try {
      // Fetch user profile on the server for every request to populate locals
      const user: User = await api.get('users/me', event.fetch, token);
      event.locals.user = user;
    } catch (error) {
      console.error('Hooks: Error fetching user profile:', error);
      // If token is invalid, clear it
      if (!isAuthPage) {
        event.cookies.delete('access_token', { path: '/' });
      }
    }
  }

  // Handle protected routes
  if (!event.locals.user && !isAuthPage && path !== '/') {
    return redirect(303, '/auth/login');
  }

  // If already logged in, don't allow access to auth pages (except logout)
  if (event.locals.user && isAuthPage && path !== '/auth/logout') {
    return redirect(303, '/');
  }

  const response = await resolve(event);
  return response;
};
