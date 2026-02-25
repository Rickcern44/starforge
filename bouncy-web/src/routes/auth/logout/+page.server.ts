import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
  // Clear the auth cookie on the server
  cookies.delete('access_token', { path: '/' });
  
  // Redirect to login page
  throw redirect(303, '/auth/login');
};
