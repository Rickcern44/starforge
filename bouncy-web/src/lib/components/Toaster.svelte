<script lang="ts">
  import { toastService } from '$lib/services/toast.svelte';
  import { flip } from 'svelte/animate';
  import { fade, fly } from 'svelte/transition';

  const typeStyles = {
    success: 'bg-green-600 text-white shadow-green-100',
    error: 'bg-red-600 text-white shadow-red-100',
    info: 'bg-indigo-600 text-white shadow-indigo-100'
  };
</script>

<div class="fixed bottom-6 right-6 z-[100] flex flex-col items-end space-y-3 pointer-events-none">
  {#each toastService.toasts as toast (toast.id)}
    <div
      animate:flip={{ duration: 300 }}
      in:fly={{ y: 20, opacity: 0, duration: 400 }}
      out:fade={{ duration: 200 }}
      class="pointer-events-auto px-6 py-4 rounded-2xl shadow-xl flex items-center space-x-3 min-w-[280px] max-w-md {typeStyles[toast.type]}"
    >
      {#if toast.type === 'success'}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      {:else if toast.type === 'error'}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      {/if}
      
      <p class="text-sm font-bold leading-tight">{toast.message}</p>
      
      <button 
        onclick={() => toastService.remove(toast.id)}
        class="ml-auto p-1 hover:bg-black/10 rounded-lg transition-colors"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 opacity-60" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/each}
</div>
