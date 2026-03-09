<script lang="ts">
  import { Users, Info } from 'lucide-svelte';
</script>

<div class="space-y-8">
  <header>
    <h2 class="text-3xl font-black tracking-tight leading-tight">Users</h2>
    <p class="text-sm font-bold opacity-40 uppercase tracking-wider">User Administration & RBAC</p>
  </header>

  <div class="card bg-base-100 p-12 border border-dashed border-base-300 items-center justify-center text-center opacity-40">
    <Users size={48} class="mb-4" />
    <p class="font-black uppercase tracking-widest">User management coming soon</p>
    <p class="text-sm">Manage global roles and permissions here.</p>
  </div>

  <div class="alert bg-base-300 border-none shadow-lg rounded-2xl">
    <Info size={20} class="opacity-40" />
    <div class="text-xs font-bold leading-tight opacity-60">
      Global role management allows you to assign system-wide administrative privileges.
    </div>
  </div>
</div>
