<script lang="ts">
    import './layout.css';
    import { page } from "$app/state";
    import { authState } from '$lib/auth';
    import { onMount } from 'svelte';
    import UserMenu from '$lib/components/UserMenu.svelte';

    $: currentPath = page.url.pathname;

    onMount(async () => {
        // Initialize auth on app startup
        await authState.initialize();
    });
</script>

<div class="min-h-screen bg-slate-950 text-slate-50">
    {#if currentPath.startsWith("/account")}
        <!-- No nav for account pages -->
    {:else}
        <!-- Top Navigation -->
        <nav class="sticky top-0 z-50 border-b border-slate-800 bg-slate-900/95 backdrop-blur-xl supports-[backdrop-filter:blur(20px)]:bg-slate-900/90">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div class="flex h-16 items-center justify-between">
                    <!-- Logo -->
                    <a href="/" class="flex items-center gap-2 group">
                        <div class="w-10 h-10 bg-gradient-to-br from-emerald-500/20 to-sky-500/20 rounded-xl flex items-center justify-center group-hover:scale-105 transition-all">
                            <svg class="w-6 h-6 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <circle cx="12" cy="12" r="10" stroke-width="2"/>
                            </svg>
                        </div>
                        <span class="text-xl font-bold bg-gradient-to-r from-emerald-400 via-sky-400 to-indigo-400 bg-clip-text text-transparent hidden md:inline">
                            Pickup Manager
                        </span>
                    </a>

                    <!-- Desktop Navigation -->
                    <div class="hidden md:flex items-center gap-6">
                        <a href="/" class={`px-3 py-2 text-sm font-medium rounded-xl transition-all ${
                            currentPath === '/'
                                ? 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40'
                                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'
                        }`}>
                            Dashboard
                        </a>
                        <a href="/leagues" class={`px-3 py-2 text-sm font-medium rounded-xl transition-all ${
                            currentPath.startsWith('/leagues')
                                ? 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40'
                                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'
                        }`}>
                            Leagues
                        </a>
                        <a href="/games" class={`px-3 py-2 text-sm font-medium rounded-xl transition-all ${
                            currentPath.startsWith('/games')
                                ? 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40'
                                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'
                        }`}>
                            All Games
                        </a>
                    </div>

                    <!-- Right side actions -->
                    <div class="flex items-center gap-3">
                        <!-- Search -->
                        <div class="hidden lg:block">
                            <div class="relative">
                                <span class="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                    <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                                    </svg>
                                </span>
                                <input
                                    type="text"
                                    placeholder="Search games..."
                                    class="w-64 pl-10 pr-4 py-2 bg-slate-800/50 border border-slate-700 rounded-xl text-sm text-slate-200 placeholder-slate-500 focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500/50 transition-all"
                                />
                            </div>
                        </div>

                        {#if authState.isAuthenticated}
                            <!-- Notifications (only for authenticated users) -->
                            <button class="p-2 rounded-xl text-slate-400 hover:bg-slate-800 hover:text-slate-200 relative">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
                                </svg>
                                <span class="absolute -top-1 -right-1 w-3 h-3 bg-rose-500 rounded-full"/>
                            </button>

                            <!-- New Game CTA -->
                            <a
                                href="/games/new"
                                class="hidden md:inline-flex items-center gap-2 bg-emerald-500/90 hover:bg-emerald-500 text-slate-950 px-4 py-2.5 rounded-xl text-sm font-semibold shadow-lg shadow-emerald-500/25 transition-all"
                            >
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                                </svg>
                                New Game
                            </a>

                            <!-- User Menu -->
                            <UserMenu />
                        {:else if !authState.isLoading}
                            <!-- Sign In / Register for unauthenticated users -->
                            <a
                                href="/account/signin"
                                class="px-4 py-2 text-sm font-medium text-slate-300 hover:text-slate-100 transition-colors"
                            >
                                Sign In
                            </a>
                            <a
                                href="/account/register"
                                class="hidden sm:inline-flex items-center gap-2 bg-emerald-500/90 hover:bg-emerald-500 text-slate-950 px-4 py-2.5 rounded-xl text-sm font-semibold shadow-lg shadow-emerald-500/25 transition-all"
                            >
                                Get Started
                            </a>
                        {:else}
                            <!-- Loading state -->
                            <div class="w-9 h-9 rounded-full bg-slate-800 animate-pulse"></div>
                        {/if}

                        <!-- Mobile menu -->
                        <button class="md:hidden p-2 rounded-xl text-slate-400 hover:bg-slate-800 hover:text-slate-200">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </nav>
    {/if}

    <!-- Main content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 lg:py-12">
        <slot/>
    </main>
</div>
