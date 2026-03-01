<script lang="ts">
  import { createGame } from '$lib/services/game';
  import type { League } from '$lib/models';
  import { authService } from '$lib/services/auth.svelte';
  import { goto } from '$app/navigation';

  let { data } = $props();

  let leagues = $state<League[]>(data.leagues || []);
  let editableLeagues = $derived.by(() => {
    if (!authService.user) return [];
    return leagues.filter(l => {
      const member = l.members.find(m => m.playerId === authService.user!.id);
      if (!member) return false;
      const role = member.role.toLowerCase();
      return role.includes('admin') || role.includes('owner');
    });
  });

  let title = $state('');
  let selectedLeagueId = $state('');
  let date = $state('');
  let time = $state('');
  let price = $state(0);
  let description = $state('');
  let isRecurring = $state(false);
  let recurrenceInterval = $state('Weekly');
  let isLoading = $state(false);

  const recurrenceOptions = ['Daily', 'Weekly', 'Bi-weekly', 'Monthly'];

  $effect(() => {
    if (data.leagues) leagues = data.leagues;
  });

  $effect(() => {
    if (editableLeagues.length === 1 && !selectedLeagueId) {
      selectedLeagueId = editableLeagues[0].id;
    }
  });

  async function handleCreateEvent(e: SubmitEvent) {
    e.preventDefault();
    if (!selectedLeagueId || !date || !time) {
      alert('Please fill in all required fields.');
      return;
    }

    isLoading = true;
    const startTime = new Date(`${date}T${time}`).toISOString();
    
    const success = await createGame({
      leagueId: selectedLeagueId,
      location: title, // Using title as location for now to match API expectations if location is required
      startTime,
      costInCents: Math.round(price * 100)
    });

    if (success) {
      alert('Event created successfully!');
      goto('/');
    } else {
      alert('Failed to create event.');
    }
    isLoading = false;
  }
</script>

<svelte:head>
  <title>Create Event | League Manager</title>
</svelte:head>

<div class="max-w-2xl mx-auto py-8">
  <div class="flex items-center space-x-4 mb-8">
    <button onclick={() => history.back()} class="text-gray-500 hover:text-indigo-600 transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
    <h1 class="text-3xl font-bold text-gray-800">Create New Event</h1>
  </div>

  <form onsubmit={handleCreateEvent} class="bg-white p-8 rounded-xl shadow-md space-y-6">
    <div class="space-y-2">
      <label for="title" class="block text-sm font-semibold text-gray-700">Event Title / Location</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
        </div>
        <input
          id="title"
          type="text"
          bind:value={title}
          required
          placeholder="e.g. Game at Memorial Park"
          class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
    </div>

    <div class="space-y-2">
      <label for="league" class="block text-sm font-semibold text-gray-700">League</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        </div>
        <select
          id="league"
          bind:value={selectedLeagueId}
          required
          disabled={editableLeagues.length <= 1}
          class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 appearance-none bg-white disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed"
        >
          <option value="" disabled>Select a league</option>
          {#each editableLeagues as league}
            <option value={league.id}>{league.name}</option>
          {/each}
        </select>
        {#if editableLeagues.length === 0}
          <p class="mt-1 text-xs text-red-500 italic">You don't have permission to create events in any leagues.</p>
        {/if}
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="space-y-2">
        <label for="date" class="block text-sm font-semibold text-gray-700">Date</label>
        <input
          id="date"
          type="date"
          bind:value={date}
          required
          class="block w-full px-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
      <div class="space-y-2">
        <label for="time" class="block text-sm font-semibold text-gray-700">Time</label>
        <input
          id="time"
          type="time"
          bind:value={time}
          required
          class="block w-full px-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
    </div>

    <div class="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
      <input
        id="recurring"
        type="checkbox"
        bind:checked={isRecurring}
        class="h-5 w-5 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
      />
      <label for="recurring" class="flex-grow">
        <span class="block text-sm font-semibold text-gray-700">Recurring Event</span>
        <span class="block text-xs text-gray-500">Repeat this event automatically</span>
      </label>
    </div>

    {#if isRecurring}
      <div class="space-y-2 animate-in fade-in slide-in-from-top-2 duration-200">
        <label for="interval" class="block text-sm font-semibold text-gray-700">Repeat Interval</label>
        <select
          id="interval"
          bind:value={recurrenceInterval}
          class="block w-full px-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 bg-white"
        >
          {#each recurrenceOptions as opt}
            <option value={opt}>{opt}</option>
          {/each}
        </select>
      </div>
    {/if}

    <div class="space-y-2">
      <label for="price" class="block text-sm font-semibold text-gray-700">Price ($)</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8V9m0 3v-1m0 4v-1m0-8a9 9 0 110 18 9 9 0 010-18z" />
          </svg>
        </div>
        <input
          id="price"
          type="number"
          step="0.01"
          min="0"
          bind:value={price}
          class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
          placeholder="0.00"
        />
      </div>
    </div>

    <div class="space-y-2">
      <label for="description" class="block text-sm font-semibold text-gray-700">Description (Optional)</label>
      <textarea
        id="description"
        bind:value={description}
        rows="3"
        placeholder="Add details about the event..."
        class="block w-full px-3 py-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
      ></textarea>
    </div>

    <button
      type="submit"
      disabled={isLoading || editableLeagues.length === 0}
      class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-4 px-6 rounded-xl shadow-lg transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {#if isLoading}
        <div class="flex items-center justify-center">
          <div class="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-white mr-2"></div>
          Creating Event...
        </div>
      {:else}
        Create Event
      {/if}
    </button>
  </form>
</div>
