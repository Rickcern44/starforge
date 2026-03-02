import type { PageServerLoad } from './$types';
import { getLeagues } from '$lib/services/league';
import { getUserPayments, getUserCharges } from '$lib/services/payment';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ fetch, locals }) => {
  const token = locals.token;

  if (!token) {
    throw redirect(303, '/auth/login');
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
