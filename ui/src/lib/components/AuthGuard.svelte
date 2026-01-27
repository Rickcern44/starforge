<script lang="ts">
    import { authState } from '$lib/auth';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import type { Snippet } from 'svelte';

    interface Props {
        /** Require authentication to view content */
        requireAuth?: boolean;
        /** Required roles (user must have at least one) */
        roles?: string[];
        /** Redirect to this path if unauthorized */
        redirectTo?: string;
        /** Content to show when authorized */
        children: Snippet;
        /** Optional loading content */
        loading?: Snippet;
        /** Optional unauthorized content */
        unauthorized?: Snippet;
    }

    let { 
        requireAuth = true, 
        roles = [], 
        redirectTo = '/account/signin',
        children,
        loading,
        unauthorized
    }: Props = $props();

    let isAuthorized = $state(false);
    let isChecking = $state(true);

    onMount(async () => {
        await authState.initialize();
        checkAuthorization();
    });

    function checkAuthorization() {
        isChecking = false;

        // Check if authentication is required
        if (requireAuth && !authState.isAuthenticated) {
            isAuthorized = false;
            
            // Redirect if no custom unauthorized content
            if (!unauthorized) {
                goto(redirectTo);
            }
            return;
        }

        // Check roles if specified
        if (roles.length > 0) {
            const hasRequiredRole = roles.some(role => authState.hasRole(role));
            
            if (!hasRequiredRole) {
                isAuthorized = false;
                return;
            }
        }

        isAuthorized = true;
    }

    // React to auth state changes
    $effect(() => {
        if (authState.isInitialized) {
            checkAuthorization();
        }
    });

    $inspect(authState)
</script>

{#if isChecking || authState.isLoading}
    {#if loading}
        {@render loading()}
    {:else}
        <div class="flex items-center justify-center min-h-[200px]">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
        </div>
    {/if}
{:else if isAuthorized}
    {@render children()}
{:else if unauthorized}
    {@render unauthorized()}
{/if}
