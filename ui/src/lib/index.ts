// Auth exports
export { authState, keycloakService, AuthMethod, authRoutes } from './auth';

// Component exports
export { default as AuthGuard } from './components/AuthGuard.svelte';
export { default as UserMenu } from './components/UserMenu.svelte';
export { default as Card } from './components/Card.svelte';
export { default as GameCard } from './components/GameCard.svelte';
export { default as GamesGrid } from './components/GamesGrid.svelte';

// State exports
export { useUser } from './state/user.svelte';

// Domain exports
export { User } from './domain/user';

// API exports
export { api } from './api/client';
