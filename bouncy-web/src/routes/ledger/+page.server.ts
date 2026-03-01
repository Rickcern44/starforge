import type { PageServerLoad } from './$types';
import { getUserPayments, getUserCharges } from '$lib/services/payment';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ fetch, locals }) => {
  const token = locals.token;

  if (!token) {
    throw redirect(303, '/auth/login');
  }

  const [payments, charges] = await Promise.all([
    getUserPayments(fetch, token),
    getUserCharges(fetch, token)
  ]);

  return {
    payments,
    charges
  };
};
