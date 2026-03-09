<script lang="ts">
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';
  import { Sun, Moon } from 'lucide-svelte';

  let theme = $state('light');

  onMount(() => {
    if (browser) {
      const storedTheme = localStorage.getItem('theme');
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      
      theme = storedTheme || (prefersDark ? 'dark' : 'light');
      document.documentElement.setAttribute('data-theme', theme);
    }
  });

  function toggleTheme() {
    theme = theme === 'light' ? 'dark' : 'light';
    if (browser) {
      localStorage.setItem('theme', theme);
      document.documentElement.setAttribute('data-theme', theme);
    }
  }
</script>

<label class="swap swap-rotate h-10 w-10 rounded-xl bg-base-200 border border-base-300 hover:bg-base-300 transition-colors">
  <!-- this hidden checkbox controls the state -->
  <input type="checkbox" class="theme-controller" checked={theme === 'dark'} onchange={toggleTheme} />

  <!-- sun icon -->
  <div class="swap-on">
    <Sun size={20} />
  </div>

  <!-- moon icon -->
  <div class="swap-off">
    <Moon size={20} />
  </div>
</label>
