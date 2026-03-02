<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import type { League, Game } from '$lib/models';
  import { goto } from '$app/navigation';

  let { data } = $props();

  let allLeagues = $state<League[]>([]);
  let selectedLeagueName = $state<string>('All Leagues');

  let user = $derived(authService.user);

  let managedLeagues = $derived.by(() => {
    if (!user) return [];
    return allLeagues.filter(league => {
      const member = league.members?.find(m => m.playerId === user.id);
      return member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'));
    });
  });

  $effect(() => {
    if (data.leagues) {
      allLeagues = data.leagues;
    }
  });

  let upcomingGames = $derived.by(() => {
    let games: Game[] = [];
    const now = new Date();
    // Set to start of today to ensure today's evening games are included
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
    const oneWeekFromNow = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000).getTime();
    
    if (selectedLeagueName === 'All Leagues') {
      allLeagues.forEach((league) => {
        if (league.games) {
          games.push(...league.games);
        }
      });
    } else {
      const selected = allLeagues.find((l) => l.name === selectedLeagueName);
      if (selected && selected.games) {
        games.push(...selected.games);
      }
    }
    
    return games
      .filter(g => {
        const gTime = g.startTime instanceof Date ? g.startTime.getTime() : new Date(g.startTime).getTime();
        return gTime >= startOfToday && gTime <= oneWeekFromNow;
      })
      .sort((a, b) => {
        const aTime = a.startTime instanceof Date ? a.startTime.getTime() : new Date(a.startTime).getTime();
        const bTime = b.startTime instanceof Date ? b.startTime.getTime() : new Date(b.startTime).getTime();
        return aTime - bTime;
      });
  });

  let financialSummary = $derived.by(() => {
    let balanceInCents = 0;
    let unpaidCount = 0;

    if (data.charges) {
      data.charges.forEach((c: any) => {
        const allocated = c.allocations?.reduce((sum: number, a: any) => sum + a.amountInCents, 0) || 0;
        if (allocated < c.amountCents) {
          balanceInCents += (c.amountCents - allocated);
          unpaidCount++;
        }
      });
    }

    return {
      balance: balanceInCents / 100,
      count: unpaidCount
    };
  });

  let canCreateEvent = $derived.by(() => {
    if (!user) return false;
    for (const league of allLeagues) {
      if (league.members) {
        for (const member of league.members) {
          if (member.playerId === user.id) {
            const role = member.role.toLowerCase();
            if (role.includes('admin') || role.includes('owner')) return true;
          }
        }
      }
    }
    return false;
  });

</script>

<svelte:head>
  <title>Home | League Manager</title>
</svelte:head>

