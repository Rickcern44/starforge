import { writable } from 'svelte/store';
import { api } from '$lib/api/client';

export type Game = {
    id: string;
    leagueId: string;
    startTime: string;
    location: string;
    costCents: number;
};

export const games = writable<Game[]>([]);
export const loading = writable(false);

export async function loadGames(leagueId: string) {
    loading.set(true);
    try {
        const res = await api.get(`/leagues/${leagueId}/games`);
        games.set(res.data);
    } finally {
        loading.set(false);
    }
}