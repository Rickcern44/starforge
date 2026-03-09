<script lang="ts">
  import { createLeague } from '$lib/services/league';
  import { toastService } from '$lib/services/toast.svelte';
  import { goto } from '$app/navigation';
  import { ChevronLeft, Info, LayoutGrid, Rocket } from 'lucide-svelte';

  let name = $state('');
  let isLoading = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!name) return;

    isLoading = true;
    try {
      const league = await createLeague(name);
      if (league) {
        toastService.success(`League "${name}" created successfully!`);
        goto(`/leagues/${league.id}/admin`);
      } else {
        toastService.error('Failed to create league.');
      }
    } catch (err) {
      toastService.error('An error occurred.');
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>New League | Bouncy</title>
</svelte:head>

<div class="max-w-md mx-auto py-12 px-4 space-y-8">
  <header>
    <button onclick={() => history.back()} class="btn btn-ghost btn-sm gap-2 opacity-40 hover:opacity-100 mb-4">
      <ChevronLeft size={16} />
      Back
    </button>
    <h1 class="text-3xl font-black tracking-tight leading-tight">Create a League</h1>
    <p class="text-sm font-bold opacity-40 uppercase tracking-widest mt-1">Setup your new organization</p>
  </header>

  <div class="card bg-base-100 shadow-xl border border-base-300 rounded-[32px]">
    <div class="card-body p-8">
      <form onsubmit={handleSubmit} class="space-y-6">
        <div class="form-control w-full">
          <label class="label" for="league-name">
            <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">League Name</span>
          </label>
          <div class="relative">
            <input
              id="league-name"
              type="text"
              bind:value={name}
              placeholder="e.g. Tuesday Night Hoops"
              required
              class="input input-bordered focus:outline-none w-full pl-10"
            />
            <LayoutGrid size={18} class="absolute left-3 top-1/2 -translate-y-1/2 opacity-20" />
          </div>
          <label class="label">
            <span class="label-text-alt opacity-40 italic">This will be the primary name for your players.</span>
          </label>
        </div>

        <div class="pt-4">
          <button
            type="submit"
            class="btn btn-neutral btn-block font-black uppercase tracking-widest shadow-lg shadow-neutral/20 h-auto py-4"
            disabled={isLoading || !name}
          >
            {#if isLoading}
              <span class="loading loading-spinner"></span>
            {:else}
              <Rocket size={18} />
            {/if}
            Create League
          </button>
        </div>
      </form>
    </div>
  </div>

  <div class="alert bg-base-300 border-none shadow-lg rounded-2xl">
    <Info size={20} class="opacity-40" />
    <div class="text-xs font-bold leading-tight opacity-60">
      As the creator, you will automatically be assigned as the <strong>League Owner</strong>.
    </div>
  </div>
</div>
