import Keycloak from 'keycloak-js';
import { keycloakConfig, keycloakInitOptions, AuthMethod, authRoutes } from './keycloak.config';
import { browser } from '$app/environment';

/**
 * Keycloak service singleton for managing authentication
 */
class KeycloakService {
    private keycloak: Keycloak | null = null;
    private initialized = false;
    private initPromise: Promise<boolean> | null = null;

    /**
     * Get the Keycloak instance
     */
    get instance(): Keycloak | null {
        return this.keycloak;
    }

    /**
     * Check if user is authenticated
     */
    get isAuthenticated(): boolean {
        return this.keycloak?.authenticated ?? false;
    }

    /**
     * Get current access token
     */
    get token(): string | undefined {
        return this.keycloak?.token;
    }

    /**
     * Get token parsed claims
     */
    get tokenParsed(): Keycloak.KeycloakTokenParsed | undefined {
        return this.keycloak?.tokenParsed;
    }

    /**
     * Get user profile from token
     */
    get userProfile(): { id: string; email: string; name: string; roles: string[] } | null {
        if (!this.keycloak?.tokenParsed) return null;
        
        const token = this.keycloak.tokenParsed;
        return {
            id: token.sub || '',
            email: token.email || '',
            name: token.name || token.preferred_username || '',
            roles: token.realm_access?.roles || []
        };
    }

    /**
     * Initialize Keycloak - should be called once on app startup
     */
    async init(): Promise<boolean> {
        if (!browser) return false;
        
        // Return existing promise if already initializing
        if (this.initPromise) {
            return this.initPromise;
        }

        // Return true if already initialized and authenticated
        if (this.initialized && this.keycloak) {
            return this.keycloak.authenticated ?? false;
        }

        this.initPromise = this._doInit();
        return this.initPromise;
    }

    private async _doInit(): Promise<boolean> {
        try {
            this.keycloak = new Keycloak(keycloakConfig);

            // Set up event handlers
            this.keycloak.onTokenExpired = () => {
                console.log('[Keycloak] Token expired, refreshing...');
                this.refreshToken();
            };

            this.keycloak.onAuthSuccess = () => {
                console.log('[Keycloak] Auth success');
            };

            this.keycloak.onAuthError = (error) => {
                console.error('[Keycloak] Auth error:', error);
            };

            this.keycloak.onAuthLogout = () => {
                console.log('[Keycloak] Logged out');
            };

            $inspect("Key Cloak Options: ", keycloakConfig);
            const authenticated = await this.keycloak.init(keycloakInitOptions);
            this.initialized = true;
            
            console.log('[Keycloak] Initialized, authenticated:', authenticated);
            return authenticated;
        } catch (error) {
            console.error('[Keycloak] Init failed:', error);
            this.initialized = true; // Mark as initialized even on failure
            return false;
        }
    }

    /**
     * Login with standard redirect
     */
    async login(options?: { 
        redirectUri?: string;
        idpHint?: string;
        loginHint?: string;
        action?: string;
    }): Promise<void> {
        if (!this.keycloak) {
            await this.init();
        }

        const defaultRedirectUri = browser 
            ? `${window.location.origin}${authRoutes.afterLogin}`
            : undefined;

        await this.keycloak?.login({
            redirectUri: options?.redirectUri || defaultRedirectUri,
            idpHint: options?.idpHint,
            loginHint: options?.loginHint,
            action: options?.action
        });
    }

    /**
     * Login with specific authentication method
     */
    async loginWithMethod(method: AuthMethod, email?: string): Promise<void> {
        if (!this.keycloak) {
            await this.init();
        }

        const redirectUri = browser 
            ? `${window.location.origin}${authRoutes.afterLogin}`
            : undefined;

        switch (method) {
            case AuthMethod.PASSWORD:
                // Standard login flow
                await this.keycloak?.login({ redirectUri });
                break;

            case AuthMethod.PASSWORDLESS:
                // Trigger WebAuthn/Passkey authentication
                // This requires Keycloak to be configured with WebAuthn authenticator
                await this.keycloak?.login({ 
                    redirectUri,
                    // Use custom authentication flow if configured
                    acr: { values: ['webauthn'], essential: true }
                } as any);
                break;

            case AuthMethod.MAGIC_LINK:
                // Magic link requires a custom Keycloak authenticator
                // Redirect to a custom endpoint or use action token
                if (email) {
                    // Store email for magic link flow
                    sessionStorage.setItem('magic_link_email', email);
                }
                await this.keycloak?.login({ 
                    redirectUri,
                    loginHint: email,
                    // Custom action for magic link if configured
                    action: 'magic-link'
                } as any);
                break;

            case AuthMethod.GOOGLE:
            case AuthMethod.FACEBOOK:
            case AuthMethod.APPLE:
                // Social login via Identity Provider hint
                await this.keycloak?.login({ 
                    redirectUri,
                    idpHint: method 
                });
                break;

            default:
                await this.keycloak?.login({ redirectUri });
        }
    }

    /**
     * Register new user
     */
    async register(options?: { 
        redirectUri?: string;
        loginHint?: string;
    }): Promise<void> {
        if (!this.keycloak) {
            await this.init();
        }

        const defaultRedirectUri = browser 
            ? `${window.location.origin}${authRoutes.afterLogin}`
            : undefined;

        await this.keycloak?.register({
            redirectUri: options?.redirectUri || defaultRedirectUri,
            loginHint: options?.loginHint
        });
    }

    /**
     * Logout user
     */
    async logout(redirectUri?: string): Promise<void> {
        if (!this.keycloak) return;

        const defaultRedirectUri = browser 
            ? `${window.location.origin}${authRoutes.afterLogout}`
            : undefined;

        await this.keycloak.logout({
            redirectUri: redirectUri || defaultRedirectUri
        });
    }

    /**
     * Refresh the access token
     */
    async refreshToken(minValidity = 30): Promise<boolean> {
        if (!this.keycloak) return false;

        try {
            const refreshed = await this.keycloak.updateToken(minValidity);
            if (refreshed) {
                console.log('[Keycloak] Token refreshed');
            }
            return refreshed;
        } catch (error) {
            console.error('[Keycloak] Token refresh failed:', error);
            // Token refresh failed, user needs to re-authenticate
            return false;
        }
    }

    /**
     * Get a valid token, refreshing if necessary
     */
    async getValidToken(): Promise<string | null> {
        if (!this.keycloak?.token) return null;

        // Check if token is expired or about to expire (within 30 seconds)
        const isExpired = this.keycloak.isTokenExpired(30);
        
        if (isExpired) {
            const refreshed = await this.refreshToken();
            if (!refreshed) {
                // Couldn't refresh, need to re-login
                return null;
            }
        }

        return this.keycloak.token;
    }

    /**
     * Check if user has a specific role
     */
    hasRole(role: string): boolean {
        return this.keycloak?.hasRealmRole(role) ?? false;
    }

    /**
     * Check if user has a specific resource role
     */
    hasResourceRole(role: string, resource?: string): boolean {
        return this.keycloak?.hasResourceRole(role, resource) ?? false;
    }

    /**
     * Account management - opens Keycloak account page
     */
    async accountManagement(): Promise<void> {
        if (!this.keycloak) return;
        await this.keycloak.accountManagement();
    }
}

// Export singleton instance
export const keycloakService = new KeycloakService();
