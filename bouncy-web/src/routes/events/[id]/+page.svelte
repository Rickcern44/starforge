<script lang="ts">
  import { page } from '$app/stores';
  import { getGameById, updateAttendance, removeAttendance } from '$lib/services/game';
  import { getLeagues } from '$lib/services/league';
  import type { Game, League } from '$lib/models';
  import { authService } from '$lib/services/auth.svelte';
  import { goto } from '$app/navigation';
  import { 
    ChevronLeft, 
    Calendar, 
    Clock, 
    MapPin, 
    DollarSign, 
    Save, 
    Trash2,
    CheckCircle2,
    HelpCircle
  } from 'lucide-svelte';

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
    const now = new Date();
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
    const gTime = game.startTime instanceof Date ? game.startTime.getTime() : new Date(game.startTime).getTime();
    return gTime < startOfToday;
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
      class="btn btn-ghost btn-sm gap-2 opacity-40 hover:opacity-100"
    >
      <ChevronLeft size={16} />
      Back
    </button>
    {#if game}
      <div class="flex flex-wrap items-center gap-2">
        {#if isAdmin}
          <button 
            onclick={() => goto(`/leagues/${game!.leagueId}/admin`)}
            class="btn btn-neutral btn-xs rounded-full font-black uppercase tracking-widest px-3"
          >
            Manage League
          </button>
        {/if}
        {#if game.isCanceled}
          <span class="badge badge-error badge-sm font-black uppercase tracking-widest py-3 border-none">Canceled</span>
        {:else}
          <span class="badge badge-success badge-sm font-black uppercase tracking-widest py-3 border-none text-white">Active</span>
        {/if}
      </div>
    {/if}
  </header>

  {#if game}
    <div>
      <h1 class="text-3xl font-black tracking-tight leading-tight">Game at {game.location}</h1>
      <p class="text-sm font-bold opacity-40 mt-1 uppercase tracking-wider">{leagueName}</p>
    </div>
  {/if}

  {#if !game}
    <div class="card bg-base-100 p-12 shadow-sm border border-base-200 text-center space-y-4 rounded-[32px]">
      <div class="bg-base-200 w-20 h-20 rounded-full flex items-center justify-center mx-auto text-base-content/20">
        <HelpCircle size={40} />
      </div>
      <h3 class="text-xl font-bold">Event not found</h3>
      <button onclick={() => goto('/')} class="btn btn-neutral px-8">Return Home</button>
    </div>
  {:else}
    <!-- Compact Info Grid -->
    <div class="grid grid-cols-2 gap-3">
      <div class="card bg-base-100 p-3 border border-base-200 shadow-sm flex-row items-center gap-3 rounded-2xl">
        <div class="bg-base-200 p-2 rounded-xl text-base-content opacity-40">
          <Calendar size={20} />
        </div>
        <div>
          <p class="text-[9px] uppercase font-black tracking-widest opacity-40 leading-none mb-0.5">Date</p>
          <p class="text-xs font-bold">{formatDateTime(game.startTime).shortDate}</p>
        </div>
      </div>
      
      <div class="card bg-base-100 p-3 border border-base-200 shadow-sm flex-row items-center gap-3 rounded-2xl">
        <div class="bg-base-200 p-2 rounded-xl text-base-content opacity-40">
          <Clock size={20} />
        </div>
        <div>
          <p class="text-[9px] uppercase font-black tracking-widest opacity-40 leading-none mb-0.5">Time</p>
          <p class="text-xs font-bold">{formatDateTime(game.startTime).time}</p>
        </div>
      </div>

      <div class="card bg-base-100 p-3 border border-base-200 shadow-sm flex-row items-center gap-3 rounded-2xl">
        <div class="bg-base-200 p-2 rounded-xl text-base-content opacity-40">
          <MapPin size={20} />
        </div>
        <div class="overflow-hidden">
          <p class="text-[9px] uppercase font-black tracking-widest opacity-40 leading-none mb-0.5">Location</p>
          <p class="text-xs font-bold truncate">{game.location}</p>
        </div>
      </div>

      <div class="card bg-base-100 p-3 border border-base-200 shadow-sm flex-row items-center gap-3 rounded-2xl">
        <div class="bg-base-200 p-2 rounded-xl text-base-content opacity-40">
          <DollarSign size={20} />
        </div>
        <div>
          <p class="text-[9px] uppercase font-black tracking-widest opacity-40 leading-none mb-0.5">Cost</p>
          <p class="text-xs font-bold">${(game.costInCents / 100).toFixed(2)}</p>
        </div>
      </div>
    </div>

    <!-- Attendance RSVP (More Compact) -->
    <section class="card bg-base-100 p-6 rounded-[32px] border border-base-200 shadow-sm space-y-4 relative overflow-hidden">
      {#if !canEditRSVP}
        <div class="absolute inset-0 bg-base-100/60 backdrop-blur-[1px] z-10 flex items-center justify-center p-6 text-center">
          <div class="card bg-base-100 p-4 shadow-xl border border-base-200 max-w-[200px] rounded-2xl">
            <p class="text-xs font-black opacity-40 uppercase tracking-widest leading-tight">Registration closed for this past event</p>
          </div>
        </div>
      {/if}
      
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-black">Your RSVP</h3>
        {#if userStatus}
          <span class="text-[10px] font-black uppercase tracking-widest opacity-40">Current: {userStatus}</span>
        {/if}
      </div>
      
      <div class="grid grid-cols-3 gap-2">
        {#each ['Yes', 'No', 'Tentative'] as status}
          <button 
            class="btn btn-outline btn-md font-black uppercase tracking-widest text-[10px]"
            class:btn-active={selectedStatus === status}
            class:btn-neutral={selectedStatus === status}
            onclick={() => selectedStatus = status}
            disabled={!canEditRSVP}
          >
            {status}
          </button>
        {/each}
      </div>

      <div class="space-y-1.5">
        <textarea
          bind:value={comment}
          placeholder="Optional comment..."
          class="textarea textarea-bordered w-full focus:outline-none placeholder:opacity-30 text-sm font-medium"
          rows="1"
          disabled={!canEditRSVP}
        ></textarea>
      </div>

      <button
        class="btn btn-neutral btn-block shadow-lg shadow-neutral/20 font-black uppercase tracking-widest h-auto py-4"
        disabled={isUpdating || !selectedStatus || !canEditRSVP}
        onclick={handleUpdateAttendance}
      >
        {#if isUpdating}
          <span class="loading loading-spinner loading-sm"></span>
          Saving...
        {:else}
          <Save size={18} />
          Save RSVP
          {#if !canEditRSVP} (Admin Only){/if}
        {/if}
      </button>
    </section>

    <!-- Attendees List -->
    <section class="space-y-3">
      <div class="flex justify-between items-end px-1">
        <h3 class="text-xl font-black">Attendees</h3>
        <span class="badge badge-neutral badge-sm font-black uppercase tracking-widest py-3 px-3">{game.attendance.length} playing</span>
      </div>
      
      {#if game.attendance.length === 0}
        <div class="bg-base-200/50 border-2 border-dashed border-base-300 rounded-3xl py-10 text-center">
          <p class="opacity-40 font-bold italic text-sm">No RSVPs yet</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 gap-2">
          {#each game.attendance as attendance}
            <div class="card bg-base-100 p-3 border border-base-200 shadow-sm flex-row items-center justify-between group rounded-2xl text-left">
              <div class="flex items-center space-x-3">
                <div class="avatar placeholder opacity-40">
                  <div class="bg-base-200 text-base-content rounded-full w-8 h-8 font-black text-[10px] uppercase border border-base-300">
                    {(attendance.userName || attendance.userId).substring(0, 1)}
                  </div>
                </div>
                <div>
                  <p class="font-bold text-sm leading-none mb-1">
                    {attendance.userName || `User ${attendance.userId.substring(0, 5)}`}
                  </p>
                  {#if attendance.checkInComment}
                    <p class="text-[10px] opacity-40 italic truncate max-w-[150px]">"{attendance.checkInComment}"</p>
                  {/if}
                </div>
              </div>
              
              <div class="flex items-center space-x-2">
                {#if isAdmin}
                  <button 
                    onclick={() => handleRemoveAttendance(attendance.userId, attendance.userName || 'this user')}
                    class="btn btn-ghost btn-xs btn-circle text-base-content/20 hover:text-error transition-colors md:opacity-0 group-hover:opacity-100"
                    title="Remove attendee"
                  >
                    <Trash2 size={14} />
                  </button>
                {/if}
                
                {#if attendance.status === 0}
                  <span class="badge badge-success badge-sm font-black uppercase tracking-widest py-3 text-white border-none text-[9px]">Yes</span>
                {:else if attendance.status === 1}
                  <span class="badge badge-error badge-sm font-black uppercase tracking-widest py-3 border-none text-[9px]">No</span>
                {:else}
                  <span class="badge badge-warning badge-sm font-black uppercase tracking-widest py-3 border-none text-[9px]">Maybe</span>
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
    background: var(--color-base-300);
    border-radius: 10px;
  }
  ::-webkit-scrollbar-thumb:hover {
    background: var(--color-base-content);
    opacity: 0.1;
  }
</style>
