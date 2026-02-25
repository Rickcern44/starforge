import type { PageServerLoad } from './$types';
import { getGameById } from '$lib/services/game';
import { getLeagues } from '$lib/services/league';
import { error, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch, locals }) => {
  const token = locals.token;
  const gameId = params.id;

  if (!token) {
    return redirect(303, '/auth/login');
  }

  const [game, leagues] = await Promise.all([
    getGameById(gameId, fetch, token),
    getLeagues(fetch, token)
  ]);

  if (!game) {
    throw error(404, 'Game not found');
  }

  return {
    game,
    leagues
  };
};
