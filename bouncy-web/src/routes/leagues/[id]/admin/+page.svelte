<script lang="ts">
  import { inviteUser } from '$lib/services/league';
  import { addPayment, addAllocation } from '$lib/services/payment';
  import type { League, Payment, Game, GameCharge, Invitation } from '$lib/models';
  import { authService } from '$lib/services/auth.svelte';
  import { featureFlagService } from '$lib/services/feature-flag.svelte';
  import { toastService } from '$lib/services/toast.svelte';
  import { 
    ChevronLeft, 
    Users, 
    DollarSign, 
    TrendingUp, 
    TrendingDown, 
    Send, 
    Plus, 
    CheckCircle2, 
    CreditCard
  } from 'lucide-svelte';

  let { data } = $props();

  let league = $state<League>(data.league);
  let payments = $state<Payment[]>(data.payments || []);
  let invitations = $state<Invitation[]>(data.invitations || []);
  let financialSummary = $derived(data.financialSummary);
  let activeTab = $state<'overview' | 'members' | 'payments'>('overview');

  let availableTabs = $derived.by(() => {
    const tabs = ['overview', 'members'];
    if (featureFlagService.isEnabled('payments')) {
      tabs.push('payments');
    }
    return tabs;
  });

  // Member management
  let inviteEmail = $state('');
  let isInviting = $state(false);

  // Payment management
  let newPaymentAmount = $state<number>(0);
  let newPaymentMethod = $state<'venmo' | 'cash'>('venmo');
  let selectedMemberId = $state<string>('');
  let newPaymentExternalName = $state('');
  let newPaymentReference = $state('');
  let isAddingPayment = $state(false);
  let isExternalPayer = $state(false);

  // Allocation management
  let selectedPayment = $state<Payment | null>(null);
  let selectedCharge = $state<GameCharge | null>(null);
  let allocationAmount = $state<number>(0);
  let isAllocating = $state(false);

  // Derived data
  let allCharges = $derived.by(() => {
    const charges: GameCharge[] = [];
    league.games.forEach(game => {
      if (game.charges) {
        charges.push(...game.charges);
      }
    });
    return charges;
  });

  let unallocatedPayments = $derived(payments.filter(p => {
    const allocated = p.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return allocated < p.amountCents;
  }));

  let unpaidCharges = $derived(allCharges.filter(c => {
    const paid = c.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return paid < c.amountCents;
  }));

  function getMemberBalance(playerId: string) {
    let balanceCents = 0;
    
    // Sum unpaid portions of all charges for this user in this league
    allCharges.forEach(charge => {
      if (charge.userId === playerId) {
        const paid = charge.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
        balanceCents += (charge.amountCents - paid);
      }
    });

    return balanceCents / 100;
  }

  async function handleInvite() {
    if (!inviteEmail) return;
    isInviting = true;
    const success = await inviteUser(inviteEmail, league.id);
    if (success) {
      toastService.success(`Invitation sent to ${inviteEmail}`);
      inviteEmail = '';
    } else {
      toastService.error('Failed to send invitation.');
    }
    isInviting = false;
  }

  async function handleAddPayment() {
    if (newPaymentAmount <= 0) return;
    
    let payerName = '';
    let userId: string | null = null;

    if (isExternalPayer) {
      if (!newPaymentExternalName) return;
      payerName = newPaymentExternalName;
    } else {
      if (!selectedMemberId) return;
      const member = league.members.find(m => m.playerId === selectedMemberId);
      if (!member) return;
      payerName = member.playerName;
      userId = member.playerId;
    }

    isAddingPayment = true;
    const paymentData = {
      amountCents: Math.round(newPaymentAmount * 100),
      method: newPaymentMethod,
      externalName: payerName,
      userId: userId,
      reference: newPaymentReference,
      recordedBy: authService.user?.id
    };
    const success = await addPayment(league.id, paymentData);
    if (success) {
      toastService.success('Payment added successfully.');
      // Refresh payments
      location.reload(); 
    } else {
      toastService.error('Failed to add payment.');
    }
    isAddingPayment = false;
  }

  async function handleAllocate() {
    if (!selectedPayment || !selectedCharge || allocationAmount <= 0) return;
    isAllocating = true;
    
    const success = await addAllocation(selectedPayment.id, {
      gameChargeId: selectedCharge.id,
      amountInCents: Math.round(allocationAmount * 100)
    });

    if (success) {
      toastService.success('Allocation successful.');
      location.reload();
    } else {
      toastService.error('Failed to allocate payment.');
    }
    isAllocating = false;
  }

  function getUnallocatedAmount(payment: Payment) {
    const allocated = payment.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return (payment.amountCents - allocated) / 100;
  }

  function getUnpaidAmount(charge: GameCharge) {
    const paid = charge.allocations?.reduce((sum, a) => sum + a.amountInCents, 0) || 0;
    return (charge.amountCents - paid) / 100;
  }
