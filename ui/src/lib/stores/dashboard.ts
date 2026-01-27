import {writable} from 'svelte/store';
import {api} from '$lib/api/client';

export const leagues = writable<any[]>([]);
export const games = writable<any[]>([]);
export const loading = writable(false);
export const error = writable<string | null>(null);

export async function loadDashboard() {
    loading.set(true);
    error.set(null);

    try {
        const [leagueRes, gameRes] = await Promise.all([
            api.get('/leagues'),
            api.get('/games/upcoming')
        ]);

        leagues.set(leagueRes.data);
        games.set(gameRes.data);
    } catch (err: any) {
        error.set(
            err?.response?.data?.error ??
            'Something went wrong loading your dashboard. Please try again.'
        );
    } finally {
        loading.set(false);
    }
}