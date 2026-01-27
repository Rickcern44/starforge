<script lang="ts">
    import { authState } from '$lib/auth';

    let isOpen = $state(false);

    function toggleMenu() {
        isOpen = !isOpen;
    }

    function closeMenu() {
        isOpen = false;
    }

    async function handleLogout() {
        closeMenu();
        await authState.logout();
    }

    async function handleManageAccount() {
        closeMenu();
        await authState.manageAccount();
    }

    // Get user initials for avatar
    const initials = $derived(() => {
        const name = authState.user?.name || authState.user?.email || '';
        const parts = name.split(' ');
        if (parts.length >= 2) {
            return (parts[0][0] + parts[1][0]).toUpperCase();
        }
        return name.substring(0, 2).toUpperCase();
    });
</script>

<svelte:window onclick={(e) => {
    if (isOpen && !(e.target as HTMLElement).closest('.user-menu')) {
        closeMenu();
    }
}} />

<div class="relative user-menu">
    <button
        onclick={toggleMenu}
        class="flex items-center gap-2 p-2 rounded-xl hover:bg-slate-800 text-slate-300 transition-all"
    >
        <div class="w-9 h-9 bg-gradient-to-r from-emerald-500/20 to-sky-500/20 rounded-full flex items-center justify-center">
            <span class="font-semibold text-sm">{initials()}</span>
        </div>
        <span class="hidden md:inline text-sm font-medium max-w-[120px] truncate">
            {authState.user?.name || authState.user?.email || 'User'}
        </span>
        <svg 
            class="w-4 h-4 text-slate-400 transition-transform" 
            class:rotate-180={isOpen}
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
        >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
        </svg>
    </button>

    {#if isOpen}
        <div class="absolute right-0 mt-2 w-64 bg-slate-900/95 backdrop-blur-xl border border-slate-800 rounded-xl shadow-sm shadow-slate-950/60 overflow-hidden z-50">
            <!-- User Info -->
            <div class="p-4 border-b border-slate-800">
                <p class="font-medium text-sm text-slate-200 truncate">
                    {authState.user?.name || 'User'}
                </p>
                <p class="text-xs text-slate-500 truncate">
                    {authState.user?.email || ''}
                </p>
            </div>

            <!-- Menu Items -->
            <div class="py-1">
                <a
                    href="/profile"
                    onclick={closeMenu}
                    class="flex items-center gap-3 px-4 py-2 text-sm text-slate-300 hover:bg-slate-800 hover:text-slate-100 transition-all"
                >
                    <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                    </svg>
                    My Profile
                </a>
                <a
                    href="/settings"
                    onclick={closeMenu}
                    class="flex items-center gap-3 px-4 py-2 text-sm text-slate-300 hover:bg-slate-800 hover:text-slate-100 transition-all"
                >
                    <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                    Settings
                </a>
                <button
                    onclick={handleManageAccount}
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-slate-300 hover:bg-slate-800 hover:text-slate-100 transition-all"
                >
                    <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                    </svg>
                    Manage Account
                </button>
            </div>

            <!-- Logout -->
            <div class="border-t border-slate-800 py-1">
                <button
                    onclick={handleLogout}
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-rose-400 hover:bg-slate-800 hover:text-rose-300 transition-all"
                >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
                    </svg>
                    Sign Out
                </button>
            </div>
        </div>
    {/if}
</div>
