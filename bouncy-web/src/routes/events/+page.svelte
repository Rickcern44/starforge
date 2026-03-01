<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import type { League, Game } from '$lib/models';
  import { goto } from '$app/navigation';

  let { data } = $props();

  let allLeagues = $state<League[]>(data.leagues || []);
  let user = $derived(authService.user);

  let allGames = $derived.by(() => {
    let games: Game[] = [];
    allLeagues.forEach((league) => {
      if (league.games) {
        games.push(...league.games);
      }
    });
    return games.sort((a, b) => {
      const aTime = a.startTime instanceof Date ? a.startTime.getTime() : new Date(a.startTime).getTime();
      const bTime = b.startTime instanceof Date ? b.startTime.getTime() : new Date(b.startTime).getTime();
      return bTime - aTime; // Newest first for full list
    });
  });

  let upcomingGames = $derived(allGames.filter(g => {
    const gTime = g.startTime instanceof Date ? g.startTime.getTime() : new Date(g.startTime).getTime();
    return gTime >= new Date().getTime();
  }).reverse()); // Reverse to show soonest first

  let pastGames = $derived(allGames.filter(g => {
    const gTime = g.startTime instanceof Date ? g.startTime.getTime() : new Date(g.startTime).getTime();
    return gTime < new Date().getTime();
  }));

</script>

<svelte:head>
  <title>All Events | Bouncy</title>
</svelte:head>

<div class="max-w-3xl mx-auto space-y-8 py-4 px-4">
  <header class="flex items-center space-x-4">
    <button onclick={() => goto('/')} class="text-gray-500 hover:text-indigo-600 transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
    <h1 class="text-3xl font-black text-gray-900 tracking-tight leading-tight">All Events</h1>
  </header>

  {#if upcomingGames.length > 0}
    <section class="space-y-4">
      <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest px-1">Upcoming</h3>
      <div class="space-y-3">
        {#each upcomingGames as game (game.id)}
          {@const league = allLeagues.find(l => l.id === game.leagueId)}
          {@const member = league?.members.find(m => m.playerId === user?.id)}
          {@const isAdmin = member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'))}
          
          <button
            class="w-full text-left bg-white p-4 rounded-2xl shadow-sm border border-gray-100 cursor-pointer flex items-center space-x-4 hover:shadow-md hover:border-indigo-100 transition-all group"
            onclick={() => goto(`/events/${game.id}`)}
          >
            <div class="bg-indigo-50 p-3 rounded-xl text-indigo-600 group-hover:bg-indigo-600 group-hover:text-white transition-colors duration-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <div class="flex-grow">
              <p class="font-bold text-gray-900 leading-tight">Game at {game.location}</p>
              <p class="text-sm text-gray-500 font-medium mt-0.5">
                {league ? league.name : 'Unknown'} • {new Date(game.startTime).toLocaleDateString()}
              </p>
            </div>
            {#if isAdmin}
              <span class="px-2 py-1 bg-indigo-50 text-indigo-600 rounded-md text-[9px] font-black uppercase tracking-widest">Admin</span>
            {/if}
          </button>
        {/each}
      </div>
    </section>
  {/if}

  {#if pastGames.length > 0}
    <section class="space-y-4">
      <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest px-1">Past Events</h3>
      <div class="space-y-3">
        {#each pastGames as game (game.id)}
          {@const league = allLeagues.find(l => l.id === game.leagueId)}
          
          <button
            class="w-full text-left bg-gray-50 p-4 rounded-2xl border border-transparent cursor-pointer flex items-center space-x-4 hover:bg-white hover:border-gray-100 transition-all group opacity-80 hover:opacity-100"
            onclick={() => goto(`/events/${game.id}`)}
          >
            <div class="bg-gray-100 p-3 rounded-xl text-gray-400 group-hover:bg-gray-200 transition-colors duration-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="flex-grow">
              <p class="font-bold text-gray-600 leading-tight">Game at {game.location}</p>
              <p class="text-sm text-gray-400 font-medium mt-0.5">
                {league ? league.name : 'Unknown'} • {new Date(game.startTime).toLocaleDateString()}
              </p>
            </div>
          </button>
        {/each}
      </div>
    </section>
  {/if}
</div>
