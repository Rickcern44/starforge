<script lang="ts">
  import type { League } from '$lib/models';
  import { Trophy, Plus, Settings, Users, ArrowRight } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { authService } from '$lib/services/auth.svelte';
  import { featureFlagService } from '$lib/services/feature-flag.svelte';

  let { data } = $props();
  let leagues = $derived(data.leagues || []);
  let user = $derived(authService.user);
</script>

<div class="space-y-8">
  <header class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-black tracking-tight leading-tight">Leagues</h2>
      <p class="text-sm font-bold opacity-40 uppercase tracking-wider">Management & Configuration</p>
    </div>
    {#if (user?.roles?.includes('league_creator') || user?.roles?.includes('admin')) && featureFlagService.isEnabled('league_creation')}
      <button onclick={() => goto('/leagues/create')} class="btn btn-primary btn-sm gap-2 rounded-xl px-4">
        <Plus size={16} />
        Create League
      </button>
    {/if}
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#each leagues as league}
      <div class="card bg-base-100 border border-base-300 hover:border-primary/50 transition-all group overflow-hidden">
        <div class="p-6 space-y-4">
          <div class="flex justify-between items-start">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 bg-primary/10 text-primary rounded-2xl flex items-center justify-center">
                <Trophy size={24} />
              </div>
              <div>
                <h3 class="font-black text-xl group-hover:text-primary transition-colors">{league.name}</h3>
                <p class="text-xs font-bold opacity-40 uppercase tracking-widest">{league.id}</p>
              </div>
            </div>
          </div>

          <div class="flex gap-2">
            <button 
              onclick={() => goto(`/leagues/${league.id}/admin`)}
              class="btn btn-ghost btn-sm flex-1 gap-2 rounded-xl"
            >
              <Settings size={14} />
              Admin
            </button>
            <button 
              class="btn btn-ghost btn-sm flex-1 gap-2 rounded-xl"
              onclick={() => goto(`/leagues/${league.id}/admin?tab=members`)}
            >
              <Users size={14} />
              Members
            </button>
          </div>
        </div>
        
        <button 
          onclick={() => goto(`/leagues/${league.id}/admin`)}
          class="bg-base-200 py-2 px-6 flex justify-between items-center text-xs font-black uppercase tracking-widest opacity-0 group-hover:opacity-100 transition-opacity"
        >
          Manage Dashboard
          <ArrowRight size={14} />
        </button>
      </div>
    {/each}

    {#if leagues.length === 0}
      <div class="col-span-full p-12 card bg-base-100 border border-dashed border-base-300 items-center justify-center text-center opacity-40">
        <Trophy size={48} class="mb-4" />
        <p class="font-black uppercase tracking-widest">No leagues found</p>
        <p class="text-sm">Create your first league to get started.</p>
      </div>
    {/if}
  </div>
</div>
