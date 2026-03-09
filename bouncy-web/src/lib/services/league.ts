import * as api from './api';
import type { League, Invitation } from '$lib/models';

export async function getLeagues(fetch?: api.Fetch, token?: string): Promise<League[]> {
  try {
    const leagues = await api.get('me/leagues', fetch, token);
    return leagues as League[];
  } catch (error) {
    console.error('Error fetching leagues:', error);
    return [];
  }
}

export async function getLeagueById(leagueId: string, fetch?: api.Fetch, token?: string): Promise<League | null> {
  try {
    const league = await api.get(`league/${leagueId}`, fetch, token);
    return league as League;
  } catch (error) {
    console.error(`Error fetching league ${leagueId}:`, error);
    return null;
  }
}

export async function inviteUser(email: string, leagueId: string): Promise<boolean> {
  try {
    await api.post('admin/invite', { email, leagueId });
    return true;
  } catch (error) {
    console.error('Error inviting user:', error);
    return false;
  }
}

export async function getLeagueInvitations(leagueId: string, fetch?: api.Fetch, token?: string): Promise<Invitation[]> {
  try {
    const invites = await api.get(`admin/league/${leagueId}/invitations`, fetch, token);
    return invites as Invitation[];
  } catch (error) {
    console.error(`Error fetching invitations for league ${leagueId}:`, error);
    return [];
  }
}

export async function createLeague(name: string): Promise<League | null> {
  try {
    const league = await api.post('league', { name });
    return league as League;
  } catch (error) {
    console.error('Error creating league:', error);
    return null;
  }
}
