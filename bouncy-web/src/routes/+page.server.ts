import type { PageServerLoad } from './$types';
import { getLeagues } from '$lib/services/league';

export const load: PageServerLoad = async ({ fetch, locals }) => {
  const token = locals.token;
  
  if (!token) {
    return {
      leagues: []
    };
  }

  const leagues = await getLeagues(fetch, token);

  return {
    leagues
  };
};
