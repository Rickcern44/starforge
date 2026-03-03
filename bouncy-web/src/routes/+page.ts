import type { PageLoad } from './$types';
import { getLeagues } from '$lib/services/league';
import { getUserPayments, getUserCharges } from '$lib/services/payment';
import { authService } from '$lib/services/auth.svelte';

export const load: PageLoad = async ({ fetch }) => {
  const token = authService.token;
  
  if (!token) {
    return {
      leagues: [],
      payments: [],
      charges: []
    };
  }

  const [leagues, payments, charges] = await Promise.all([
    getLeagues(fetch, token),
    getUserPayments(fetch, token),
    getUserCharges(fetch, token)
  ]);

  return {
    leagues,
    payments,
    charges
  };
};
