<script lang="ts">
  import { onMount } from 'svelte';
  import { getLeagues } from '$lib/services/league';
  import { authService } from '$lib/services/auth.svelte';
  import type { League, Game, User } from '$lib/models';
  import { goto } from '$app/navigation';

  let { data } = $props();

  let allLeagues = $state<League[]>([]);
  let selectedLeagueName = $state<string>('All Leagues');
  let isLoading = $state(false); 
  let error: string | null = $state(null);

  let user = $derived(authService.user);

  $effect(() => {
    if (data.leagues) {
      allLeagues = data.leagues;
    }
  });

  let filteredGames = $derived.by(() => {
    let games: Game[] = [];
    
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
    
    return games.sort((a, b) => {
      const aTime = a.startTime instanceof Date ? a.startTime.getTime() : new Date(a.startTime).getTime();
      const bTime = b.startTime instanceof Date ? b.startTime.getTime() : new Date(b.startTime).getTime();
      return aTime - bTime;
    });
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

{#if isLoading}
  <div class="flex flex-col justify-center items-center h-screen space-y-4">
    <div class="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-indigo-500"></div>
    <p class="text-gray-500 animate-pulse">Loading your leagues...</p>
  </div>
{:else if error}
  <div class="flex flex-col justify-center items-center h-screen space-y-4">
    <div class="bg-red-100 p-6 rounded-lg text-center shadow-md">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-red-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-lg font-bold text-red-800">Oops! Something went wrong</h3>
      <p class="text-red-600 mt-2">{error}</p>
      <button 
        onclick={() => window.location.reload()} 
        class="mt-6 bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-6 rounded-lg transition-colors shadow-sm"
      >
        Try Again
      </button>
    </div>
  </div>
{:else}
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
          <p class="text-4xl font-bold tracking-tight">$30.00</p>
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
            Unpaid Charges: 2 Items
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
        <button class="text-indigo-600 hover:text-indigo-800 text-sm font-bold transition-colors" onclick={() => console.log('See all events')}>See all</button>
      </div>
      
      <div class="space-y-3">
        {#if filteredGames.length === 0}
          <div class="bg-gray-50 border-2 border-dashed border-gray-200 rounded-2xl py-12 text-center">
            <p class="text-gray-400 font-medium italic">No upcoming events scheduled</p>
          </div>
        {:else}
          {#each filteredGames as game (game.id)}
            {@const league = allLeagues.find(l => l.id === game.leagueId)}
            <button
              class="w-full text-left bg-white p-4 rounded-2xl shadow-sm border border-gray-100 cursor-pointer flex items-center space-x-4 hover:shadow-md hover:border-indigo-100 transition-all duration-200 group"
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
          {/each}
        {/if}
      </div>
    </section>

    <!-- Recent Activity -->
    <section class="space-y-4">
      <div class="flex justify-between items-center">
        <h3 class="text-xl font-bold text-gray-900">Recent Activity</h3>
        <button class="text-indigo-600 hover:text-indigo-800 text-sm font-bold transition-colors" onclick={() => console.log('See all activity')}>See all</button>
      </div>
      
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 divide-y divide-gray-50">
        <div class="p-4 flex items-center space-x-4">
          <div class="bg-green-50 p-2.5 rounded-full text-green-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8V9m0 3v-1m0 4v-1m0-8a9 9 0 110 18 9 9 0 010-18z" />
            </svg>
          </div>
          <div class="flex-grow">
            <p class="font-bold text-gray-800 leading-tight">Registration Fee</p>
            <p class="text-xs text-gray-400 font-medium mt-0.5">2 hours ago</p>
          </div>
          <p class="font-black text-green-600">+$50.00</p>
        </div>
        
        <div class="p-4 flex items-center space-x-4">
          <div class="bg-red-50 p-2.5 rounded-full text-red-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div class="flex-grow">
            <p class="font-bold text-gray-800 leading-tight">Equipment Purchase</p>
            <p class="text-xs text-gray-400 font-medium mt-0.5">Yesterday</p>
          </div>
          <p class="font-black text-red-600">-$120.00</p>
        </div>
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
{/if}

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
