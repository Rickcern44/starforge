import { getFeatureFlags } from '$lib/services/feature-flag';
import type { PageLoad } from './$types';
import { browser } from '$app/environment';

export const load: PageLoad = async ({ fetch }) => {
    if (!browser) return { features: [] };
    
    try {
        const features = await getFeatureFlags(fetch);
        return {
            features
        };
    } catch (error) {
        console.error('Failed to load feature flags:', error);
        return {
            features: []
        };
    }
};
