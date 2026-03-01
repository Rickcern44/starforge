<script lang="ts">
  import Header from '$lib/components/Header.svelte';
  import './layout.css';
  import { authService } from '$lib/services/auth.svelte';
  import { browser } from '$app/environment';

  let { children, data } = $props();

  $effect(() => {
    if (browser) {
      if (data.user) {
        authService.user = data.user;
        if (data.token) authService.token = data.token;
      }
    }
  });
</script>

<div class="min-h-screen bg-gray-100">
  <Header />
  <main class="container mx-auto p-4">
    {@render children()}
  </main>
</div>
