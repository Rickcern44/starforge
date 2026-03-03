import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { authService } from '$lib/services/auth.svelte';

export const load: PageLoad = async () => {
  // Clear the auth cookie/token on the client
  authService.logout();
  
  // Redirect to login page
  throw redirect(303, '/auth/login');
};
