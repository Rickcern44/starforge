<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import { featureFlagService } from '$lib/services/feature-flag.svelte';
  import type { League, Game } from '$lib/models';
  import { goto } from '$app/navigation';
  import { 
    Bell, 
    Wallet, 
    ChevronRight, 
    Settings, 
    Calendar, 
    Plus, 
    SlidersHorizontal
  } from 'lucide-svelte';

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
    const thirtyDaysFromNow = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000).getTime();
    
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
        return gTime >= startOfToday && gTime <= thirtyDaysFromNow;
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
    if (user.roles?.includes('admin')) return true;
    
    for (const league of (data.leagues || [])) {
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
  <title>Dashboard | Bouncy</title>
</svelte:head>

<div class="max-w-xl mx-auto space-y-8 py-2">
    <header class="flex justify-between items-end px-1">
      <div>
        <p class="text-[10px] font-black opacity-40 uppercase tracking-widest mb-1">Welcome back,</p>
        <h2 class="text-3xl font-black tracking-tight">{user ? user.name : 'Guest'}</h2>
      </div>
      <button class="btn btn-ghost btn-circle border border-base-300">
        <Bell size={20} />
      </button>
    </header>

    {#if user && (user.roles?.includes('admin') || managedLeagues.length > 0)}
      <section class="space-y-4 px-1">
        <div class="flex justify-between items-center">
          <h3 class="text-xs font-black uppercase tracking-widest opacity-40">Quick Management</h3>
          <button class="btn btn-link btn-xs no-underline hover:no-underline font-black uppercase tracking-widest text-[9px] opacity-40 hover:opacity-100" onclick={() => goto('/manage')}>See All</button>
        </div>
        <div class="flex space-x-3 overflow-x-auto pb-2 no-scrollbar -mx-1 px-1">
          {#if user.roles?.includes('admin')}
            <button
              class="flex-shrink-0 w-48 p-4 bg-neutral text-neutral-content rounded-[22px] hover:opacity-90 shadow-lg active:scale-[0.98] transition-all group border-none text-left"
              onclick={() => goto('/manage')}
            >
              <div class="bg-white/10 w-10 h-10 rounded-xl flex items-center justify-center mb-3">
                <SlidersHorizontal size={20} />
              </div>
              <p class="font-bold text-sm leading-tight">System Portal</p>
              <p class="text-[9px] font-black opacity-40 uppercase tracking-widest mt-1">Platform Control</p>
            </button>
          {/if}
          
          {#each managedLeagues as league}
            <button
              class="flex-shrink-0 w-48 p-4 bg-base-100 border border-base-300 rounded-[22px] hover:border-neutral shadow-sm active:scale-[0.98] transition-all group text-left"
              onclick={() => goto(`/leagues/${league.id}/admin`)}
            >
              <div class="bg-base-200 w-10 h-10 rounded-xl flex items-center justify-center mb-3 group-hover:bg-neutral group-hover:text-neutral-content transition-colors">
                <Settings size={20} />
              </div>
              <p class="font-bold text-sm leading-tight truncate">{league.name}</p>
              <p class="text-[9px] font-black opacity-40 uppercase tracking-widest mt-1">League Admin</p>
            </button>
          {/each}

          {#if featureFlagService.isEnabled('league_creation')}
          <button
            class="flex-shrink-0 w-48 p-4 bg-base-100 border-2 border-dashed border-base-300 rounded-[22px] hover:border-neutral active:scale-[0.98] transition-all group flex flex-col items-center justify-center text-center space-y-2"
            onclick={() => goto('/leagues/create')}
          >
            <div class="w-10 h-10 rounded-full bg-base-200 flex items-center justify-center group-hover:bg-neutral group-hover:text-neutral-content transition-colors">
              <Plus size={20} />
            </div>
            <p class="font-black text-[10px] uppercase tracking-widest opacity-40">New League</p>
          </button>
          {/if}
        </div>
      </section>
    {/if}

    <!-- Financial Summary (DaisyUI Card) -->
    {#if featureFlagService.isEnabled('payments')}
    <button 
      class="card w-full bg-neutral text-neutral-content shadow-2xl shadow-neutral/20 cursor-pointer hover:opacity-90 transition-all duration-300 relative overflow-hidden group border-none text-left" 
      onclick={() => goto('/ledger')}
    >
      <div class="card-body p-6">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-bold uppercase tracking-widest opacity-60 mb-1">Outstanding Balance</p>
            <p class="text-5xl font-black tracking-tighter">${financialSummary.balance.toFixed(2)}</p>
          </div>
          <div class="bg-white/10 p-3 rounded-2xl backdrop-blur-md border border-white/10 text-white">
            <Wallet size={24} />
          </div>
        </div>
        
        <div class="card-actions justify-between items-center text-neutral-content/90 mt-2">
          <div class="badge badge-warning border-none gap-2 px-3 py-3 font-black text-[11px] uppercase tracking-wider text-warning-content">
            <span class="w-2 h-2 rounded-full bg-warning-content animate-pulse"></span>
            <span>{financialSummary.count} {financialSummary.count === 1 ? 'Charge' : 'Charges'} pending</span>
          </div>
          <div class="flex items-center text-xs font-bold uppercase tracking-widest opacity-60 group-hover:opacity-100 transition-opacity">
            Details
            <ChevronRight size={16} class="ml-1 group-hover:translate-x-1 transition-transform" />
          </div>
        </div>
      </div>
    </button>
    {/if}

    <!-- League Selector -->
    <section class="space-y-4 px-1">
      <div class="flex space-x-2 overflow-x-auto pb-2 no-scrollbar">
        {#each ['All Leagues', ...allLeagues.map(l => l.name)] as name}
          <button
            class="btn btn-sm h-10 px-6 rounded-2xl font-black whitespace-nowrap uppercase tracking-wide border transition-all"
            class:btn-neutral={selectedLeagueName === name}
            class:btn-ghost={selectedLeagueName !== name}
            class:bg-base-100={selectedLeagueName !== name}
            class:border-base-300={selectedLeagueName !== name}
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
        <h3 class="text-xl font-black tracking-tight">Schedule</h3>
        <button class="btn btn-link btn-xs opacity-40 hover:opacity-100 font-black uppercase tracking-widest no-underline hover:no-underline" onclick={() => goto('/events')}>View All</button>
      </div>
      
      <div class="space-y-4">
        {#if upcomingGames.length === 0}
          <div class="bg-base-300/30 border-2 border-dashed border-base-300 rounded-[28px] py-16 text-center">
            <p class="opacity-40 font-bold italic text-sm text-base-content/40">No games on the horizon</p>
          </div>
        {:else}
          {#each upcomingGames as game (game.id)}
            {@const league = allLeagues.find(l => l.id === game.leagueId)}
            {@const member = league?.members.find(m => m.playerId === user?.id)}
            {@const isAdmin = member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'))}
            
            <button
              class="card w-full bg-base-100 border border-base-300 shadow-sm hover:border-neutral active:scale-[0.98] transition-all duration-300 group overflow-hidden text-left"
              onclick={() => goto(`/events/${game.id}`)}
            >
              <div class="card-body p-5 flex-row items-center gap-4">
                <div class="bg-base-200 p-3 rounded-2xl text-base-content/40 group-active:bg-neutral group-active:text-neutral-content transition-colors duration-200">
                  <Calendar size={24} />
                </div>
                <div class="flex-grow">
                  <div class="flex items-center space-x-2">
                    <span class="text-[10px] font-black opacity-40 uppercase tracking-widest">{league?.name || 'League'}</span>
                    <span class="w-1 h-1 rounded-full bg-base-300"></span>
                    <span class="text-[10px] font-bold opacity-40 uppercase tracking-widest">{new Date(game.startTime).toLocaleDateString('en-US', { weekday: 'short' })}</span>
                  </div>
                  <p class="font-black text-lg leading-tight mt-0.5 group-hover:text-neutral transition-colors">{game.location}</p>
                </div>
                <div class="text-right flex flex-col items-end">
                  <p class="font-black text-sm">{new Date(game.startTime).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/New_York' })}</p>
                  <div class="mt-1 flex items-center text-[10px] font-black opacity-20 uppercase tracking-widest">
                    Details
                    <ChevronRight size={12} strokeWidth={4} class="ml-0.5" />
                  </div>
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
        class="btn btn-neutral btn-circle btn-lg fixed bottom-6 right-6 shadow-2xl shadow-neutral/30 z-50 border-2 border-white/10 active:scale-90"
        onclick={() => goto('/events/create')}
        title="New Event"
      >
        <Plus size={32} strokeWidth={3} />
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
