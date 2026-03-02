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
  
      if (data.charges) {
        data.charges.forEach((c: GameCharge) => {
          const leagueId = c.game?.leagueId || 'unknown';
          if (!leagueMap.has(leagueId)) return;
  
          const unpaid = getUnpaidAmount(c);
          const chargeItem = {
            id: c.id,
            leagueName: leagueMap.get(leagueId)?.league.name,
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
            amount: p.amountCents / 100,
            status: 'completed'
          });
        });
      }
  
      leagueMap.forEach(entry => {
        entry.recentActivity.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
        entry.unpaidCharges.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
        entry.allCharges.sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime());
      });
  
      return Array.from(leagueMap.values())
        .filter(entry => entry.recentActivity.length > 0 || entry.league.isActive)
        .sort((a, b) => a.league.name.localeCompare(b.league.name));
    });
  
    let allUnpaidCharges = $derived(
      leaguesWithData
        .flatMap(l => l.unpaidCharges)
        .sort((a, b) => b.rawDate.getTime() - a.rawDate.getTime())
    );
  
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

  <!-- 1. The Big Picture Card -->
  <div class="bg-indigo-600 rounded-[32px] p-8 mb-10 text-white shadow-xl shadow-indigo-100 flex justify-between items-center relative overflow-hidden">
    <div class="relative z-10">
      <p class="text-sm font-bold text-indigo-100 uppercase tracking-widest opacity-80 leading-none mb-2">Total Outstanding</p>
      <p class="text-5xl font-black">${totalBalance.toFixed(2)}</p>
    </div>
    <div class="bg-white/10 p-4 rounded-2xl backdrop-blur-md relative z-10 border border-white/10">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
    </div>
    <div class="absolute -right-10 -bottom-10 w-48 h-48 bg-indigo-500 rounded-full blur-3xl opacity-50"></div>
  </div>

  <!-- 2. Priority: All Unpaid Charges (Action Items) -->
  <section class="mb-12">
    <div class="flex items-center justify-between mb-4 px-1">
      <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest">Priority: Unpaid Charges</h3>
      <span class="bg-orange-100 text-orange-600 px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-widest">
        {allUnpaidCharges.length} {allUnpaidCharges.length === 1 ? 'Item' : 'Items'}
      </span>
    </div>
    
    <div class="bg-white rounded-[32px] shadow-sm border border-gray-100 overflow-hidden">
      {#if allUnpaidCharges.length === 0}
        <div class="p-12 text-center">
          <div class="bg-green-50 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 text-green-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <p class="text-gray-800 font-bold">You're all caught up!</p>
          <p class="text-gray-400 text-xs mt-1">No outstanding charges from any of your leagues.</p>
        </div>
      {:else}
        <div class="divide-y divide-gray-50">
          {#each allUnpaidCharges as charge}
            <div class="p-5 flex items-center space-x-4 hover:bg-gray-50 transition-colors">
              <div class="p-3 rounded-2xl bg-orange-50 text-orange-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="flex-grow min-w-0">
                <h4 class="font-bold text-gray-900 truncate leading-tight">{charge.title}</h4>
                <div class="flex items-center space-x-2 mt-1">
                  <span class="text-[9px] font-black text-indigo-400 uppercase tracking-widest">{charge.leagueName}</span>
                  <span class="text-gray-300 text-[9px]">•</span>
                  <span class="text-[9px] font-bold text-gray-400 uppercase tracking-widest">{charge.date}</span>
                </div>
              </div>
              <div class="text-right flex-shrink-0">
                <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Due</p>
                <p class="text-xl font-black text-red-600">${charge.unpaid.toFixed(2)}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>

  <!-- 3. League Breakdown (History) -->
  <section class="space-y-4">
    <h3 class="text-sm font-black text-gray-400 uppercase tracking-widest px-1">League Breakdown</h3>
    <div class="space-y-3">
      {#each leaguesWithData as { league, allCharges, recentActivity, balance }}
        {@const isExpanded = expandedLeagues.has(league.id)}
        <div class="bg-white rounded-[24px] border border-gray-100 shadow-sm overflow-hidden transition-all duration-300">
          <!-- Accordion Header -->
          <button 
            onclick={() => toggleExpand(league.id)}
            class="w-full p-5 flex items-center justify-between hover:bg-gray-50 transition-colors"
          >
            <div class="flex items-center space-x-4">
              <div class="h-10 w-10 rounded-xl bg-gray-50 flex items-center justify-center text-indigo-600 font-black text-xs uppercase border border-gray-100">
                {league.name.charAt(0)}
              </div>
              <div class="text-left">
                <h4 class="font-bold text-gray-900 leading-none mb-1.5">{league.name}</h4>
                <div class="flex items-center space-x-2">
                  {#if balance > 0}
                    <span class="text-[9px] font-black text-red-500 uppercase tracking-widest">Owes ${balance.toFixed(2)}</span>
                  {:else}
                    <span class="text-[9px] font-black text-green-500 uppercase tracking-widest">Settled</span>
                  {/if}
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-3">
              <span class="text-[10px] font-black text-gray-300 uppercase tracking-widest">
                {isExpanded ? 'Collapse' : 'History'}
              </span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 transition-transform duration-300 {isExpanded ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </button>

          <!-- Accordion Content -->
          {#if isExpanded}
            <div class="bg-gray-50/50 border-t border-gray-50">
              <div class="divide-y divide-gray-100">
                {#each recentActivity as item}
                  {@const isCharge = item.type === 'charge'}
                  <div class="p-4 flex items-center justify-between group">
                    <div class="flex items-center space-x-3">
                      <div class="p-1.5 rounded-lg {isCharge ? 'text-gray-400 bg-gray-100' : 'text-green-600 bg-green-50'}">
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
                      <div>
                        <p class="text-[13px] font-bold text-gray-700 leading-none mb-1">{item.title}</p>
                        <p class="text-[9px] text-gray-400 font-bold uppercase tracking-widest">{item.date}</p>
                      </div>
                    </div>
                    <p class="text-sm font-black {isCharge ? 'text-gray-400' : 'text-green-600'}">
                      {isCharge ? '-' : '+'}${item.amount.toFixed(2)}
                    </p>
                  </div>
                {/each}
                
                {#if recentActivity.length === 0}
                  <div class="p-8 text-center">
                    <p class="text-[10px] text-gray-400 font-bold italic uppercase tracking-widest">No activity recorded.</p>
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </section>
</div>
