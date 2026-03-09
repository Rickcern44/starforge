<script lang="ts">
  import type { Payment, GameCharge, League } from '$lib/models';
  import { 
    ChevronLeft, 
    Wallet, 
    AlertCircle, 
    CheckCircle2, 
    ArrowUpRight, 
    ArrowDownLeft
  } from 'lucide-svelte';
  import { goto } from '$app/navigation';

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
  <title>Wallet | Bouncy</title>
</svelte:head>

<div class="max-w-3xl mx-auto flex flex-col min-h-screen pb-24 px-4">
  <div class="flex items-center space-x-4 py-6 text-left">
    <button onclick={() => history.back()} class="btn btn-ghost btn-circle border border-base-300">
      <ChevronLeft size={20} />
    </button>
    <h1 class="text-3xl font-black tracking-tight">Financial Ledger</h1>
  </div>

  <!-- 1. The Big Picture Card -->
  <div class="stats shadow-2xl shadow-neutral/20 bg-neutral text-neutral-content w-full rounded-[32px] mb-10 overflow-hidden border-none">
    <div class="stat p-8">
      <div class="stat-title text-sm font-bold uppercase tracking-widest opacity-60 text-neutral-content">Total Outstanding</div>
      <div class="stat-value text-5xl font-black">${totalBalance.toFixed(2)}</div>
      <div class="stat-figure text-white opacity-10">
        <Wallet size={64} />
      </div>
    </div>
  </div>

  <!-- 2. Priority: All Unpaid Charges (Action Items) -->
  <section class="mb-12">
    <div class="flex items-center justify-between mb-4 px-1">
      <h3 class="text-xs font-black uppercase tracking-widest opacity-40">Priority: Unpaid Charges</h3>
      <span class="badge badge-warning badge-sm font-black uppercase tracking-widest border-none px-3 py-3 text-warning-content">
        {allUnpaidCharges.length} {allUnpaidCharges.length === 1 ? 'Item' : 'Items'}
      </span>
    </div>
    
    <div class="card bg-base-100 shadow-sm border border-base-300 overflow-hidden rounded-[32px]">
      {#if allUnpaidCharges.length === 0}
        <div class="card-body p-12 text-center items-center">
          <div class="bg-success/10 w-16 h-16 rounded-full flex items-center justify-center mb-4 text-success">
            <CheckCircle2 size={32} />
          </div>
          <p class="font-bold">You're all caught up!</p>
          <p class="text-xs opacity-40 mt-1 uppercase font-bold tracking-widest">No outstanding charges</p>
        </div>
      {:else}
        <div class="divide-y divide-base-200">
          {#each allUnpaidCharges as charge}
            <div class="p-5 flex items-center space-x-4 hover:bg-base-200 transition-colors text-left">
              <div class="p-3 rounded-2xl bg-warning/10 text-warning">
                <AlertCircle size={24} />
              </div>
              <div class="flex-grow min-w-0">
                <h4 class="font-bold truncate leading-tight">{charge.title}</h4>
                <div class="flex items-center space-x-2 mt-1">
                  <span class="text-[9px] font-black text-primary uppercase tracking-widest">{charge.leagueName}</span>
                  <span class="opacity-20 text-[9px]">•</span>
                  <span class="text-[9px] font-bold opacity-40 uppercase tracking-widest">{charge.date}</span>
                </div>
              </div>
              <div class="text-right flex-shrink-0">
                <p class="text-[9px] font-black uppercase tracking-widest leading-none mb-1 opacity-40">Due</p>
                <p class="text-xl font-black text-error">${charge.unpaid.toFixed(2)}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>

  <!-- 3. League Breakdown (History) -->
  <section class="space-y-4">
    <h3 class="text-xs font-black uppercase tracking-widest px-1 opacity-40">League Breakdown</h3>
    <div class="space-y-3">
      {#each leaguesWithData as { league, allCharges, recentActivity, balance }}
        <div class="collapse collapse-arrow bg-base-100 border border-base-300 shadow-sm rounded-[24px]">
          <input type="checkbox" /> 
          <div class="collapse-title p-5 flex items-center gap-4 pr-12 text-left">
            <div class="h-10 w-10 rounded-xl bg-base-200 flex items-center justify-center text-base-content font-black text-xs uppercase border border-base-300">
              {league.name.charAt(0)}
            </div>
            <div class="flex-grow">
              <h4 class="font-bold leading-none mb-1.5">{league.name}</h4>
              <div class="flex items-center space-x-2">
                {#if balance > 0}
                  <span class="badge badge-ghost text-error badge-xs font-black uppercase tracking-widest py-2 border-none">Owes ${balance.toFixed(2)}</span>
                {:else}
                  <span class="badge badge-ghost text-success badge-xs font-black uppercase tracking-widest py-2 border-none">Settled</span>
                {/if}
              </div>
            </div>
          </div>
          <div class="collapse-content bg-base-200/20 px-0">
            <div class="divide-y divide-base-200">
              {#each recentActivity as item}
                {@const isCharge = item.type === 'charge'}
                <div class="px-6 py-4 flex items-center justify-between group bg-base-100 text-left">
                  <div class="flex items-center space-x-3">
                    <div class="p-1.5 rounded-lg {isCharge ? 'text-base-content/40 bg-base-200' : 'text-success bg-success/10'}">
                      {#if isCharge}
                        <ArrowUpRight size={16} />
                      {:else}
                        <ArrowDownLeft size={16} />
                      {/if}
                    </div>
                    <div>
                      <p class="text-[13px] font-bold leading-none mb-1">{item.title}</p>
                      <p class="text-[9px] opacity-40 font-bold uppercase tracking-widest">{item.date}</p>
                    </div>
                  </div>
                  <p class="text-sm font-black {isCharge ? 'opacity-40' : 'text-success'}">
                    {isCharge ? '-' : '+'}${item.amount.toFixed(2)}
                  </p>
                </div>
              {/each}
              
              {#if recentActivity.length === 0}
                <div class="p-8 text-center bg-base-100">
                  <p class="text-[10px] opacity-40 font-bold italic uppercase tracking-widest">No activity recorded.</p>
                </div>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>
  </section>
</div>
