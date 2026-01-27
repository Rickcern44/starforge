/**
 * Keycloak Configuration
 *
 * These values should match your Keycloak server setup.
 * For production, use environment variables.
 */
export const keycloakConfig = {
    url: import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:8080',
    realm: import.meta.env.VITE_KEYCLOAK_REALM || 'bouncy',
    clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'bouncy',
    clientSecret: import.meta.env.VITE_KEYCLOAK_CLIENT_SECRET
};

/**
 * Keycloak initialization options
 */
export const keycloakInitOptions = {
    // Check for existing session on load
    onLoad: 'check-sso' as const,

    // Use PKCE for enhanced security (recommended for SPAs)
    pkceMethod: 'S256' as const,

    // Silent SSO check using iframe
    silentCheckSsoRedirectUri: typeof window !== 'undefined'
        ? `${window.location.origin}/silent-check-sso.html`
        : undefined,

    // Check login iframe - useful for session management
    checkLoginIframe: true,
    checkLoginIframeInterval: 5
};

/**
 * Authentication methods available
 * These correspond to Keycloak authentication flows
 */
export enum AuthMethod {
    /** Standard username/password login */
    PASSWORD = 'password',

    /** Passwordless via WebAuthn/Passkeys */
    PASSWORDLESS = 'webauthn',

    /** Magic link sent to email */
    MAGIC_LINK = 'magic-link',

    /** Social login providers */
    GOOGLE = 'google',
    FACEBOOK = 'facebook',
    APPLE = 'apple'
}

/**
 * Routes configuration
 */
export const authRoutes = {
    signIn: '/account/signin',
    register: '/account/register',
    afterLogin: '/',
    afterLogout: '/account/signin'
};
