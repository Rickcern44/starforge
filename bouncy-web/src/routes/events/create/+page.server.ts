import type { PageServerLoad } from './$types';
import { getLeagues } from '$lib/services/league';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ fetch, locals }) => {
  const token = locals.token;

  if (!token) {
    return redirect(303, '/auth/login');
  }

  const leagues = await getLeagues(fetch, token);

  return {
    leagues
  };
};
