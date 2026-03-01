<script lang="ts">
  import { page } from '$app/stores';
  import { getGameById, updateAttendance, removeAttendance } from '$lib/services/game';
  import { getLeagues } from '$lib/services/league';
  import type { Game, League } from '$lib/models';
  import { authService } from '$lib/services/auth.svelte';
  import { goto } from '$app/navigation';

  let { data } = $props();

  let gameId = $derived($page.params.id);
  let game = $state<Game | null>(data.game);
  let leagues = $state<League[]>(data.leagues || []);
  let isUpdating = $state(false);

  let user = $derived(authService.user);

  let selectedStatus = $state<string | null>(null);
  let comment = $state('');

  let leagueName = $derived.by(() => {
    if (!game || leagues.length === 0) return 'Loading...';
    const league = leagues.find(l => l.id === game!.leagueId);
    return league ? league.name : 'Unknown League';
  });

  let currentLeague = $derived(leagues.find(l => l.id === game?.leagueId));
  let isAdmin = $derived.by(() => {
    if (!user || !currentLeague) return false;
    const member = currentLeague.members.find(m => m.playerId === user.id);
    return member && (member.role.toLowerCase().includes('admin') || member.role.toLowerCase().includes('owner'));
  });

  let isPastGame = $derived.by(() => {
    if (!game) return false;
    const gTime = game.startTime instanceof Date ? game.startTime.getTime() : new Date(game.startTime).getTime();
    return gTime < new Date().getTime();
  });

  let canEditRSVP = $derived(!isPastGame || isAdmin);

  // Keep this for when navigation happens client-side without full reload
  $effect(() => {
    if (data.game) game = data.game;
    if (data.leagues) leagues = data.leagues;
  });

  async function handleRemoveAttendance(targetUserId: string, targetUserName: string) {
    if (!game) return;
    if (!confirm(`Are you sure you want to remove ${targetUserName || 'this user'} from the attendance list?`)) return;

    const previousAttendance = [...game.attendance];
    // Optimistic Update
    game.attendance = game.attendance.filter(a => a.userId !== targetUserId);

    try {
      const success = await removeAttendance(gameId, targetUserId);
      if (!success) {
        game.attendance = previousAttendance;
        alert('Failed to remove attendance.');
      }
    } catch (err) {
      console.error('[Event Page] Error in handleRemoveAttendance:', err);
      game.attendance = previousAttendance;
      alert('An error occurred while removing attendance.');
    }
  }

  async function handleUpdateAttendance() {
    if (!selectedStatus || !user || !game) {
      alert('Please select a status first!');
      return;
    }

    const previousAttendance = [...game.attendance];
    isUpdating = true;
    
    // Status Map: Yes=0, No=1, Tentative=2 (matches Go iota)
    const statusMap: Record<string, number> = { 'Yes': 0, 'No': 1, 'Tentative': 2 };
    const status = statusMap[selectedStatus];

    // Optimistic Update
    const newRecord = {
      userId: user.id,
      checkedIn: false,
      status: status,
      checkInComment: comment,
      createdAt: new Date(),
      updatedAt: new Date()
    };

    const existingIndex = game.attendance.findIndex(a => a.userId === user.id);
    if (existingIndex > -1) {
      game.attendance[existingIndex] = newRecord;
    } else {
      game.attendance.push(newRecord);
    }

    try {
      const success = await updateAttendance(gameId, status, comment);

      if (success) {
        // Fetch fresh game data to sync with server state (names, etc)
        const updatedGame = await getGameById(gameId);
        if (updatedGame) {
          game = updatedGame;
        }
      } else {
        // Rollback on failure
        game.attendance = previousAttendance;
        alert('Failed to update attendance.');
      }
    } catch (err) {
      console.error('[Event Page] Error in handleUpdateAttendance:', err);
      // Rollback on error
      game.attendance = previousAttendance;
      alert('An error occurred while updating attendance.');
    } finally {
      isUpdating = false;
    }
  }

  function formatDateTime(dateStr: string) {
    const d = new Date(dateStr);
    return {
      date: d.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric', timeZone: 'America/New_York' }),
      time: d.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/New_York' }),
      shortDate: d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', timeZone: 'America/New_York' })
    };
  }

  let userStatus = $derived.by(() => {
    if (!user || !game || !game.attendance) return null;
    const att = game.attendance.find(a => a.userId === user.id);
    if (!att) return null;
    const statusMap: Record<number, string> = { 0: 'Yes', 1: 'No', 2: 'Tentative' };
    return statusMap[att.status];
  });

  // Sync selected status with current attendance if available
  $effect(() => {
    if (userStatus && !selectedStatus) {
      selectedStatus = userStatus;
    }
  });