<div class="max-w-3xl mx-auto space-y-8 py-4">
    <header class="flex justify-between items-end">
      <div>
        <p class="text-sm font-medium text-gray-500 uppercase tracking-wider">Welcome back,</p>
        <h2 class="text-3xl font-bold text-gray-900">{user ? user.name : 'Guest'}</h2>
      </div>
      <button class="p-2 text-gray-400 hover:text-indigo-600 transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </button>
    </header>

    <!-- Managed Leagues Section (Intuitive Admin Entry) -->
    {#if managedLeagues.length > 0}
      <section class="space-y-4">
        <div class="flex justify-between items-center">
          <h3 class="text-xl font-black text-gray-900">Manage Your Leagues</h3>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          {#each managedLeagues as league}
            <button
              class="flex items-center justify-between p-4 bg-indigo-50 border border-indigo-100 rounded-2xl hover:bg-indigo-100 transition-all group"
              onclick={() => goto(`/leagues/${league.id}/admin`)}
            >
              <div class="flex items-center space-x-3">
                <div class="bg-indigo-600 p-2 rounded-xl text-white">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37a1.724 1.724 0 002.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
                <div class="text-left">
                  <p class="font-black text-indigo-900 text-sm">{league.name}</p>
                  <p class="text-[10px] font-bold text-indigo-400 uppercase tracking-widest">Admin Dashboard</p>
                </div>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-400 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          {/each}
        </div>
      </section>
    {/if}

    <!-- League Selector (Horizontal Chips) -->
    <section class="space-y-3">
      <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider">Select League</h3>
      <div class="flex space-x-2 overflow-x-auto pb-2 scrollbar-hide no-scrollbar">
        {#each ['All Leagues', ...allLeagues.map(l => l.name)] as name}
          <button
            class="px-5 py-2 rounded-full text-sm font-semibold whitespace-nowrap transition-all duration-200 border"
            class:bg-indigo-600={selectedLeagueName === name}
            class:text-white={selectedLeagueName === name}
            class:border-indigo-600={selectedLeagueName === name}
            class:bg-white={selectedLeagueName !== name}
            class:text-gray-600={selectedLeagueName !== name}
            class:border-gray-200={selectedLeagueName !== name}
            class:hover:border-indigo-300={selectedLeagueName !== name}
            onclick={() => (selectedLeagueName = name)}
          >
            {name}
          </button>
        {/each}
      </div>
    </section>

    <!-- Financial Summary (Minimalist) -->
    <button 
      class="w-full text-left bg-indigo-600 text-white p-8 rounded-[20px] shadow-lg shadow-indigo-200 cursor-pointer hover:bg-indigo-700 transition-all duration-300 transform hover:-translate-y-1" 
      onclick={() => goto('/ledger')}
    >
      <div class="flex justify-between items-start mb-8">
        <div>
          <p class="text-sm font-medium text-indigo-100 opacity-80 mb-1">Outstanding Balance</p>
          <p class="text-4xl font-bold tracking-tight">${financialSummary.balance.toFixed(2)}</p>
        </div>
        <div class="bg-white/20 p-3 rounded-xl backdrop-blur-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
      </div>
      
      <div class="flex items-center justify-between text-indigo-100">
        <div class="flex items-center space-x-3 text-sm font-medium">
          <div class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1 text-orange-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Unpaid Charges: {financialSummary.count} {financialSummary.count === 1 ? 'Item' : 'Items'}
          </div>
        </div>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 opacity-60" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      </div>
    </button>

    <!-- Upcoming Events -->
    <section class="space-y-4">
      <div class="flex justify-between items-center">
        <h3 class="text-xl font-bold text-gray-900">Upcoming Events</h3>
        <button class="text-indigo-600 hover:text-indigo-800 text-sm font-bold transition-colors" onclick={() => goto('/events')}>See all</button>
      </div>
      
      <div class="space-y-3">
        {#if upcomingGames.length === 0}
          <div class="bg-gray-50 border-2 border-dashed border-gray-200 rounded-2xl py-12 text-center">
            <p class="text-gray-400 font-medium italic">No upcoming events scheduled</p>
          </div>
        {:else}
          {#each upcomingGames as game (game.id)}
            {@const league = allLeagues.find(l => l.id === game.leagueId)}
            {@const member = league?.members.find(m => m.playerId === user?.id)}
            {@const isAdmin = member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'))}
            
            <div class="relative group">
              <button
                class="w-full text-left bg-white p-4 rounded-2xl shadow-sm border border-gray-100 cursor-pointer flex items-center space-x-4 hover:shadow-md hover:border-indigo-100 transition-all duration-200"
                onclick={() => goto(`/events/${game.id}`)}
              >
                <div class="bg-indigo-50 p-3 rounded-xl text-indigo-600 group-hover:bg-indigo-600 group-hover:text-white transition-colors duration-200">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 17l-3 3m0 0l-3-3m3 3V10" />
                  </svg>
                </div>
                <div class="flex-grow">
                  <p class="font-bold text-gray-900 leading-tight">Game at {game.location}</p>
                  <p class="text-sm text-gray-500 font-medium mt-0.5">
                    {league ? league.name : 'Unknown League'} • {new Date(game.startTime).toLocaleDateString('en-US', { month: 'short', day: 'numeric', timeZone: 'America/New_York' })}
                  </p>
                </div>
                <div class="text-right">
                  <p class="font-black text-indigo-600">{new Date(game.startTime).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/New_York' })}</p>
                </div>
              </button>
              
              {#if isAdmin}
                <button 
                  class="absolute -top-2 -right-2 bg-indigo-100 text-indigo-600 p-1.5 rounded-lg border border-white shadow-sm opacity-0 group-hover:opacity-100 transition-opacity hover:bg-indigo-600 hover:text-white"
                  onclick={(e) => { e.stopPropagation(); goto(`/leagues/${game.leagueId}/admin`); }}
                  title="Manage League"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37a1.724 1.724 0 002.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              {/if}
            </div>
          {/each}
        {/if}
      </div>
    </section>

    <!-- Floating Action Button -->
    {#if canCreateEvent}
      <button
        class="fixed bottom-8 right-8 bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3.5 px-6 rounded-2xl shadow-xl shadow-indigo-200 flex items-center space-x-2 transition-all duration-200 transform hover:scale-105 active:scale-95"
        onclick={() => goto('/events/create')}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        <span class="text-sm">New Event</span>
      </button>
    {/if}
  </div>

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .no-scrollbar {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }
</style>
