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

<div class="max-w-xl mx-auto space-y-8 py-2">
    <header class="flex justify-between items-end px-1">
      <div>
        <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-1">Welcome back,</p>
        <h2 class="text-3xl font-black text-gray-900 tracking-tight">{user ? user.name : 'Guest'}</h2>
      </div>
      <button class="w-10 h-10 bg-white border border-gray-100 rounded-xl shadow-sm flex items-center justify-center text-gray-400 active:text-indigo-600 active:bg-indigo-50 transition-all">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </button>
    </header>

    <!-- Financial Summary (Premium Card) -->
    <button 
      class="w-full text-left bg-indigo-600 text-white p-6 rounded-[28px] shadow-2xl shadow-indigo-200 cursor-pointer hover:bg-indigo-700 active:scale-[0.97] transition-all duration-300 relative overflow-hidden group" 
      onclick={() => goto('/ledger')}
    >
      <!-- Background Ornament -->
      <div class="absolute -right-8 -top-8 w-32 h-32 bg-white/10 rounded-full blur-2xl group-hover:scale-150 transition-transform duration-700"></div>
      
      <div class="relative z-10 flex justify-between items-start mb-6">
        <div>
          <p class="text-xs font-bold text-indigo-100 uppercase tracking-widest mb-1 opacity-80">Outstanding Balance</p>
          <p class="text-5xl font-black tracking-tighter">${financialSummary.balance.toFixed(2)}</p>
        </div>
        <div class="bg-white/10 p-3 rounded-2xl backdrop-blur-md border border-white/20">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
      </div>
      
      <div class="relative z-10 flex items-center justify-between text-indigo-100">
        <div class="flex items-center space-x-2 bg-black/10 px-3 py-1.5 rounded-full text-[11px] font-black uppercase tracking-wider border border-white/10">
          <span class="w-2 h-2 rounded-full bg-orange-400 animate-pulse"></span>
          <span>{financialSummary.count} {financialSummary.count === 1 ? 'Charge' : 'Charges'} pending</span>
        </div>
        <div class="flex items-center text-xs font-bold uppercase tracking-widest opacity-70 group-hover:opacity-100 transition-opacity">
          Details
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </button>

    <!-- Managed Leagues Section -->
    {#if managedLeagues.length > 0}
      <section class="space-y-4 px-1">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-black text-gray-900 tracking-tight">Admin Dashboard</h3>
        </div>
        <div class="grid grid-cols-1 gap-3">
          {#each managedLeagues as league}
            <button
              class="flex items-center justify-between p-4 bg-white border border-gray-100 rounded-[22px] hover:border-indigo-100 shadow-sm active:scale-[0.98] active:bg-gray-50 transition-all group"
              onclick={() => goto(`/leagues/${league.id}/admin`)}
            >
              <div class="flex items-center space-x-4">
                <div class="bg-indigo-50 p-2.5 rounded-xl text-indigo-600">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37a1.724 1.724 0 002.572-1.065z" />
                  </svg>
                </div>
                <div class="text-left">
                  <p class="font-bold text-gray-900 leading-tight">{league.name}</p>
                  <p class="text-[10px] font-black text-indigo-500 uppercase tracking-widest mt-0.5">Settings & Finances</p>
                </div>
              </div>
              <div class="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center group-hover:bg-indigo-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400 group-hover:text-indigo-600 group-hover:translate-x-0.5 transition-all" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </button>
          {/each}
        </div>
      </section>
    {/if}

    <!-- League Selector -->
    <section class="space-y-4 px-1">
      <div class="flex space-x-2 overflow-x-auto pb-2 no-scrollbar">
        {#each ['All Leagues', ...allLeagues.map(l => l.name)] as name}
          <button
            class="px-6 py-2.5 rounded-2xl text-[13px] font-black whitespace-nowrap transition-all duration-300 border uppercase tracking-wide"
            class:bg-indigo-600={selectedLeagueName === name}
            class:text-white={selectedLeagueName === name}
            class:border-indigo-600={selectedLeagueName === name}
            class:shadow-lg={selectedLeagueName === name}
            class:shadow-indigo-100={selectedLeagueName === name}
            class:bg-white={selectedLeagueName !== name}
            class:text-gray-400={selectedLeagueName !== name}
            class:border-gray-100={selectedLeagueName !== name}
            class:active:bg-gray-50={selectedLeagueName !== name}
            onclick={() => (selectedLeagueName = name)}
          >
            {name}
          </button>
        {/each}
      </div>
    </section>

    <!-- Upcoming Events -->
    <section class="space-y-5 px-1 pb-20">
      <div class="flex justify-between items-center">
        <h3 class="text-xl font-black text-gray-900 tracking-tight">Schedule</h3>
        <button class="text-[11px] font-black text-indigo-600 uppercase tracking-widest border-b-2 border-indigo-100 pb-0.5" onclick={() => goto('/events')}>View All</button>
      </div>
      
      <div class="space-y-4">
        {#if upcomingGames.length === 0}
          <div class="bg-gray-50 border-2 border-dashed border-gray-200 rounded-[28px] py-16 text-center">
            <p class="text-gray-400 font-bold italic text-sm">No games on the horizon</p>
          </div>
        {:else}
          {#each upcomingGames as game (game.id)}
            {@const league = allLeagues.find(l => l.id === game.leagueId)}
            {@const member = league?.members.find(m => m.playerId === user?.id)}
            {@const isAdmin = member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'))}
            
            <button
              class="w-full text-left bg-white p-5 rounded-[28px] border border-gray-100 shadow-sm active:scale-[0.98] active:bg-gray-50 transition-all duration-300 flex items-center group relative overflow-hidden"
              onclick={() => goto(`/events/${game.id}`)}
            >
              <div class="bg-gray-50 p-3 rounded-2xl text-gray-400 group-active:bg-indigo-600 group-active:text-white transition-colors duration-200">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <div class="ml-4 flex-grow">
                <div class="flex items-center space-x-2">
                  <span class="text-[10px] font-black text-indigo-500 uppercase tracking-widest">{league?.name || 'League'}</span>
                  <span class="w-1 h-1 rounded-full bg-gray-300"></span>
                  <span class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">{new Date(game.startTime).toLocaleDateString('en-US', { weekday: 'short' })}</span>
                </div>
                <p class="font-black text-gray-900 text-lg leading-tight mt-0.5 group-hover:text-indigo-600 transition-colors">{game.location}</p>
              </div>
              <div class="text-right flex flex-col items-end">
                <p class="font-black text-gray-900 text-sm">{new Date(game.startTime).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/New_York' })}</p>
                <div class="mt-1 flex items-center text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  Details
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 ml-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
            </button>
          {/each}
        {/if}
      </div>
    </section>

    <!-- Floating Action Button -->
    {#if canCreateEvent}
      <button
        class="fixed bottom-6 right-6 w-14 h-14 bg-indigo-600 hover:bg-indigo-700 text-white rounded-2xl shadow-2xl shadow-indigo-300 flex items-center justify-center active:scale-90 active:rotate-12 transition-all duration-300 z-50 border-2 border-white/20 backdrop-blur-sm"
        onclick={() => goto('/events/create')}
        title="New Event"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
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