</script>

<svelte:head>
  <title>{game ? `Game at ${game.location}` : 'Event Details'} | Bouncy</title>
</svelte:head>

<div class="max-w-2xl mx-auto py-4 px-4 space-y-6">
  <!-- Top Navigation & Title -->
  <header class="flex justify-between items-center">
    <button 
      onclick={() => history.back()} 
      class="flex items-center text-sm font-bold text-gray-500 hover:text-indigo-600 transition-colors group"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1 transform group-hover:-translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7" />
      </svg>
      Back
    </button>
    {#if game}
      <div class="flex items-center space-x-2">
        {#if isAdmin}
          <button 
            onclick={() => goto(`/leagues/${game!.leagueId}/admin`)}
            class="px-3 py-1 bg-indigo-50 text-indigo-600 rounded-full text-[10px] font-black uppercase tracking-widest hover:bg-indigo-100 transition-colors"
          >
            Manage League
          </button>
        {/if}
        {#if game.isCanceled}
          <span class="px-3 py-1 bg-red-100 text-red-600 rounded-full text-[10px] font-black uppercase tracking-widest">Canceled</span>
        {:else}
          <span class="px-3 py-1 bg-green-100 text-green-600 rounded-full text-[10px] font-black uppercase tracking-widest">Active</span>
        {/if}
      </div>
    {/if}
  </header>

  {#if game}
    <div>
      <h1 class="text-3xl font-black text-gray-900 tracking-tight leading-tight">Game at {game.location}</h1>
      <p class="text-sm font-bold text-indigo-600 mt-1 uppercase tracking-wider">{leagueName}</p>
    </div>
  {/if}

  {#if !game}
    <div class="bg-white p-12 rounded-3xl shadow-sm border border-gray-100 text-center space-y-4">
      <div class="bg-gray-50 w-20 h-20 rounded-full flex items-center justify-center mx-auto text-gray-300">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 9.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-gray-800">Event not found</h3>
      <button onclick={() => goto('/')} class="bg-indigo-600 text-white font-bold px-8 py-3 rounded-2xl hover:bg-indigo-700 transition-all">Return Home</button>
    </div>
  {:else}
    <!-- Compact Info Grid -->
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-white p-3 rounded-2xl border border-gray-100 shadow-sm flex items-center space-x-3">
        <div class="bg-indigo-50 p-2 rounded-xl text-indigo-600">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <div>
          <p class="text-[9px] text-gray-400 uppercase font-black tracking-widest leading-none mb-0.5">Date</p>
          <p class="text-xs font-bold text-gray-800">{formatDateTime(game.startTime).shortDate}</p>
        </div>
      </div>
      
      <div class="bg-white p-3 rounded-2xl border border-gray-100 shadow-sm flex items-center space-x-3">
        <div class="bg-indigo-50 p-2 rounded-xl text-indigo-600">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div>
          <p class="text-[9px] text-gray-400 uppercase font-black tracking-widest leading-none mb-0.5">Time</p>
          <p class="text-xs font-bold text-gray-800">{formatDateTime(game.startTime).time}</p>
        </div>
      </div>

      <div class="bg-white p-3 rounded-2xl border border-gray-100 shadow-sm flex items-center space-x-3">
        <div class="bg-indigo-50 p-2 rounded-xl text-indigo-600">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
        </div>
        <div class="overflow-hidden">
          <p class="text-[9px] text-gray-400 uppercase font-black tracking-widest leading-none mb-0.5">Location</p>
          <p class="text-xs font-bold text-gray-800 truncate">{game.location}</p>
        </div>
      </div>

      <div class="bg-white p-3 rounded-2xl border border-gray-100 shadow-sm flex items-center space-x-3">
        <div class="bg-indigo-50 p-2 rounded-xl text-indigo-600">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8V9m0 3v-1m0 4v-1m0-8a9 9 0 110 18 9 9 0 010-18z" />
          </svg>
        </div>
        <div>
          <p class="text-[9px] text-gray-400 uppercase font-black tracking-widest leading-none mb-0.5">Cost</p>
          <p class="text-xs font-bold text-gray-800">${(game.costInCents / 100).toFixed(2)}</p>
        </div>
      </div>
    </div>

    <!-- Attendance RSVP (More Compact) -->
    <section class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-4 relative overflow-hidden">
      {#if !canEditRSVP}
        <div class="absolute inset-0 bg-white/60 backdrop-blur-[1px] z-10 flex items-center justify-center p-6 text-center">
          <div class="bg-white p-4 rounded-2xl shadow-xl border border-gray-100 max-w-[200px]">
            <p class="text-xs font-black text-gray-400 uppercase tracking-widest leading-tight">Registration closed for this past event</p>
          </div>
        </div>
      {/if}
      
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-black text-gray-900">Your RSVP</h3>
        {#if userStatus}
          <span class="text-[10px] font-black uppercase tracking-widest text-indigo-400">Status: {userStatus}</span>
        {/if}
      </div>
      
      <div class="grid grid-cols-3 gap-2">
        {#each ['Yes', 'No', 'Tentative'] as status}
          <button 
            class="flex flex-col items-center justify-center py-3 rounded-xl border-2 transition-all duration-200"
            class:bg-indigo-600={selectedStatus === status}
            class:border-indigo-600={selectedStatus === status}
            class:text-white={selectedStatus === status}
            class:bg-white={selectedStatus !== status}
            class:border-gray-50={selectedStatus !== status}
            class:text-gray-400={selectedStatus !== status}
            onclick={() => selectedStatus = status}
            disabled={!canEditRSVP}
          >
            <span class="text-[10px] font-black uppercase tracking-widest">{status}</span>
          </button>
        {/each}
      </div>

      <div class="space-y-1.5">
        <textarea
          bind:value={comment}
          placeholder="Optional comment..."
          class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 transition-all placeholder:text-gray-300 text-sm font-medium"
          rows="1"
          disabled={!canEditRSVP}
        ></textarea>
      </div>

      <button
        class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-black py-3 px-6 rounded-xl shadow-lg shadow-indigo-100 transition-all duration-200 disabled:opacity-50"
        disabled={isUpdating || !selectedStatus || !canEditRSVP}
        onclick={handleUpdateAttendance}
      >
        {#if isUpdating}
          <div class="flex items-center justify-center text-sm">
            <div class="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white mr-2"></div>
            Saving...
          </div>
        {:else}
          Save RSVP
          {#if !canEditRSVP} (Admin Only){/if}
        {/if}
      </button>
    </section>

    <!-- Attendees List -->
    <section class="space-y-3">
      <div class="flex justify-between items-end px-1">
        <h3 class="text-xl font-black text-gray-900">Attendees</h3>
        <span class="text-indigo-600 font-black text-sm">{game.attendance.length} playing</span>
      </div>
      
      {#if game.attendance.length === 0}
        <div class="bg-gray-50 border-2 border-dashed border-gray-200 rounded-3xl py-10 text-center">
          <p class="text-gray-400 font-bold italic text-sm">No RSVPs yet</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 gap-2">
          {#each game.attendance as attendance}
            <div class="bg-white p-3 rounded-2xl border border-gray-100 shadow-sm flex items-center justify-between group">
              <div class="flex items-center space-x-3">
                <div class="h-8 w-8 rounded-full bg-gray-50 border border-gray-100 flex items-center justify-center text-gray-400 font-black text-[10px] uppercase">
                  {(attendance.userName || attendance.userId).substring(0, 1)}
                </div>
                <div>
                  <p class="font-bold text-gray-800 text-sm leading-none mb-1">
                    {attendance.userName || `User ${attendance.userId.substring(0, 5)}`}
                  </p>
                  {#if attendance.checkInComment}
                    <p class="text-[10px] text-gray-400 italic truncate max-w-[150px]">"{attendance.checkInComment}"</p>
                  {/if}
                </div>
              </div>
              
              <div class="flex items-center space-x-2">
                {#if isAdmin}
                  <button 
                    onclick={() => handleRemoveAttendance(attendance.userId, attendance.userName || 'this user')}
                    class="p-1.5 text-gray-300 hover:text-red-500 transition-colors opacity-0 group-hover:opacity-100"
                    title="Remove attendee"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                {/if}
                
                {#if attendance.status === 0}
                  <span class="px-2 py-0.5 bg-green-50 text-green-600 rounded-md text-[9px] font-black uppercase tracking-widest">Yes</span>
                {:else if attendance.status === 1}
                  <span class="px-2 py-0.5 bg-red-50 text-red-600 rounded-md text-[9px] font-black uppercase tracking-widest">No</span>
                {:else}
                  <span class="px-2 py-0.5 bg-orange-50 text-orange-600 rounded-md text-[9px] font-black uppercase tracking-widest">Maybe</span>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </section>
  {/if}
</div>

<style>
  /* Extra minimal scrollbar */
  ::-webkit-scrollbar {
    width: 6px;
  }
  ::-webkit-scrollbar-track {
    background: transparent;
  }
  ::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 10px;
  }
  ::-webkit-scrollbar-thumb:hover {
    background: #cbd5e1;
  }
</style>
