<script lang="ts">
  import { goto } from '$app/navigation';
  import { authService } from '$lib/services/auth.svelte';

  let name = $state('');
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    if (password !== confirmPassword) {
      errorMessage = 'Passwords do not match.';
      return;
    }
    isLoading = true;
    errorMessage = null;
    try {
      const success = await authService.register(name, email, password);
      if (success) {
        goto('/');
      } else {
        errorMessage = 'Registration failed. Please try again.';
      }
    } catch (err) {
      errorMessage = (err as Error).message;
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Register | League Manager</title>
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-80px)] py-12">
  <div class="card w-full max-w-md bg-base-100 shadow-xl border border-base-200">
    <div class="card-body p-8 space-y-6">
      <div class="text-center">
        <h2 class="text-3xl font-black tracking-tight">
          Join Bouncy
        </h2>
        <p class="mt-2 text-sm opacity-60 font-bold uppercase tracking-wider">
          Create your player account
        </p>
      </div>

      <form class="space-y-4" onsubmit={handleRegister}>
        <div class="form-control w-full">
          <label class="label" for="name">
            <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Full Name</span>
          </label>
          <input
            id="name"
            name="name"
            type="text"
            bind:value={name}
            required
            class="input input-bordered focus:outline-none w-full"
            placeholder="John Doe"
          />
        </div>

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
            minlength="8"
            class="input input-bordered focus:outline-none w-full"
            placeholder="••••••••"
          />
        </div>

        <div class="form-control w-full">
          <label class="label" for="confirm-password">
            <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Confirm Password</span>
          </label>
          <input
            id="confirm-password"
            name="confirm-password"
            type="password"
            bind:value={confirmPassword}
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

        <div class="pt-4">
          <button
            type="submit"
            class="btn btn-primary btn-block font-black uppercase tracking-widest shadow-lg shadow-primary/20 h-auto py-4"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="loading loading-spinner"></span>
            {/if}
            {isLoading ? 'Creating account...' : 'Create Account'}
          </button>
        </div>
      </form>
      
      <div class="divider text-[10px] font-black uppercase tracking-widest opacity-20">or</div>
      
      <div class="text-sm text-center">
          <p class="font-bold opacity-60 uppercase tracking-widest text-[10px]">Already have an account?</p>
          <a href="/auth/login" class="btn btn-ghost btn-sm btn-block mt-2 font-black uppercase tracking-widest">Sign In Instead</a>
      </div>
    </div>
  </div>
</div>
