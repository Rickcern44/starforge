import * as api from './api';
import type { League } from '$lib/models';

export async function getLeagues(fetch?: api.Fetch, token?: string): Promise<League[]> {
  try {
    const leagues = await api.get('me/leagues', fetch, token);
    return leagues as League[];
  } catch (error) {
    console.error('Error fetching leagues:', error);
    return [];
  }
}
