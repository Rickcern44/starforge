<script lang="ts">
  import Header from '$lib/components/Header.svelte';
  import Toaster from '$lib/components/Toaster.svelte';
  import './layout.css';
  import { authService } from '$lib/services/auth.svelte';
  import { browser } from '$app/environment';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';

  let { children } = $props();

  $effect(() => {
    if (browser) {
      const path = page.url.pathname;
      const isAuthPage = path.startsWith('/auth');

      // If user is not logged in and not on an auth page, redirect to login
      if (!authService.user && !isAuthPage) {
        goto('/auth/login');
      }

      // If logged in and on an auth page (except logout), redirect to home
      if (authService.user && isAuthPage && path !== '/auth/logout') {
        goto('/');
      }
    }
  });
</script>

<div class="min-h-screen bg-gray-100">
  <Toaster />
  {#if authService.user || page.url.pathname.startsWith('/auth') || page.url.pathname === '/'}
    <Header />
    <main class="container mx-auto p-4">
      {@render children()}
    </main>
  {:else}
    <div class="flex justify-center items-center h-screen">
      <div class="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600"></div>
    </div>
  {/if}
</div>