</script>

<svelte:head>
  <title>Admin: {league.name} | Bouncy</title>
</svelte:head>

<div class="max-w-4xl mx-auto py-8 px-4 space-y-8">
  <header class="flex justify-between items-center">
    <div>
      <h1 class="text-3xl font-black tracking-tight leading-tight">League Admin</h1>
      <p class="text-sm font-bold text-primary mt-1 uppercase tracking-wider">{league.name}</p>
    </div>
    <button onclick={() => history.back()} class="btn btn-ghost btn-sm gap-2">
      <ChevronLeft size={16} />
      Back
    </button>
  </header>

  <!-- Tabs -->
  <div class="flex border-b border-base-300 overflow-x-auto no-scrollbar whitespace-nowrap -mx-4 px-4 sm:mx-0 sm:px-0">
    {#each availableTabs as tab}
      <button
        class="px-6 py-3 text-xs sm:text-sm font-black uppercase tracking-widest transition-all duration-200 border-b-2 -mb-[2px] flex-shrink-0"
        class:border-neutral={activeTab === tab}
        class:text-neutral={activeTab === tab}
        class:border-transparent={activeTab !== tab}
        class:opacity-40={activeTab !== tab}
        onclick={() => activeTab = tab as any}
      >
        {tab}
      </button>
    {/each}
  </div>

  {#if activeTab === 'overview'}
    <div class="space-y-8">
      <div class="stats stats-vertical lg:stats-horizontal shadow bg-base-100 w-full rounded-[32px] border border-base-300">
        <div class="stat">
          <div class="stat-figure opacity-20"><Users size={24} /></div>
          <div class="stat-title text-[10px] font-black uppercase tracking-widest">Members</div>
          <div class="stat-value text-3xl font-black">{league.members.length}</div>
        </div>
        {#if featureFlagService.isEnabled('payments')}
        <div class="stat">
          <div class="stat-figure text-success opacity-20"><TrendingUp size={24} /></div>
          <div class="stat-title text-[10px] font-black uppercase tracking-widest">Collected</div>
          <div class="stat-value text-3xl font-black text-success">${((financialSummary?.totalCollected || 0) / 100).toFixed(2)}</div>
        </div>
        <div class="stat">
          <div class="stat-figure text-error opacity-20"><TrendingDown size={24} /></div>
          <div class="stat-title text-[10px] font-black uppercase tracking-widest">Charges</div>
          <div class="stat-value text-3xl font-black text-error">${((financialSummary?.totalCharges || 0) / 100).toFixed(2)}</div>
        </div>
        <div class="stat">
          <div class="stat-figure text-primary opacity-20"><DollarSign size={24} /></div>
          <div class="stat-title text-[10px] font-black uppercase tracking-widest">Balance</div>
          <div class="stat-value text-3xl font-black text-primary">${((financialSummary?.totalAvailable || 0) / 100).toFixed(2)}</div>
        </div>
        {/if}
      </div>

      {#if featureFlagService.isEnabled('payments')}
      <section class="space-y-3">
        <h3 class="text-xs font-black opacity-40 uppercase tracking-widest px-1">Treasurer Overview</h3>
        <div class="card bg-base-100 border border-base-300 shadow-sm overflow-hidden rounded-[32px]">
          <div class="card-body p-8">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-12">
              <div class="space-y-4">
                <div class="flex justify-between items-center border-b border-base-200 pb-2">
                  <span class="text-sm font-bold opacity-40">Total Funds Collected</span>
                  <span class="text-sm font-black">${((financialSummary?.totalCollected || 0) / 100).toFixed(2)}</span>
                </div>
                <div class="flex justify-between items-center border-b border-base-200 pb-2">
                  <span class="text-sm font-bold opacity-40">Allocated to Games</span>
                  <span class="text-sm font-black text-success">-${((financialSummary?.totalAllocated || 0) / 100).toFixed(2)}</span>
                </div>
                <div class="flex justify-between items-center pt-2">
                  <span class="text-sm font-black">Unallocated Funds</span>
                  <span class="text-sm font-black text-primary">${((financialSummary?.totalAvailable || 0) / 100).toFixed(2)}</span>
                </div>
              </div>
              <div class="space-y-4">
                <div class="flex justify-between items-center border-b border-base-200 pb-2">
                  <span class="text-sm font-bold opacity-40">Total Game Revenue Due</span>
                  <span class="text-sm font-black">${((financialSummary?.totalCharges || 0) / 100).toFixed(2)}</span>
                </div>
                <div class="flex justify-between items-center border-b border-base-200 pb-2">
                  <span class="text-sm font-bold opacity-40">Payments Applied</span>
                  <span class="text-sm font-black text-success">-${((financialSummary?.totalAllocated || 0) / 100).toFixed(2)}</span>
                </div>
                <div class="flex justify-between items-center pt-2">
                  <span class="text-sm font-black">Outstanding Debt</span>
                  <span class="text-sm font-black text-error">${((financialSummary?.totalUnpaid || 0) / 100).toFixed(2)}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
      {/if}
    </div>
  {/if}

  {#if activeTab === 'members'}
    <section class="space-y-6">
      {#if featureFlagService.isEnabled('admin_invites')}
      <div class="card bg-base-100 p-8 rounded-[32px] border border-base-300 shadow-sm space-y-4">
        <h3 class="text-xl font-black">Invite New Member</h3>
        <div class="join w-full">
          <input
            type="email"
            bind:value={inviteEmail}
            placeholder="player@example.com"
            class="input input-bordered join-item flex-grow focus:outline-none"
          />
          <button
            onclick={handleInvite}
            disabled={isInviting || !inviteEmail}
            class="btn btn-neutral join-item font-black uppercase tracking-widest px-8"
          >
            {#if isInviting}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              <Send size={16} />
            {/if}
            Invite
          </button>
        </div>
      </div>
      {/if}

      {#if invitations.length > 0}
        <div class="space-y-3">
          <h3 class="text-xs font-black uppercase tracking-widest opacity-40 px-1">Pending Invitations</h3>
          <div class="card bg-base-100 border border-base-300 shadow-sm overflow-x-auto no-scrollbar rounded-[32px]">
            <table class="table w-full text-left min-w-[600px]">
              <thead>
                <tr class="text-[10px] font-black uppercase tracking-widest opacity-40">
                  <th class="px-8 py-4">Email</th>
                  <th class="px-8 py-4">Status</th>
                  <th class="px-8 py-4">Expires</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-base-200">
                {#each invitations as invite}
                  {@const isExpired = new Date(invite.expiresAt) < new Date()}
                  <tr class="hover">
                    <td class="px-8 py-5 text-sm font-bold">{invite.email}</td>
                    <td class="px-8 py-5">
                      {#if invite.usedAt}
                        <span class="badge badge-success badge-sm font-black uppercase tracking-tighter text-[9px] text-white border-none">Used</span>
                      {:else if isExpired}
                        <span class="badge badge-error badge-sm font-black uppercase tracking-tighter text-[9px] border-none">Expired</span>
                      {:else}
                        <span class="badge badge-ghost badge-sm font-black uppercase tracking-tighter text-[9px] border border-base-300">Active</span>
                      {/if}
                    </td>
                    <td class="px-8 py-5 text-sm font-medium opacity-40">{new Date(invite.expiresAt).toLocaleDateString()}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}

      <div class="space-y-3">
        <h3 class="text-xs font-black uppercase tracking-widest opacity-40 px-1">League Members</h3>
        <div class="card bg-base-100 border border-base-300 shadow-sm overflow-x-auto no-scrollbar rounded-[32px]">
          <table class="table w-full text-left min-w-[600px]">
            <thead>
              <tr class="text-[10px] font-black uppercase tracking-widest opacity-40">
                <th class="px-8 py-4">Name</th>
                <th class="px-8 py-4">Role</th>
                <th class="px-8 py-4 text-right">Balance</th>
                <th class="px-8 py-4">Joined</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-base-200">
              {#each league.members as member}
                {@const balance = getMemberBalance(member.playerId)}
                <tr class="hover">
                  <td class="px-8 py-5">
                    <p class="text-sm font-bold">{member.playerName}</p>
                    <p class="text-[10px] opacity-40 font-mono tracking-tighter">{member.playerId.substring(0, 8)}...</p>
                  </td>
                  <td class="px-8 py-5">
                    <span class="badge badge-ghost badge-sm font-black uppercase tracking-tighter text-[9px] border border-base-300" 
                      class:badge-neutral={member.role === 'admin' || member.role === 'owner'}
                    >
                      {member.role}
                    </span>
                  </td>
                  <td class="px-8 py-5 text-sm font-black text-right {balance > 0 ? 'text-error' : 'text-success'}">
                    ${balance.toFixed(2)}
                  </td>
                  <td class="px-8 py-5 text-sm font-medium opacity-40">{new Date(member.joinedAt).toLocaleDateString()}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>
  {/if}

  {#if activeTab === 'payments'}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Add Payment -->
      <section class="space-y-4">
        <h3 class="text-xl font-black">Add New Payment</h3>
        <div class="card bg-base-100 p-8 rounded-[32px] border border-base-300 shadow-sm space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control w-full">
              <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Amount ($)</span></label>
              <input type="number" step="0.01" bind:value={newPaymentAmount} class="input input-bordered focus:outline-none" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Method</span></label>
              <select bind:value={newPaymentMethod} class="select select-bordered focus:outline-none">
                <option value="venmo">Venmo</option>
                <option value="cash">Cash</option>
              </select>
            </div>
          </div>
          
          <div class="form-control">
            <label class="label cursor-pointer justify-between">
              <span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Payer Type</span>
              <button 
                onclick={() => isExternalPayer = !isExternalPayer}
                class="btn btn-link btn-xs no-underline hover:no-underline font-black uppercase tracking-widest text-[9px] text-primary"
              >
                {isExternalPayer ? 'Switch to Member List' : 'Switch to Guest/External'}
              </button>
            </label>

            {#if isExternalPayer}
              <input type="text" bind:value={newPaymentExternalName} placeholder="As it appears on Venmo" class="input input-bordered focus:outline-none w-full" />
            {:else}
              <select bind:value={selectedMemberId} class="select select-bordered focus:outline-none w-full">
                <option value="">Select a member...</option>
                {#each league.members as member}
                  <option value={member.playerId}>{member.playerName}</option>
                {/each}
              </select>
            {/if}
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Reference / Note</span></label>
            <input type="text" bind:value={newPaymentReference} placeholder="Optional" class="input input-bordered focus:outline-none w-full" />
          </div>
          
          <button
            onclick={handleAddPayment}
            disabled={isAddingPayment || newPaymentAmount <= 0 || (!isExternalPayer && !selectedMemberId) || (isExternalPayer && !newPaymentExternalName)}
            class="btn btn-neutral btn-block font-black uppercase tracking-widest py-4 h-auto shadow-lg"
          >
            {#if isAddingPayment}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              <Plus size={18} />
            {/if}
            Record Payment
          </button>
        </div>
      </section>

      <!-- Allocation -->
      <section class="space-y-4">
        <h3 class="text-xl font-black">Allocate Funds</h3>
        <div class="card bg-base-100 p-8 rounded-[32px] border border-base-300 shadow-sm space-y-4">
          <div class="form-control w-full">
            <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Select Payment</span></label>
            <select bind:value={selectedPayment} class="select select-bordered focus:outline-none w-full">
              <option value={null}>Choose unallocated payment...</option>
              {#each unallocatedPayments as p}
                <option value={p}>{p.externalName} (${getUnallocatedAmount(p)})</option>
              {/each}
            </select>
          </div>

          <div class="form-control w-full">
            <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Select Unpaid Charge</span></label>
            <select bind:value={selectedCharge} class="select select-bordered focus:outline-none w-full">
              <option value={null}>Choose unpaid charge...</option>
              {#each unpaidCharges as c}
                <option value={c}>{c.externalName || 'Guest'} (${getUnpaidAmount(c)}) - Game {c.gameId.substring(0, 5)}</option>
              {/each}
            </select>
          </div>

          <div class="form-control w-full">
            <label class="label"><span class="label-text text-[10px] font-black uppercase tracking-widest opacity-40">Allocation Amount ($)</span></label>
            <input type="number" step="0.01" bind:value={allocationAmount} class="input input-bordered focus:outline-none" />
          </div>

          <button
            onclick={handleAllocate}
            disabled={isAllocating || !selectedPayment || !selectedCharge || allocationAmount <= 0}
            class="btn btn-success btn-block font-black uppercase tracking-widest py-4 h-auto shadow-lg text-white"
          >
            {#if isAllocating}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              <CheckCircle2 size={18} />
            {/if}
            Confirm Allocation
          </button>
        </div>
      </section>
    </div>

    <!-- Recent Payments Table -->
    <section class="space-y-3 pt-4">
      <h3 class="text-xs font-black uppercase tracking-widest opacity-40 px-1">Recent League Activity</h3>
      <div class="card bg-base-100 border border-base-300 shadow-sm overflow-x-auto no-scrollbar rounded-[32px]">
        <table class="table w-full text-left min-w-[600px]">
          <thead>
            <tr class="text-[10px] font-black uppercase tracking-widest opacity-40">
              <th class="px-8 py-4">Payer</th>
              <th class="px-8 py-4">Amount</th>
              <th class="px-8 py-4">Status</th>
              <th class="px-8 py-4">Date</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-base-200">
            {#each payments as payment}
              {@const unallocated = getUnallocatedAmount(payment)}
              <tr class="hover">
                <td class="px-8 py-5">
                  <div class="flex items-center gap-3">
                    <div class="p-2 rounded-xl bg-base-200 opacity-40"><CreditCard size={16} /></div>
                    <div>
                      <p class="text-sm font-bold">{payment.externalName}</p>
                      {#if payment.reference}
                        <p class="text-[10px] opacity-40 italic">"{payment.reference}"</p>
                      {/if}
                    </div>
                  </div>
                </td>
                <td class="px-8 py-5 text-sm font-black">${(payment.amountCents / 100).toFixed(2)}</td>
                <td class="px-8 py-5">
                  {#if unallocated === 0}
                    <span class="badge badge-success badge-sm font-black uppercase tracking-tighter text-[9px] text-white border-none">Fully Allocated</span>
                  {:else}
                    <span class="badge badge-warning badge-sm font-black uppercase tracking-tighter text-[9px] border-none">${unallocated.toFixed(2)} Remaining</span>
                  {/if}
                </td>
                <td class="px-8 py-5 text-sm font-medium opacity-40">{new Date(payment.receivedAt).toLocaleDateString()}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </section>
  {/if}
</div>
