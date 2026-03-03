import type { PageLoad } from './$types';
import { getLeagues } from '$lib/services/league';
import { redirect } from '@sveltejs/kit';
import { authService } from '$lib/services/auth.svelte';

export const load: PageLoad = async ({ fetch }) => {
  const token = authService.token;

  if (!token) {
    throw redirect(303, '/auth/login');
  }

  const leagues = await getLeagues(fetch, token);

  return {
    leagues
  };
};
