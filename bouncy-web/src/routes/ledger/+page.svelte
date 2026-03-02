<script lang="ts">
  import type { Payment, GameCharge, League } from '$lib/models';

  let { data } = $props();

  function getUnpaidAmount(charge: GameCharge) {
    const paid = charge.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return (charge.amountCents - paid) / 100;
  }

  let expandedLeagues = $state(new Set<string>());

  function toggleExpand(leagueId: string) {
    if (expandedLeagues.has(leagueId)) {
      expandedLeagues.delete(leagueId);
    } else {
      expandedLeagues.add(leagueId);
    }
    expandedLeagues = new Set(expandedLeagues);
  }

  let leaguesWithData = $derived.by(() => {
    const leagueMap = new Map<string, { 
      league: League, 
      unpaidCharges: any[], 
      allCharges: any[],
      recentActivity: any[],
      balance: number 
    }>();

    // Initialize map with all user leagues
    if (data.leagues) {
      data.leagues.forEach((l: League) => {
        leagueMap.set(l.id, { 
          league: l, 
          unpaidCharges: [], 
          allCharges: [],
          recentActivity: [],
          balance: 0 
        });
      });
    }

    // Process Charges
    if (data.charges) {
      data.charges.forEach((c: GameCharge) => {
        const leagueId = c.game?.leagueId || 'unknown';
        if (!leagueMap.has(leagueId)) return;

        const unpaid = getUnpaidAmount(c);
        const chargeItem = {
          id: c.id,
          type: 'charge',
          title: c.game ? `Game at ${c.game.location}` : `Game ${c.gameId.substring(0, 5)}`,
          date: new Date(c.createdAt).toLocaleDateString(),
          rawDate: new Date(c.createdAt),
          amount: c.amountCents / 100,
          unpaid: unpaid,
          status: unpaid > 0 ? 'unpaid' : 'paid'
        };

        const entry = leagueMap.get(leagueId)!;
        entry.allCharges.push(chargeItem);
        entry.recentActivity.push(chargeItem);
        if (unpaid > 0) {
          entry.unpaidCharges.push(chargeItem);
          entry.balance += unpaid;
        }
      });
    }

    // Process Payments
    if (data.payments) {
      data.payments.forEach((p: Payment) => {
        const entry = leagueMap.get(p.leagueId);
        if (!entry) return;

        entry.recentActivity.push({
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

    // Sort recent activity for each league
    leagueMap.forEach(entry => {
      entry.recentActivity.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
      entry.unpaidCharges.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
      entry.allCharges.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
    });

    // Return as array, only including leagues that have some data or are active
    return Array.from(leagueMap.values())
      .filter(entry => entry.recentActivity.length > 0 || entry.league.isActive)
      .sort((a, b) => a.league.name.localeCompare(b.league.name));
  });

  let totalBalance = $derived(leaguesWithData.reduce((sum, l) => sum + l.balance, 0));
</script>

<svelte:head>
  <title>League Ledger | League Manager</title>
</svelte:head>

<div class="max-w-3xl mx-auto flex flex-col min-h-screen pb-24 px-4">
  <div class="flex items-center space-x-4 py-6">
    <button onclick={() => history.back()} class="text-gray-500 hover:text-indigo-600 transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
    <h1 class="text-3xl font-black text-gray-900 tracking-tight">Financial Ledger</h1>
  </div>

  <!-- Total Balance Header -->
  <div class="bg-indigo-600 rounded-[32px] p-8 mb-10 text-white shadow-xl shadow-indigo-100 flex justify-between items-center relative overflow-hidden">
    <div class="relative z-10">
      <p class="text-sm font-bold text-indigo-100 uppercase tracking-widest opacity-80">Total Outstanding</p>
      <p class="text-5xl font-black mt-2">${totalBalance.toFixed(2)}</p>
    </div>
    <div class="bg-white/10 p-4 rounded-2xl backdrop-blur-md relative z-10">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
    </div>
    <!-- Abstract background shape -->
    <div class="absolute -right-10 -bottom-10 w-40 h-40 bg-indigo-500 rounded-full blur-3xl opacity-50"></div>
  </div>

  {#each leaguesWithData as { league, unpaidCharges, allCharges, recentActivity, balance }}
    <div class="mb-12 space-y-6">
      <!-- League Header -->
      <div class="flex items-end justify-between px-1">
        <div>
          <h2 class="text-2xl font-black text-gray-900 tracking-tight">{league.name}</h2>
          <p class="text-[10px] font-black text-indigo-400 uppercase tracking-widest">League Specific Ledger</p>
        </div>
        <div class="text-right">
          <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">League Balance</p>
          <p class="text-xl font-black {balance > 0 ? 'text-red-600' : 'text-green-600'}">${balance.toFixed(2)}</p>
        </div>
      </div>

      <!-- Unpaid for this League (Actionable) -->
      {#if unpaidCharges.length > 0}
        <div class="space-y-3">
          <h3 class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Unpaid in {league.name}</h3>
          <div class="bg-white rounded-[24px] border border-gray-100 shadow-sm overflow-hidden">
            <div class="divide-y divide-gray-50">
              {#each unpaidCharges as charge}
                <div class="p-4 flex items-center space-x-4">
                  <div class="p-2.5 rounded-xl bg-orange-50 text-orange-600">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div class="flex-grow">
                    <h4 class="text-sm font-bold text-gray-800">{charge.title}</h4>
                    <p class="text-[9px] text-gray-400 font-bold uppercase tracking-widest">{charge.date}</p>
                  </div>
                  <div class="text-right">
                    <p class="text-[9px] text-gray-400 font-black uppercase tracking-widest leading-none mb-0.5">Due</p>
                    <p class="text-sm font-black text-red-600">${charge.unpaid.toFixed(2)}</p>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/if}

      <!-- History for this League (Summarized) -->
      <div class="space-y-3">
        <div class="flex justify-between items-center px-1">
          <h3 class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Recent Charges</h3>
          {#if allCharges.length > 5 || recentActivity.length > allCharges.length}
            <button 
              onclick={() => toggleExpand(league.id)}
              class="text-[9px] font-black text-indigo-600 uppercase tracking-widest hover:text-indigo-800 transition-colors"
            >
              {expandedLeagues.has(league.id) ? 'Show Less' : 'See all activity'}
            </button>
          {/if}
        </div>
        
        <div class="bg-white rounded-[24px] border border-gray-100 shadow-sm overflow-hidden">
          <div class="divide-y divide-gray-50">
            {#each (expandedLeagues.has(league.id) ? recentActivity : allCharges.slice(0, 5)) as item}
              {@const isCharge = item.type === 'charge'}
              <div class="p-4 flex items-center space-x-4 group transition-colors">
                <div class="p-2 rounded-lg {isCharge ? 'bg-gray-50 text-gray-400' : 'bg-green-50 text-green-600'}">
                  {#if isCharge}
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  {:else}
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8V9m0 3v-1m0 4v-1m0-8a9 9 0 110 18 9 9 0 010-18z" />
                    </svg>
                  {/if}
                </div>
                <div class="flex-grow">
                  <h4 class="text-[13px] font-bold text-gray-700">{item.title}</h4>
                  <p class="text-[9px] text-gray-400 font-bold uppercase tracking-widest">{item.date}</p>
                </div>
                <div class="text-right">
                  <p class="text-sm font-black {isCharge ? 'text-gray-400' : 'text-green-600'}">
                    {isCharge ? '-' : '+'}${item.amount.toFixed(2)}
                  </p>
                </div>
              </div>
            {/each}
            
            {#if !expandedLeagues.has(league.id) && allCharges.length === 0}
              <div class="p-8 text-center">
                <p class="text-gray-400 font-bold italic text-sm">No charges recorded.</p>
              </div>
            {:else if expandedLeagues.has(league.id) && recentActivity.length === 0}
              <div class="p-8 text-center">
                <p class="text-gray-400 font-bold italic text-sm">No activity recorded.</p>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {:else}
    <div class="py-20 text-center bg-white rounded-[32px] border-2 border-dashed border-gray-100">
      <p class="text-gray-400 font-bold italic">No financial history found.</p>
    </div>
  {/each}
</div>
