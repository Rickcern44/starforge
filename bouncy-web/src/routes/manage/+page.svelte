<script lang="ts">
  import { Trophy, Users, Shield, ArrowRight, Activity, Plus, ChevronRight, Send } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { authService } from '$lib/services/auth.svelte';
  import { toastService } from '$lib/services/toast.svelte';
  import { featureFlagService } from '$lib/services/feature-flag.svelte';
  
  let { data } = $props();
  let leagues = $derived(data.leagues || []);
  let user = $derived(authService.user);
  
  let totalMembers = $derived.by(() => {
    return leagues.reduce((sum: number, l: any) => sum + (l.members?.length || 0), 0);
  });

  // Invitation state
  let inviteEmail = $state('');
  let isInviting = $state(false);

  async function handleInvite() {
    if (!inviteEmail) return;
    isInviting = true;
    const success = await authService.inviteLeagueCreator(inviteEmail);
    if (success) {
      toastService.success(`League creator invitation sent to ${inviteEmail}`);
      inviteEmail = '';
    } else {
      toastService.error('Failed to send invitation.');
    }
    isInviting = false;
  }
</script>

<div class="space-y-8">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <div class="card bg-base-100 border border-base-300 shadow-sm rounded-[32px] p-8 space-y-2">
      <div class="flex justify-between items-start">
        <div class="bg-primary/10 text-primary p-3 rounded-2xl">
          <Trophy size={24} />
        </div>
      </div>
      <div>
        <p class="text-xs font-black uppercase tracking-widest opacity-40">Total Leagues</p>
        <p class="text-4xl font-black">{leagues.length}</p>
      </div>
    </div>

    <div class="card bg-base-100 border border-base-300 shadow-sm rounded-[32px] p-8 space-y-2">
      <div class="flex justify-between items-start">
        <div class="bg-success/10 text-success p-3 rounded-2xl">
          <Users size={24} />
        </div>
      </div>
      <div>
        <p class="text-xs font-black uppercase tracking-widest opacity-40">Total Players</p>
        <p class="text-4xl font-black">{totalMembers}</p>
      </div>
    </div>

    <div class="card bg-base-100 border border-base-300 shadow-sm rounded-[32px] p-8 space-y-2">
      <div class="flex justify-between items-start">
        <div class="bg-warning/10 text-warning p-3 rounded-2xl">
          <Activity size={24} />
        </div>
      </div>
      <div>
        <p class="text-xs font-black uppercase tracking-widest opacity-40">System Status</p>
        <p class="text-4xl font-black">Healthy</p>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <section class="space-y-4">
      <h3 class="text-xs font-black uppercase tracking-widest opacity-40 px-1">Quick Actions</h3>
      <div class="grid grid-cols-1 gap-3">
        {#if user?.roles?.includes('admin') && featureFlagService.isEnabled('admin_invites')}
          <div class="card bg-base-100 border border-base-300 rounded-[28px] p-6 space-y-4 shadow-sm">
            <div>
              <p class="font-black">Invite League Creator</p>
              <p class="text-xs opacity-40">Grants permission to create and manage leagues</p>
            </div>
            <div class="join w-full">
              <input
                type="email"
                bind:value={inviteEmail}
                placeholder="email@example.com"
                class="input input-bordered join-item flex-grow focus:outline-none"
              />
              <button
                onclick={handleInvite}
                disabled={isInviting || !inviteEmail}
                class="btn btn-neutral join-item font-black uppercase tracking-widest px-8"
              >
                {#if isInviting}
                  <span class="loading loading-spinner loading-xs"></span>
                {:else}
                  <Send size={16} />
                {/if}
                Invite
              </button>
            </div>
          </div>
        {/if}

        {#if featureFlagService.isEnabled('league_creation')}
        <button 
          onclick={() => goto('/leagues/create')}
          class="flex items-center justify-between p-6 bg-base-100 border border-base-300 rounded-[28px] hover:border-primary/50 transition-all group"
        >
          <div class="flex items-center gap-4">
            <div class="bg-base-200 p-3 rounded-2xl group-hover:bg-primary group-hover:text-primary-content transition-colors">
              <Plus size={20} />
            </div>
            <div class="text-left">
              <p class="font-black">Create New League</p>
              <p class="text-xs opacity-40">Initialize a new competition</p>
            </div>
          </div>
          <ArrowRight size={18} class="opacity-0 group-hover:opacity-100 transition-all" />
        </button>
        {/if}

        <button 
          onclick={() => goto('/manage/platform')}
          class="flex items-center justify-between p-6 bg-base-100 border border-base-300 rounded-[28px] hover:border-primary/50 transition-all group"
        >
          <div class="flex items-center gap-4">
            <div class="bg-base-200 p-3 rounded-2xl group-hover:bg-primary group-hover:text-primary-content transition-colors">
              <Shield size={20} />
            </div>
            <div class="text-left">
              <p class="font-black">Platform Settings</p>
              <p class="text-xs opacity-40">Manage global feature flags</p>
            </div>
          </div>
          <ArrowRight size={18} class="opacity-0 group-hover:opacity-100 transition-all" />
        </button>
      </div>
    </section>

    <section class="space-y-4">
      <h3 class="text-xs font-black uppercase tracking-widest opacity-40 px-1">Recent Leagues</h3>
      <div class="card bg-base-100 border border-base-300 shadow-sm rounded-[32px] overflow-hidden">
        <div class="divide-y divide-base-200">
          {#each leagues.slice(0, 5) as league}
            <button 
              onclick={() => goto(`/leagues/${league.id}/admin`)}
              class="w-full p-6 flex items-center justify-between hover:bg-base-200/50 transition-colors text-left"
            >
              <div>
                <p class="font-black">{league.name}</p>
                <p class="text-[10px] font-bold opacity-40 uppercase tracking-widest">{league.members?.length || 0} Members</p>
              </div>
              <ChevronRight size={16} class="opacity-20" />
            </button>
          {/each}
        </div>
        {#if leagues.length > 5}
          <button 
            onclick={() => goto('/manage/leagues')}
            class="w-full p-4 bg-base-200/50 text-center text-xs font-black uppercase tracking-widest opacity-60 hover:opacity-100 transition-all"
          >
            View All Leagues
          </button>
        {/if}
      </div>
    </section>
  </div>
</div>
