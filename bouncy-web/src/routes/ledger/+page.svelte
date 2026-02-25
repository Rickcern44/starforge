<script lang="ts">
  import { authService } from '$lib/services/auth.svelte';
  import { goto } from '$app/navigation';

  // Mock data to match the Flutter implementation
  let items = $state([
    {
      id: '1',
      type: 'charge',
      title: 'Game vs Red Dragons',
      date: 'Oct 10, 2023',
      amount: 15.0,
      status: 'unpaid'
    },
    {
      id: '2',
      type: 'charge',
      title: 'Game vs Blue Hawks',
      date: 'Oct 17, 2023',
      amount: 15.0,
      status: 'unpaid'
    },
    {
      id: '3',
      type: 'payment',
      title: 'Bulk Payment',
      date: 'Oct 05, 2023',
      amount: 45.0,
      status: 'completed'
    },
    {
      id: '4',
      type: 'charge',
      title: 'Court Rental Fee',
      date: 'Oct 01, 2023',
      amount: 10.0,
      status: 'paid'
    },
  ]);

  let selectedChargeIds = $state(new Set<string>());

  let selectedTotal = $derived.by(() => {
    let total = 0;
    items.forEach(item => {
      if (selectedChargeIds.has(item.id)) {
        total += item.amount;
      }
    });
    return total;
  });

  let outstandingBalance = $derived.by(() => {
    return items
      .filter(i => i.type === 'charge' && i.status === 'unpaid')
      .reduce((sum, i) => sum + i.amount, 0);
  });

  function toggleSelection(id: string) {
    if (selectedChargeIds.has(id)) {
      selectedChargeIds.delete(id);
    } else {
      selectedChargeIds.add(id);
    }
    // We need to re-assign to trigger reactivity for the Set in Svelte 5 
    // if not using a specialized reactive collection, but $state(new Set()) 
    // works if we replace the set or if Svelte 5 handles set mutations (it doesn't by default without Proxy).
    // Actually, Svelte 5 $state with Set requires reassignment or using a reactive wrapper.
    selectedChargeIds = new Set(selectedChargeIds);
  }

  function handleMakePayment() {
    alert(`Payment of $${selectedTotal.toFixed(2)} processed!`);
    // Mock update: mark selected as paid
    items = items.map(item => {
      if (selectedChargeIds.has(item.id)) {
        return { ...item, status: 'paid' };
      }
      return item;
    });
    selectedChargeIds = new Set();
  }
</script>

<svelte:head>
  <title>League Ledger | League Manager</title>
</svelte:head>

<div class="max-w-3xl mx-auto flex flex-col min-h-[calc(100-rem)] pb-24">
  <div class="flex items-center space-x-4 mb-6">
    <button onclick={() => history.back()} class="text-gray-500 hover:text-indigo-600 transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
    <h1 class="text-3xl font-bold text-gray-800">League Ledger</h1>
  </div>

  <!-- Balance Header -->
  <div class="bg-indigo-50 border border-indigo-100 rounded-xl p-8 mb-8 text-center shadow-sm">
    <p class="text-sm font-semibold text-indigo-600 uppercase tracking-wider">Outstanding Balance</p>
    <p class="text-5xl font-black text-red-600 mt-2">${outstandingBalance.toFixed(2)}</p>
  </div>

  <!-- Items List -->
  <div class="bg-white rounded-xl shadow-md overflow-hidden border border-gray-100">
    <div class="divide-y divide-gray-100">
      {#each items as item}
        {@const isCharge = item.type === 'charge'}
        {@const isUnpaid = item.status === 'unpaid'}
        {@const isSelected = selectedChargeIds.has(item.id)}

        <div class="p-4 flex items-center space-x-4 hover:bg-gray-50 transition-colors duration-150">
          <div class="p-3 rounded-full {isCharge ? (isUnpaid ? 'bg-orange-100 text-orange-600' : 'bg-red-100 text-red-600') : 'bg-green-100 text-green-600'}">
            {#if isCharge}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            {/if}
          </div>

          <div class="flex-grow">
            <h4 class="font-bold text-gray-800">{item.title}</h4>
            <p class="text-sm text-gray-500">{item.date}</p>
          </div>

          <div class="flex items-center space-x-4">
            <p class="font-black {isCharge ? 'text-red-600' : 'text-green-600'}">
              {isCharge ? '-' : '+'}${item.amount.toFixed(2)}
            </p>
            
            {#if isUnpaid}
              <input 
                type="checkbox" 
                checked={isSelected}
                onchange={() => toggleSelection(item.id)}
                class="h-6 w-6 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Floating Payment Footer -->
  {#if selectedChargeIds.size > 0}
    <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-[0_-4px_10px_rgba(0,0,0,0.05)] p-4 z-20">
      <div class="max-w-3xl mx-auto flex items-center justify-between">
        <div>
          <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Selected Total</p>
          <p class="text-2xl font-black text-gray-900">${selectedTotal.toFixed(2)}</p>
        </div>
        <button 
          onclick={handleMakePayment}
          class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-8 rounded-xl shadow-lg transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98]"
        >
          Make Payment
        </button>
      </div>
    </div>
  {/if}
</div>
