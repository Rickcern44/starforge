import { authService } from './auth.svelte';
import { getFeatureFlags } from './feature-flag';
import type { FeatureFlag } from '$lib/models';

class FeatureFlagService {
    #flags = $state<Record<string, boolean>>({});
    #loading = $state(true);

    constructor() {
        this.initialize();
    }

    async initialize() {
        if (!authService.user) {
            this.#loading = false;
            return;
        }

        try {
            const flags = await getFeatureFlags();
            const flagMap: Record<string, boolean> = {};
            flags.forEach(flag => {
                flagMap[flag.key] = flag.enabled;
            });
            this.#flags = flagMap;
        } catch (error) {
            console.error('Failed to load feature flags:', error);
        } finally {
            this.#loading = false;
        }
    }

    isEnabled(key: string): boolean {
        return this.#flags[key] || false;
    }

    get loading() {
        return this.#loading;
    }
}

export const featureFlagService = new FeatureFlagService();
