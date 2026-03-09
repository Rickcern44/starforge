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

<div class="flex items-center justify-center min-h-[calc(100vh-80px)]">
  <div class="card w-full max-w-md bg-base-100 shadow-xl border border-base-200">
    <div class="card-body p-8 space-y-6">
      <div class="text-center">
        <h2 class="text-3xl font-black tracking-tight">
          Welcome Back!
        </h2>
        <p class="mt-2 text-sm opacity-60 font-bold uppercase tracking-wider">
          Sign in to your account
        </p>
      </div>

      <form class="space-y-4" onsubmit={handleLogin}>
        <div class="form-control w-full">
          <label class="label" for="email-address">
            <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Email address</span>
          </label>
          <input
            id="email-address"
            name="email"
            type="email"
            bind:value={email}
            required
            class="input input-bordered focus:outline-none w-full"
            placeholder="name@example.com"
          />
        </div>
        
        <div class="form-control w-full">
          <label class="label" for="password">
            <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Password</span>
          </label>
          <input
            id="password"
            name="password"
            type="password"
            bind:value={password}
            required
            class="input input-bordered focus:outline-none w-full"
            placeholder="••••••••"
          />
        </div>

        {#if errorMessage}
          <div class="alert alert-error text-xs font-bold py-2 rounded-xl border-none">
            {errorMessage}
          </div>
        {/if}

        <div class="flex items-center justify-between pt-2">
          <button class="btn btn-link btn-xs no-underline hover:no-underline font-black uppercase tracking-widest opacity-40" onclick={(e) => { e.preventDefault(); alert('Password recovery is coming soon!'); }}>
            Forgot password?
          </button>
        </div>

        <div class="pt-4">
          <button
            type="submit"
            class="btn btn-primary btn-block font-black uppercase tracking-widest shadow-lg shadow-primary/20 h-auto py-4"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="loading loading-spinner"></span>
            {/if}
            {isLoading ? 'Signing in...' : 'Sign in'}
          </button>
        </div>
      </form>
      
      <div class="divider text-[10px] font-black uppercase tracking-widest opacity-20">or</div>
      
      <div class="text-sm text-center">
          <p class="font-bold opacity-60 uppercase tracking-widest text-[10px]">New to Bouncy?</p>
          <a href="/auth/register" class="btn btn-ghost btn-sm btn-block mt-2 font-black uppercase tracking-widest">Create an Account</a>
      </div>
    </div>
  </div>
</div>
