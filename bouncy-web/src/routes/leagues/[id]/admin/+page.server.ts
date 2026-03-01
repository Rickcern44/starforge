import type { PageServerLoad } from './$types';
import { getLeagueById, getLeagueInvitations } from '$lib/services/league';
import { getPaymentsForLeague } from '$lib/services/payment';
import { error, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch, locals }) => {
  const token = locals.token;
  const leagueId = params.id;

  if (!token) {
    throw redirect(303, '/auth/login');
  }

  const [league, payments, invitations] = await Promise.all([
    getLeagueById(leagueId, fetch, token),
    getPaymentsForLeague(leagueId, fetch, token),
    getLeagueInvitations(leagueId, fetch, token)
  ]);

  if (!league) {
    throw error(404, 'League not found');
  }

  // Check if user is admin of this league
  const member = league.members.find(m => m.playerId === locals.user?.id);
  const isAdmin = member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'));

  if (!isAdmin) {
    throw error(403, 'You do not have permission to access the admin page for this league.');
  }

  return {
    league,
    payments,
    invitations
  };
};
