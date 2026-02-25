<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let user = $derived(authService.user);

  function handleLogout() {
    authService.logout();
    goto('/auth/logout');
  }
</script>

<header class="bg-white border-b border-gray-100 sticky top-0 z-30">
  <nav class="container mx-auto px-4 py-4 flex justify-between items-center">
    <a href="/" class="text-xl font-extrabold tracking-tight text-indigo-600">Bouncy</a>
    <div class="flex items-center space-x-3">
      {#if user}
        <div class="flex items-center space-x-4">
          <button 
            onclick={handleLogout} 
            class="text-sm font-bold text-gray-500 hover:text-red-500 transition-colors duration-200"
          >
            Logout
          </button>
          <div class="h-8 w-8 rounded-full bg-indigo-50 border border-indigo-100 flex items-center justify-center text-indigo-600 font-bold text-xs uppercase">
            {user?.name?.charAt(0) || 'U'}
          </div>
        </div>
      {:else if !$page.url.pathname.startsWith('/auth')}
        <a href="/auth/login" class="text-sm font-bold text-indigo-600 hover:text-indigo-800 transition-colors">Login</a>
      {/if}
    </div>
  </nav>
</header>
