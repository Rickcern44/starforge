import type { PageLoad } from './$types';
import { getGameById } from '$lib/services/game';
import { getLeagues } from '$lib/services/league';
import { error, redirect } from '@sveltejs/kit';
import { authService } from '$lib/services/auth.svelte';

export const load: PageLoad = async ({ params, fetch }) => {
  const token = authService.token;
  const gameId = params.id;

  if (!token) {
    throw redirect(303, '/auth/login');
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
