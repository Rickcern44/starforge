<script lang="ts">
  import { goto } from '$app/navigation';
  import { authService } from '$lib/services/auth.svelte';

  let email = $state('');
  let password = $state('');
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);

  // Redirect if already logged in
  $effect(() => {
    if (authService.user) {
      goto('/');
    }
  });

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    isLoading = true;
    errorMessage = null;
    try {
      const success = await authService.login(email, password);
      if (success) {
        goto('/');
      } else {
        errorMessage = 'Invalid email or password.';
      }
    } catch (err) {
      errorMessage = (err as Error).message;
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Login | League Manager</title>
</svelte:head>

<div class="flex items-center justify-center min-h-screen">
  <div class="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-lg">
    <div class="text-center">
      <h2 class="text-3xl font-extrabold text-gray-900">
        Welcome Back!
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Sign in to continue to League Manager
      </p>
    </div>

    <form class="mt-8 space-y-6" onsubmit={handleLogin}>
      <div class="space-y-4 rounded-md shadow-sm">
        <div>
          <label for="email-address" class="sr-only">Email address</label>
          <input
            id="email-address"
            name="email"
            type="email"
            bind:value={email}
            required
            class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Email address"
          />
        </div>
        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            bind:value={password}
            required
            class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
          />
        </div>
      </div>

      {#if errorMessage}
        <p class="text-sm text-red-600">{errorMessage}</p>
      {/if}

      <div class="flex items-center justify-between">
        <div class="text-sm">
          <a href="/auth/forgot-password" class="font-medium text-indigo-600 hover:text-indigo-500">
            Forgot your password?
          </a>
        </div>
      </div>

      <div>
        <button
          type="submit"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          disabled={isLoading}
        >
          {#if isLoading}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {/if}
          {isLoading ? 'Signing in...' : 'Sign in'}
        </button>
      </div>
    </form>
    <div class="text-sm text-center">
        <p>Don't have an account? <a href="/auth/register" class="font-medium text-indigo-600 hover:text-indigo-500">Sign up</a></p>
    </div>
  </div>
</div>
