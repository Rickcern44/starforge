<script lang="ts">
  import { toastService } from '$lib/services/toast.svelte';
  import { flip } from 'svelte/animate';
  import { fade, fly } from 'svelte/transition';

  const alertTypes = {
    success: 'alert-success',
    error: 'alert-error',
    info: 'alert-info'
  };
</script>

<div class="toast toast-end toast-bottom z-[100] p-6 space-y-2">
  {#each toastService.toasts as toast (toast.id)}
    <div
      animate:flip={{ duration: 300 }}
      in:fly={{ y: 20, opacity: 0, duration: 400 }}
      out:fade={{ duration: 200 }}
      class="alert {alertTypes[toast.type]} shadow-xl rounded-2xl min-w-[280px] border-none text-white font-bold"
    >
      <span>{toast.message}</span>
      <button 
        onclick={() => toastService.remove(toast.id)}
        class="btn btn-ghost btn-xs btn-circle"
      >
        ✕
      </button>
    </div>
  {/each}
</div>
