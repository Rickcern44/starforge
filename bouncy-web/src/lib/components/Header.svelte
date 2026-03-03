<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';

  let user = $derived(authService.user);

  function handleLogout() {
    authService.logout();
    goto('/auth/logout');
  }
</script>

<header class="bg-white border-b border-gray-100 sticky top-0 z-40 backdrop-blur-md bg-white/80">
  <nav class="container mx-auto px-4 py-3 flex justify-between items-center h-16">
    <a href="/" class="flex items-center space-x-2 group">
      <div class="w-8 h-8 bg-indigo-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-indigo-200 group-active:scale-95 transition-transform">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
        </svg>
      </div>
      <span class="text-xl font-black tracking-tight text-gray-900 group-hover:text-indigo-600 transition-colors">Bouncy</span>
    </a>
    <div class="flex items-center space-x-2">
      {#if user}
        <div class="flex items-center space-x-3">
          <div class="hidden sm:block text-right">
            <p class="text-xs font-bold text-gray-900 leading-none">{user.name}</p>
            <button 
              onclick={handleLogout} 
              class="text-[10px] font-black text-gray-400 hover:text-red-500 uppercase tracking-widest transition-colors"
            >
              Sign Out
            </button>
          </div>
          <button 
            onclick={() => goto('/auth/logout')}
            class="h-10 w-10 rounded-2xl bg-indigo-50 border-2 border-white shadow-sm flex items-center justify-center text-indigo-600 font-black text-sm uppercase ring-1 ring-indigo-100 group active:scale-95 transition-all"
          >
            {user?.name?.charAt(0) || 'U'}
          </button>
        </div>
      {:else if !page.url.pathname.startsWith('/auth')}
        <a href="/auth/login" class="px-4 py-2 bg-indigo-600 text-white text-sm font-black rounded-xl shadow-lg shadow-indigo-100 active:scale-95 transition-all">
          Sign In
        </a>
      {/if}
    </div>
  </nav>
</header>
