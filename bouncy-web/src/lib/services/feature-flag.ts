import * as api from './api';
import type { FeatureFlag } from '$lib/models';

export async function getFeatureFlags(fetch?: api.Fetch, token?: string): Promise<FeatureFlag[]> {
    const flags = await api.get('admin/features', fetch, token);
    return flags as FeatureFlag[];
}

export async function toggleFeatureFlag(key: string, enabled: boolean): Promise<boolean> {
    try {
        await api.patch(`admin/features/${key}`, { enabled });
        return true;
    } catch (error) {
        console.error(`Error toggling feature flag ${key}:`, error);
        return false;
    }
}
