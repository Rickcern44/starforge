<script lang="ts">
  import { toggleFeatureFlag } from '$lib/services/feature-flag';
  import { toastService } from '$lib/services/toast.svelte';
  import type { FeatureFlag } from '$lib/models';
  import { ChevronLeft, Info, Activity } from 'lucide-svelte';

  let { data } = $props();
  let features = $state<FeatureFlag[]>(data.features || []);

  async function handleToggle(key: string, currentStatus: boolean) {
    const newStatus = !currentStatus;
    const success = await toggleFeatureFlag(key, newStatus);
    
    if (success) {
      toastService.success(`Feature ${key} updated`);
      // Update local state
      features = features.map(f => f.key === key ? { ...f, enabled: newStatus } : f);
    } else {
      toastService.error(`Failed to update feature ${key}`);
    }
  }
</script>

<svelte:head>
  <title>Features | Admin</title>
</svelte:head>

<div class="space-y-8">
  <header class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-black tracking-tight leading-tight">Platform Features</h2>
      <p class="text-sm font-bold opacity-40 uppercase tracking-wider">System Administration</p>
    </div>
  </header>

  <div class="card bg-base-100 rounded-[32px] border border-base-300 shadow-xl overflow-hidden">
    {#if features.length === 0}
      <div class="p-12 text-center opacity-20">
        <Activity size={48} class="mx-auto mb-4" />
        <p class="italic">No feature flags found.</p>
      </div>
    {:else}
      <div class="divide-y divide-base-200">
        {#each features as feature}
          <div class="p-6 flex items-center justify-between hover:bg-base-200/50 transition-colors">
            <div class="space-y-1 pr-4 text-left">
              <h3 class="font-black flex items-center">
                {feature.name}
                <span class="ml-2 text-[10px] font-mono bg-base-200 px-1.5 py-0.5 rounded opacity-40">
                  {feature.key}
                </span>
              </h3>
              <p class="text-sm opacity-60 leading-snug">{feature.description}</p>
            </div>
            
            <div class="flex-shrink-0">
              <input 
                type="checkbox" 
                class="toggle toggle-md" 
                checked={feature.enabled} 
                onchange={() => handleToggle(feature.key, feature.enabled)}
              />
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="alert bg-base-300 border-none shadow-lg rounded-2xl">
    <Info size={20} class="opacity-40" />
    <div class="text-xs font-bold leading-tight opacity-60">
      Changes to feature flags take effect immediately. Use with caution.
    </div>
  </div>
</div>
