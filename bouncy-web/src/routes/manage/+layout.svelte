<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { LayoutGrid, Shield, Users, Trophy } from 'lucide-svelte';

  let { children } = $props();
  
  const tabs = [
    { id: 'leagues', name: 'Leagues', icon: Trophy, path: '/manage/leagues' },
    { id: 'platform', name: 'Platform', icon: Shield, path: '/manage/platform' },
    { id: 'users', name: 'Users', icon: Users, path: '/manage/users' }
  ];

  let currentTab = $derived(page.url.pathname.split('/')[2] || 'leagues');
</script>

<div class="max-w-6xl mx-auto py-8 px-4 space-y-8">
  <header>
    <h1 class="text-4xl font-black tracking-tight leading-tight">Management</h1>
    <p class="text-sm font-bold opacity-40 uppercase tracking-wider">System & League Administration</p>
  </header>

  <div class="flex flex-col md:flex-row gap-8">
    <aside class="w-full md:w-64 space-y-1">
      {#each tabs as tab}
        <button
          onclick={() => goto(tab.path)}
          class="btn btn-ghost w-full justify-start gap-3 rounded-2xl px-4 py-3 h-auto font-black {currentTab === tab.id ? 'bg-primary/10 text-primary hover:bg-primary/20' : 'opacity-60 hover:opacity-100 hover:bg-base-200'}"
        >
          <tab.icon size={20} />
          {tab.name}
        </button>
      {/each}
    </aside>

    <main class="flex-1">
      {@render children()}
    </main>
  </div>
</div>
