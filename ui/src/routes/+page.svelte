<script lang="ts">
    import GamesGrid from '$lib/components/GamesGrid.svelte';
    import type {UpcomingGame} from '$lib/types/gameTypes';
    import {getUpcomingGames} from '$lib/api/data';
    import {onMount} from 'svelte';
    import {AuthGuard} from "$lib";

    let upcomingGames: UpcomingGame[] = [];
    let loading = true;
    let weekGroups: { title: string; games: UpcomingGame[] }[] = [];

    onMount(() => {
        upcomingGames = getUpcomingGames();
        loading = false;

        // Group games by week after data loads
        if (upcomingGames.length > 0) {
            groupGamesByWeek();
        }
    });

    $: if (upcomingGames.length > 0 && !loading) {
        groupGamesByWeek();
    }

    function groupGamesByWeek() {
        const weeks: Record<string, { startDate: Date; games: UpcomingGame[] }> = {};

        upcomingGames.forEach((game) => {
            const date = new Date(game.game.startTime);
            const startOfWeek = new Date(date);
            startOfWeek.setDate(date.getDate() - date.getDay()); // Start of week (Sunday)
            startOfWeek.setHours(0, 0, 0, 0);

            const weekKey = startOfWeek.toISOString().split('T')[0];

            if (!weeks[weekKey]) {
                weeks[weekKey] = {
                    startDate: startOfWeek,
                    games: []
                };
            }
            weeks[weekKey].games.push(game);
        });

        weekGroups = Object.values(weeks)
            .sort((a, b) => a.startDate.getTime() - b.startDate.getTime())
            .map(group => ({
                title: formatWeekTitle(group.startDate),
                games: group.games.sort((a, b) => new Date(a.game.startTime).getTime() - new Date(b.game.startTime).getTime())
            }));
    }

    function formatWeekTitle(date: Date): string {
        const now = new Date();
        const diffTime = date.getTime() - now.getTime();
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays <= 0) return 'This Week';
        if (diffDays < 7) return 'Next Week';
        if (diffDays < 14) return 'Week of ' + date.toLocaleDateString('en-US', {month: 'short', day: 'numeric'});

        return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric'
            }) + ' - ' +
            new Date(date.getTime() + 6 * 86400000).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
    }
</script>

<AuthGuard>
    <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <div>
                <h1 class="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-slate-100 to-slate-300 bg-clip-text text-transparent">
                    Upcoming Games
                </h1>
                <p class="text-slate-400 mt-2 text-sm">Your schedule across all leagues</p>
            </div>

            <div class="flex items-center gap-3">
                <button class="px-4 py-2.5 bg-slate-800/50 hover:bg-slate-800 border border-slate-700 rounded-xl text-sm font-medium text-slate-300 transition-all">
                    This Month
                </button>
                <button class="px-4 py-2.5 bg-emerald-500/90 hover:bg-emerald-500 text-slate-950 rounded-xl text-sm font-semibold shadow-lg shadow-emerald-500/25 transition-all">
                    This Week
                </button>
            </div>
        </div>

        {#if loading}
            <!-- Loading state -->
            <div class="flex items-center justify-center py-24">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-500"></div>
            </div>
        {:else if weekGroups.length === 0}
            <!-- Empty state -->
            <div class="text-center py-24">
                <div class="w-24 h-24 mx-auto mb-6 bg-slate-800/50 rounded-3xl flex items-center justify-center">
                    <span class="i-lucide-calendar-x w-12 h-12 text-slate-600"/>
                </div>
                <h3 class="text-lg font-semibold text-slate-200 mb-2">No games scheduled</h3>
                <p class="text-sm text-slate-400 mb-6 max-w-md mx-auto">
                    Join a league or create a new game to get started. Your schedule will appear here.
                </p>
                <a href="/leagues"
                   class="inline-flex items-center gap-2 bg-emerald-500/90 hover:bg-emerald-500 text-slate-950 px-4 py-2.5 rounded-xl font-semibold shadow-lg shadow-emerald-500/25 transition-all">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                    </svg>
                    Find Leagues
                </a>
            </div>
        {:else}
            <!-- Weekly groups -->
            {#each weekGroups as week}
                <section class="space-y-4">
                    <!-- Week header -->
                    <div class="flex items-center gap-4">
                        <div class="w-2 h-12 bg-gradient-to-b from-emerald-400 to-emerald-500 rounded-full"/>
                        <div>
                            <h2 class="text-lg font-semibold text-slate-100">{week.title}</h2>
                            <p class="text-xs text-slate-500">{week.games.length}
                                game{week.games.length === 1 ? '' : 's'}</p>
                        </div>
                    </div>

                    <!-- Games grid for this week -->
                    <GamesGrid games={week.games}/>
                </section>
            {/each}
        {/if}
    </div>
</AuthGuard>