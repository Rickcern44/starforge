<script lang="ts">
    import { authState } from '$lib/auth';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    let email = $state('');
    let activeTab = $state<'password' | 'passwordless' | 'magic-link'>('password');
    let magicLinkSent = $state(false);

    onMount(async () => {
        // If already authenticated, redirect to home
        await authState.initialize();
        if (authState.isAuthenticated) {
            goto('/');
        }
    });

    async function handlePasswordLogin() {
        await authState.login();
    }

    async function handlePasswordlessLogin() {
        await authState.loginPasswordless();
    }

    async function handleMagicLink() {
        if (!email) return;
        await authState.requestMagicLink(email);
        magicLinkSent = true;
    }

    async function handleSocialLogin(provider: 'google' | 'facebook' | 'apple') {
        await authState.loginWithProvider(provider);
    }
</script>

<div class="min-h-screen flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
        <!-- Logo & Header -->
        <div class="text-center mb-8">
            <div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-emerald-500/20 to-sky-500/20 rounded-2xl flex items-center justify-center">
                <svg class="w-8 h-8 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="12" r="10" stroke-width="2"/>
                    <path stroke-width="2" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10"/>
                </svg>
            </div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-slate-100 to-slate-300 bg-clip-text text-transparent">
                Welcome Back
            </h1>
            <p class="text-slate-400 mt-2">Sign in to Pickup Manager</p>
        </div>

        <!-- Auth Card -->
        <div class="bg-slate-900/60 border border-slate-800 rounded-xl p-6 shadow-sm shadow-slate-950/60">
            <!-- Error Display -->
            {#if authState.error}
                <div class="mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
                    {authState.error}
                    <button 
                        onclick={() => authState.clearError()}
                        class="ml-2 underline hover:no-underline"
                    >
                        Dismiss
                    </button>
                </div>
            {/if}

            <!-- Auth Method Tabs -->
            <div class="flex gap-1 p-1 bg-slate-800/50 rounded-xl mb-6">
                <button
                    onclick={() => activeTab = 'password'}
                    class={`flex-1 py-2.5 px-3 text-sm font-medium rounded-lg transition-all ${
                        activeTab === 'password'
                            ? 'bg-slate-700 text-slate-100 shadow-lg'
                            : 'text-slate-400 hover:text-slate-200'
                    }`}
                >
                    Password
                </button>
                <button
                    onclick={() => activeTab = 'passwordless'}
                    class={`flex-1 py-2.5 px-3 text-sm font-medium rounded-lg transition-all ${
                        activeTab === 'passwordless'
                            ? 'bg-slate-700 text-slate-100 shadow-lg'
                            : 'text-slate-400 hover:text-slate-200'
                    }`}
                >
                    Passkey
                </button>
                <button
                    onclick={() => { activeTab = 'magic-link'; magicLinkSent = false; }}
                    class={`flex-1 py-2.5 px-3 text-sm font-medium rounded-lg transition-all ${
                        activeTab === 'magic-link'
                            ? 'bg-slate-700 text-slate-100 shadow-lg'
                            : 'text-slate-400 hover:text-slate-200'
                    }`}
                >
                    Magic Link
                </button>
            </div>

            <!-- Password Login -->
            {#if activeTab === 'password'}
                <div class="space-y-4">
                    <p class="text-sm text-slate-400 mb-4">
                        Sign in with your email and password through our secure authentication.
                    </p>
                    <button
                        onclick={handlePasswordLogin}
                        disabled={authState.isLoading}
                        class="w-full py-2.5 px-4 bg-emerald-500/90 hover:bg-emerald-500 disabled:bg-emerald-500/50 
                               text-slate-950 font-semibold rounded-xl shadow-lg shadow-emerald-500/25 
                               transition-all disabled:cursor-not-allowed
                               flex items-center justify-center gap-2"
                    >
                        {#if authState.isLoading}
                            <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                            </svg>
                        {:else}
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                            </svg>
                        {/if}
                        Continue with Email
                    </button>
                </div>
            {/if}

            <!-- Passwordless Login -->
            {#if activeTab === 'passwordless'}
                <div class="space-y-4">
                    <p class="text-sm text-slate-400 mb-4">
                        Use your device's biometrics or security key for a password-free sign in.
                    </p>
                    <button
                        onclick={handlePasswordlessLogin}
                        disabled={authState.isLoading}
                        class="w-full py-2.5 px-4 bg-emerald-500/90 hover:bg-emerald-500 disabled:bg-emerald-500/50 
                               text-slate-950 font-semibold rounded-xl shadow-lg shadow-emerald-500/25 
                               transition-all disabled:cursor-not-allowed
                               flex items-center justify-center gap-2"
                    >
                        {#if authState.isLoading}
                            <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                            </svg>
                        {:else}
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.39-2.823 1.07-4"/>
                            </svg>
                        {/if}
                        Sign in with Passkey
                    </button>
                    <p class="text-xs text-slate-500 text-center">
                        Supports Face ID, Touch ID, Windows Hello, and security keys
                    </p>
                </div>
            {/if}

            <!-- Magic Link Login -->
            {#if activeTab === 'magic-link'}
                <div class="space-y-4">
                    {#if magicLinkSent}
                        <div class="text-center py-4">
                            <div class="w-16 h-16 mx-auto mb-4 bg-emerald-500/10 rounded-full flex items-center justify-center">
                                <svg class="w-8 h-8 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                                </svg>
                            </div>
                            <h3 class="text-lg font-semibold text-slate-100 mb-2">Check your email</h3>
                            <p class="text-sm text-slate-400">
                                We've sent a magic link to <span class="text-slate-200">{email}</span>
                            </p>
                            <button
                                onclick={() => magicLinkSent = false}
                                class="mt-4 text-sm text-emerald-400 hover:text-emerald-300"
                            >
                                Use a different email
                            </button>
                        </div>
                    {:else}
                        <p class="text-sm text-slate-400 mb-4">
                            Enter your email and we'll send you a magic link to sign in instantly.
                        </p>
                        <div>
                            <label for="email" class="block text-sm font-medium text-slate-300 mb-2">
                                Email address
                            </label>
                            <input
                                id="email"
                                type="email"
                                bind:value={email}
                                placeholder="you@example.com"
                                class="w-full px-4 py-2.5 bg-slate-800/50 border border-slate-700 rounded-xl 
                                       text-slate-200 placeholder-slate-500 
                                       focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500/50 
                                       transition-all"
                            />
                        </div>
                        <button
                            onclick={handleMagicLink}
                            disabled={!email || authState.isLoading}
                            class="w-full py-2.5 px-4 bg-emerald-500/90 hover:bg-emerald-500 disabled:bg-emerald-500/50 
                                   text-slate-950 font-semibold rounded-xl shadow-lg shadow-emerald-500/25 
                                   transition-all disabled:cursor-not-allowed
                                   flex items-center justify-center gap-2"
                        >
                            {#if authState.isLoading}
                                <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                                </svg>
                            {:else}
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                                </svg>
                            {/if}
                            Send Magic Link
                        </button>
                    {/if}
                </div>
            {/if}

            <!-- Divider -->
            <div class="relative my-6">
                <div class="absolute inset-0 flex items-center">
                    <div class="w-full border-t border-slate-700"></div>
                </div>
                <div class="relative flex justify-center text-sm">
                    <span class="px-4 bg-slate-900/60 text-slate-500">or continue with</span>
                </div>
            </div>

            <!-- Social Login Options -->
            <div class="grid grid-cols-3 gap-3">
                <button
                    onclick={() => handleSocialLogin('google')}
                    disabled={authState.isLoading}
                    class="flex items-center justify-center py-2.5 px-4 bg-slate-800/50 hover:bg-slate-800 
                           border border-slate-700 rounded-xl transition-all disabled:opacity-50"
                >
                    <svg class="w-5 h-5" viewBox="0 0 24 24">
                        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                    </svg>
                </button>
                <button
                    onclick={() => handleSocialLogin('facebook')}
                    disabled={authState.isLoading}
                    class="flex items-center justify-center py-2.5 px-4 bg-slate-800/50 hover:bg-slate-800 
                           border border-slate-700 rounded-xl transition-all disabled:opacity-50"
                >
                    <svg class="w-5 h-5 text-[#1877F2]" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
                    </svg>
                </button>
                <button
                    onclick={() => handleSocialLogin('apple')}
                    disabled={authState.isLoading}
                    class="flex items-center justify-center py-2.5 px-4 bg-slate-800/50 hover:bg-slate-800 
                           border border-slate-700 rounded-xl transition-all disabled:opacity-50"
                >
                    <svg class="w-5 h-5 text-slate-100" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12.152 6.896c-.948 0-2.415-1.078-3.96-1.04-2.04.027-3.91 1.183-4.961 3.014-2.117 3.675-.546 9.103 1.519 12.09 1.013 1.454 2.208 3.09 3.792 3.039 1.52-.065 2.09-.987 3.935-.987 1.831 0 2.35.987 3.96.948 1.637-.026 2.676-1.48 3.676-2.948 1.156-1.688 1.636-3.325 1.662-3.415-.039-.013-3.182-1.221-3.22-4.857-.026-3.04 2.48-4.494 2.597-4.559-1.429-2.09-3.623-2.324-4.39-2.376-2-.156-3.675 1.09-4.61 1.09zM15.53 3.83c.843-1.012 1.4-2.427 1.245-3.83-1.207.052-2.662.805-3.532 1.818-.78.896-1.454 2.338-1.273 3.714 1.338.104 2.715-.688 3.559-1.701"/>
                    </svg>
                </button>
            </div>
        </div>

        <!-- Register Link -->
        <p class="text-center mt-6 text-slate-400">
            Don't have an account?
            <a href="/account/register" class="text-emerald-400 hover:text-emerald-300 font-medium">
                Create one
            </a>
        </p>
    </div>
</div>
