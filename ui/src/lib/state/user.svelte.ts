import { authState } from '$lib/auth';

/**
 * User state hook - delegates to authState for Keycloak integration
 * Kept for backward compatibility with existing components
 */
export function useUser() {
    return {
        get user() {
            return authState.user;
        },
        get isAuthenticated() {
            return authState.isAuthenticated;
        },
        get isLoading() {
            return authState.isLoading;
        },
        get isInitialized() {
            return authState.isInitialized;
        },
        async login() {
            await authState.login();
        },
        async logout() {
            await authState.logout();
        },
        hasRole(role: string) {
            return authState.hasRole(role);
        }
    };
}
