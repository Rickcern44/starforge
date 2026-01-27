import { keycloakService } from './keycloak.service';
import { User } from '$lib/domain/user';
import { browser } from '$app/environment';

/**
 * Reactive authentication state using Svelte 5 runes
 */
class AuthState {
    private _user = $state<User | null>(null);
    private _isInitialized = $state(false);
    private _isLoading = $state(true);
    private _error = $state<string | null>(null);

    get user() {
        return this._user;
    }

    get isAuthenticated() {
        return this._user !== null;
    }

    get isInitialized() {
        return this._isInitialized;
    }

    get isLoading() {
        return this._isLoading;
    }

    get error() {
        return this._error;
    }

    /**
     * Initialize authentication - call this on app startup
     */
    async initialize(): Promise<void> {
        if (!browser) return;
        if (this._isInitialized) return;

        this._isLoading = true;
        this._error = null;

        try {
            const authenticated = await keycloakService.init();
            
            if (authenticated) {
                this.syncUserFromKeycloak();
            }

            this._isInitialized = true;
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Authentication initialization failed';
            console.error('[Auth] Init error:', err);
        } finally {
            this._isLoading = false;
        }
    }

    /**
     * Sync user state from Keycloak token
     */
    private syncUserFromKeycloak(): void {
        const profile = keycloakService.userProfile;
        
        if (profile) {
            this._user = new User(
                profile.id,
                profile.email,
                profile.name,
                profile.roles
            );
        } else {
            this._user = null;
        }
    }

    /**
     * Login with email/password (redirects to Keycloak)
     */
    async login(): Promise<void> {
        this._isLoading = true;
        this._error = null;
        
        try {
            await keycloakService.login();
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Login failed';
            this._isLoading = false;
        }
    }

    /**
     * Login with passwordless (WebAuthn/Passkeys)
     */
    async loginPasswordless(): Promise<void> {
        this._isLoading = true;
        this._error = null;
        
        try {
            await keycloakService.loginWithMethod('passwordless' as any);
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Passwordless login failed';
            this._isLoading = false;
        }
    }

    /**
     * Request magic link login
     */
    async requestMagicLink(email: string): Promise<void> {
        this._isLoading = true;
        this._error = null;
        
        try {
            await keycloakService.loginWithMethod('magic-link' as any, email);
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Magic link request failed';
            this._isLoading = false;
        }
    }

    /**
     * Login with social provider
     */
    async loginWithProvider(provider: 'google' | 'facebook' | 'apple'): Promise<void> {
        this._isLoading = true;
        this._error = null;
        
        try {
            await keycloakService.loginWithMethod(provider as any);
        } catch (err) {
            this._error = err instanceof Error ? err.message : `${provider} login failed`;
            this._isLoading = false;
        }
    }

    /**
     * Register new user
     */
    async register(email?: string): Promise<void> {
        this._isLoading = true;
        this._error = null;
        
        try {
            await keycloakService.register({ loginHint: email });
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Registration failed';
            this._isLoading = false;
        }
    }

    /**
     * Logout
     */
    async logout(): Promise<void> {
        this._isLoading = true;
        
        try {
            await keycloakService.logout();
            this._user = null;
        } catch (err) {
            this._error = err instanceof Error ? err.message : 'Logout failed';
        } finally {
            this._isLoading = false;
        }
    }

    /**
     * Get valid access token for API calls
     */
    async getAccessToken(): Promise<string | null> {
        return keycloakService.getValidToken();
    }

    /**
     * Check if user has role
     */
    hasRole(role: string): boolean {
        return this._user?.hasRole(role) ?? false;
    }

    /**
     * Open account management page
     */
    async manageAccount(): Promise<void> {
        await keycloakService.accountManagement();
    }

    /**
     * Clear any errors
     */
    clearError(): void {
        this._error = null;
    }
}

// Export singleton instance
export const authState = new AuthState();
