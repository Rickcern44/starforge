import { getLeagues } from '$lib/services/league';
import type { PageLoad } from './$types';
import { browser } from '$app/environment';

export const load: PageLoad = async ({ fetch }) => {
    if (!browser) return { leagues: [] };
    
    try {
        const leagues = await getLeagues(fetch);
        return {
            leagues
        };
    } catch (error) {
        console.error('Failed to load leagues:', error);
        return {
            leagues: []
        };
    }
};
