import type { PageLoad } from './$types';
import { getLeagues } from '$lib/services/league';
import { authService } from '$lib/services/auth.svelte';

export const load: PageLoad = async ({ fetch }) => {
  const token = authService.token;
  
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
