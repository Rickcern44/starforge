<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import { featureFlagService } from '$lib/services/feature-flag.svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import ThemeToggle from './ThemeToggle.svelte';
  import { Trophy, Wallet, LayoutDashboard, LogOut, Bell, Shield } from 'lucide-svelte';
  import logo from '$lib/assets/logo.png';

  let user = $derived(authService.user);

  function handleLogout() {
    authService.logout();
    goto('/auth/logout');
  }
</script>

<header class="navbar bg-base-100 border-b border-base-300 sticky top-0 z-40 backdrop-blur-md bg-base-100/80">
  <div class="flex-1">
    <a href="/" class="btn btn-ghost group gap-2 px-2 sm:px-4">
      <div class="h-10 flex items-center justify-center overflow-hidden group-active:scale-95 transition-transform">
        <img src={logo} alt="Bouncy" class="h-full w-auto object-contain dark:invert" />
      </div>
    </a>
  </div>
  
  <div class="flex-none gap-2 px-2">
    {#if user && featureFlagService.isEnabled('notifications')}
      <button class="btn btn-ghost btn-circle">
        <div class="indicator">
          <Bell size={20} />
          <span class="badge badge-xs badge-primary indicator-item"></span>
        </div>
      </button>
    {/if}

    <ThemeToggle />
    {#if user}
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar border border-base-300">
          <div class="w-10 rounded-xl bg-base-200 text-base-content flex items-center justify-center font-black">
             {user?.name?.charAt(0) || 'U'}
          </div>
        </div>
        <ul tabindex="0" class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52 border border-base-300">
          <li class="menu-title text-[10px] font-black uppercase opacity-40 px-4 py-2 border-b border-base-300 mb-1">{user.name}</li>
          <li>
            <button onclick={() => goto('/')} class="py-3">
              <LayoutDashboard size={16} />
              Dashboard
            </button>
          </li>
          {#if featureFlagService.isEnabled('payments')}
          <li>
            <button onclick={() => goto('/ledger')} class="py-3">
              <Wallet size={16} />
              My Wallet
            </button>
          </li>
          {/if}
          {#if user.roles?.includes('admin') || user.roles?.includes('league_creator')}
          <li>
            <button onclick={() => goto('/manage')} class="py-3 text-primary font-black">
              <Shield size={16} />
              Management Portal
            </button>
          </li>
          {/if}
          <div class="divider my-0"></div>
          <li>
            <button onclick={handleLogout} class="text-error py-3">
              <LogOut size={16} />
              Sign Out
            </button>
          </li>
        </ul>
      </div>
    {:else if !page.url.pathname.startsWith('/auth')}
      <a href="/auth/login" class="btn btn-primary btn-sm rounded-xl font-black uppercase tracking-widest px-6">
        Sign In
      </a>
    {/if}
  </div>
</header>
