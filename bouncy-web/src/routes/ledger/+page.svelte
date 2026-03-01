<script lang="ts">
  import type { Payment, GameCharge } from '$lib/models';

  let { data } = $props();

  function getUnpaidAmount(charge: GameCharge) {
    const paid = charge.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return (charge.amountCents - paid) / 100;
  }

  let unpaidCharges = $derived.by(() => {
    if (!data.charges) return [];
    return data.charges
      .map(c => ({
        ...c,
        unpaid: getUnpaidAmount(c)
      }))
      .filter(c => c.unpaid > 0)
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  });

  let recentActivity = $derived.by(() => {
    const items: any[] = [];
    
    if (data.charges) {
      data.charges.forEach((c: GameCharge) => {
        items.push({
          id: c.id,
          type: 'charge',
          title: c.externalName || `Game ${c.gameId.substring(0, 5)}`,
          date: new Date(c.createdAt).toLocaleDateString(),
          rawDate: new Date(c.createdAt),
          amount: c.amountCents / 100,
          status: getUnpaidAmount(c) > 0 ? 'unpaid' : 'paid'
        });
      });
    }

    if (data.payments) {
      data.payments.forEach((p: Payment) => {
        items.push({
          id: p.id,
          type: 'payment',
          title: 'Payment Received',
          date: new Date(p.receivedAt).toLocaleDateString(),
          rawDate: new Date(p.receivedAt),
          amount: p.amountInCents / 100,
          status: 'completed'
        });
      });
    }

    return items.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
  });

  let outstandingBalance = $derived.by(() => {
    return unpaidCharges.reduce((sum, i) => sum + i.unpaid, 0);
  });
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

  <!-- Unpaid Charges Section -->
  <section class="mb-10">
    <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest mb-4 px-1">Unpaid Charges</h3>
    <div class="bg-white rounded-3xl shadow-sm overflow-hidden border border-gray-100">
      {#if unpaidCharges.length === 0}
        <div class="p-8 text-center">
          <p class="text-gray-400 font-bold italic text-sm">No unpaid charges. You're all caught up!</p>
        </div>
      {:else}
        <div class="divide-y divide-gray-50">
          {#each unpaidCharges as charge}
            <div class="p-4 flex items-center space-x-4 hover:bg-gray-50 transition-colors">
              <div class="p-3 rounded-2xl bg-orange-50 text-orange-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="flex-grow">
                <h4 class="font-bold text-gray-800">{charge.externalName || `Game ${charge.gameId.substring(0, 5)}`}</h4>
                <p class="text-[10px] text-gray-400 font-black uppercase tracking-widest">{new Date(charge.createdAt).toLocaleDateString()}</p>
              </div>
              <div class="text-right">
                <p class="text-[10px] text-gray-400 uppercase font-black tracking-widest leading-none mb-1">Due</p>
                <p class="text-lg font-black text-red-600">${charge.unpaid.toFixed(2)}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>

  <!-- Recent Activity Section -->
  <section>
    <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest mb-4 px-1">Recent Activity</h3>
    <div class="bg-white rounded-3xl shadow-sm overflow-hidden border border-gray-100">
      {#if recentActivity.length === 0}
        <div class="p-8 text-center">
          <p class="text-gray-400 font-bold italic text-sm">No recent activity recorded.</p>
        </div>
      {:else}
        <div class="divide-y divide-gray-50">
          {#each recentActivity as item}
            {@const isCharge = item.type === 'charge'}
            <div class="p-4 flex items-center space-x-4 hover:bg-gray-50 transition-colors">
              <div class="p-2.5 rounded-xl {isCharge ? 'bg-gray-50 text-gray-400' : 'bg-green-50 text-green-600'}">
                {#if isCharge}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8V9m0 3v-1m0 4v-1m0-8a9 9 0 110 18 9 9 0 010-18z" />
                  </svg>
                {/if}
              </div>
              <div class="flex-grow">
                <h4 class="text-sm font-bold text-gray-700">{item.title}</h4>
                <p class="text-[9px] text-gray-400 font-bold uppercase tracking-widest">{item.date}</p>
              </div>
              <div class="text-right">
                <p class="font-black {isCharge ? 'text-gray-400' : 'text-green-600'}">
                  {isCharge ? '-' : '+'}${item.amount.toFixed(2)}
                </p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>
</div>
