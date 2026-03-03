<script lang="ts">
  import { page } from '$app/stores';
  import { inviteUser } from '$lib/services/league';
  import { addPayment, addAllocation } from '$lib/services/payment';
  import type { League, Payment, Game, GameCharge, Invitation } from '$lib/models';
  import { authService } from '$lib/services/auth.svelte';
  import { toastService } from '$lib/services/toast.svelte';

  let { data } = $props();

  let league = $state<League>(data.league);
  let payments = $state<Payment[]>(data.payments || []);
  let invitations = $state<Invitation[]>(data.invitations || []);
  let financialSummary = $derived(data.financialSummary);
  let activeTab = $state<'overview' | 'members' | 'payments'>('overview');

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
      // Refresh payments (in a real app we'd fetch or update state)
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
      <h1 class="text-3xl font-black text-gray-900 tracking-tight">League Admin</h1>
      <p class="text-sm font-bold text-indigo-600 uppercase tracking-wider">{league.name}</p>
    </div>
    <button onclick={() => history.back()} class="text-sm font-bold text-gray-500 hover:text-indigo-600 transition-colors">
      Back to League
    </button>
  </header>

  <!-- Tabs -->
  <div class="flex border-b border-gray-200">
    {#each ['overview', 'members', 'payments'] as tab}
      <button
        class="px-6 py-3 text-sm font-bold uppercase tracking-widest transition-all duration-200 border-b-2 -mb-[2px]"
        class:border-indigo-600={activeTab === tab}
        class:text-indigo-600={activeTab === tab}
        class:border-transparent={activeTab !== tab}
        class:text-gray-400={activeTab !== tab}
        onclick={() => activeTab = tab as any}
      >
        {tab}
      </button>
    {/each}
  </div>

  {#if activeTab === 'overview'}
    <div class="space-y-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-2">
          <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Members</p>
          <p class="text-3xl font-black text-gray-900">{league.members.length}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-2">
          <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Total Collected</p>
          <p class="text-3xl font-black text-green-600">${((financialSummary?.totalCollected || 0) / 100).toFixed(2)}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-2">
          <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Total Charges</p>
          <p class="text-3xl font-black text-red-600">${((financialSummary?.totalCharges || 0) / 100).toFixed(2)}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-2">
          <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Available Balance</p>
          <p class="text-3xl font-black text-indigo-600">${((financialSummary?.totalAvailable || 0) / 100).toFixed(2)}</p>
        </div>
      </div>

      <section class="space-y-3">
        <h3 class="text-xl font-black text-gray-900 px-1">Treasurer Overview</h3>
        <div class="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-4">
              <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                <span class="text-sm font-bold text-gray-500">Total Funds Collected</span>
                <span class="text-sm font-black text-gray-900">${((financialSummary?.totalCollected || 0) / 100).toFixed(2)}</span>
              </div>
              <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                <span class="text-sm font-bold text-gray-500">Allocated to Games</span>
                <span class="text-sm font-black text-green-600">-${((financialSummary?.totalAllocated || 0) / 100).toFixed(2)}</span>
              </div>
              <div class="flex justify-between items-center pt-2">
                <span class="text-sm font-black text-gray-900">Unallocated Funds (Credit)</span>
                <span class="text-sm font-black text-indigo-600">${((financialSummary?.totalAvailable || 0) / 100).toFixed(2)}</span>
              </div>
            </div>
            <div class="space-y-4">
              <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                <span class="text-sm font-bold text-gray-500">Total Game Revenue Due</span>
                <span class="text-sm font-black text-gray-900">${((financialSummary?.totalCharges || 0) / 100).toFixed(2)}</span>
              </div>
              <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                <span class="text-sm font-bold text-gray-500">Payments Applied</span>
                <span class="text-sm font-black text-green-600">-${((financialSummary?.totalAllocated || 0) / 100).toFixed(2)}</span>
              </div>
              <div class="flex justify-between items-center pt-2">
                <span class="text-sm font-black text-gray-900">Total Outstanding Debt</span>
                <span class="text-sm font-black text-red-600">${((financialSummary?.totalUnpaid || 0) / 100).toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  {/if}

  {#if activeTab === 'members'}
    <section class="space-y-6">
      <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-4">
        <h3 class="text-xl font-black text-gray-900">Invite New Member</h3>
        <div class="flex space-x-2">
          <input
            type="email"
            bind:value={inviteEmail}
            placeholder="player@example.com"
            class="flex-grow p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium"
          />
          <button
            onclick={handleInvite}
            disabled={isInviting || !inviteEmail}
            class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-3 px-6 rounded-xl shadow-lg shadow-indigo-100 transition-all disabled:opacity-50"
          >
            {isInviting ? 'Sending...' : 'Send Invite'}
          </button>
        </div>
      </div>

      {#if invitations.length > 0}
        <div class="space-y-3">
          <h3 class="text-xl font-black text-gray-900 px-1">Pending Invitations</h3>
          <div class="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-gray-50 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  <th class="px-6 py-4">Email</th>
                  <th class="px-6 py-4">Status</th>
                  <th class="px-6 py-4">Expires</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-50">
                {#each invitations as invite}
                  {@const isExpired = new Date(invite.expiresAt) < new Date()}
                  <tr class="group hover:bg-gray-50 transition-colors">
                    <td class="px-6 py-4 text-sm font-bold text-gray-900">{invite.email}</td>
                    <td class="px-6 py-4 text-sm font-medium">
                      {#if invite.usedAt}
                        <span class="px-2 py-1 bg-green-50 text-green-600 rounded-md text-[9px] font-black uppercase tracking-widest">Used</span>
                      {:else if isExpired}
                        <span class="px-2 py-1 bg-red-50 text-red-600 rounded-md text-[9px] font-black uppercase tracking-widest">Expired</span>
                      {:else}
                        <span class="px-2 py-1 bg-blue-50 text-blue-600 rounded-md text-[9px] font-black uppercase tracking-widest">Active</span>
                      {/if}
                    </td>
                    <td class="px-6 py-4 text-sm font-medium text-gray-400">{new Date(invite.expiresAt).toLocaleDateString()}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}

      <div class="space-y-3">
        <h3 class="text-xl font-black text-gray-900 px-1">League Members</h3>
        <div class="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gray-50 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                <th class="px-6 py-4">Name</th>
                <th class="px-6 py-4">Role</th>
                <th class="px-6 py-4">Balance</th>
                <th class="px-6 py-4">Joined</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              {#each league.members as member}
                {@const balance = getMemberBalance(member.playerId)}
                <tr class="group hover:bg-gray-50 transition-colors">
                  <td class="px-6 py-4">
                    <p class="text-sm font-bold text-gray-900">{member.playerName}</p>
                    <p class="text-[10px] text-gray-400 font-mono tracking-tighter">{member.playerId.substring(0, 8)}...</p>
                  </td>
                  <td class="px-6 py-4 text-sm font-medium text-gray-500">
                    <span class="px-2 py-1 rounded-md text-[10px] font-black uppercase tracking-widest" 
                      class:bg-indigo-100={member.role === 'admin' || member.role === 'owner'}
                      class:text-indigo-600={member.role === 'admin' || member.role === 'owner'}
                      class:bg-gray-100={member.role !== 'admin' && member.role !== 'owner'}
                      class:text-gray-500={member.role !== 'admin' && member.role !== 'owner'}
                    >
                      {member.role}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-sm font-black {balance > 0 ? 'text-red-600' : 'text-green-600'}">
                    ${balance.toFixed(2)}
                  </td>
                  <td class="px-6 py-4 text-sm font-medium text-gray-400">{new Date(member.joinedAt).toLocaleDateString()}</td>
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
        <h3 class="text-xl font-black text-gray-900">Add New Payment</h3>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Amount ($)</label>
              <input type="number" step="0.01" bind:value={newPaymentAmount} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium" />
            </div>
            <div class="space-y-1">
              <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Method</label>
              <select bind:value={newPaymentMethod} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium appearance-none">
                <option value="venmo">Venmo</option>
                <option value="cash">Cash</option>
              </select>
            </div>
          </div>
          <div class="flex items-center justify-between mb-2">
            <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Payer</label>
            <button 
              onclick={() => isExternalPayer = !isExternalPayer}
              class="text-[9px] font-black text-indigo-600 uppercase tracking-widest hover:text-indigo-800 transition-colors"
            >
              {isExternalPayer ? 'Switch to Member List' : 'Switch to Guest/External'}
            </button>
          </div>

          {#if isExternalPayer}
            <div class="space-y-1">
              <input type="text" bind:value={newPaymentExternalName} placeholder="As it appears on Venmo/Sheet" class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium" />
            </div>
          {:else}
            <div class="space-y-1">
              <select bind:value={selectedMemberId} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium">
                <option value="">Select a member...</option>
                {#each league.members as member}
                  <option value={member.playerId}>{member.playerName}</option>
                {/each}
              </select>
            </div>
          {/if}

          <div class="space-y-1">
            <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Reference / Note</label>
            <input type="text" bind:value={newPaymentReference} placeholder="Optional" class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium" />
          </div>
          <button
            onclick={handleAddPayment}
            disabled={isAddingPayment || newPaymentAmount <= 0 || (!isExternalPayer && !selectedMemberId) || (isExternalPayer && !newPaymentExternalName)}
            class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-black py-4 rounded-xl shadow-lg transition-all disabled:opacity-50"
          >
            {isAddingPayment ? 'Processing...' : 'Record Payment'}
          </button>
        </div>
      </section>

      <!-- Allocation -->
      <section class="space-y-4">
        <h3 class="text-xl font-black text-gray-900">Allocate Funds</h3>
        <div class="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm space-y-4">
          <div class="space-y-1">
            <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Select Payment</label>
            <select bind:value={selectedPayment} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium">
              <option value={null}>Choose unallocated payment...</option>
              {#each unallocatedPayments as p}
                <option value={p}>{p.externalName} (${getUnallocatedAmount(p)})</option>
              {/each}
            </select>
          </div>

          <div class="space-y-1">
            <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Select Unpaid Charge</label>
            <select bind:value={selectedCharge} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium">
              <option value={null}>Choose unpaid charge...</option>
              {#each unpaidCharges as c}
                <option value={c}>{c.externalName || 'Guest'} (${getUnpaidAmount(c)}) - Game {c.gameId.substring(0, 5)}</option>
              {/each}
            </select>
          </div>

          <div class="space-y-1">
            <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest px-1">Allocation Amount ($)</label>
            <input type="number" step="0.01" bind:value={allocationAmount} class="w-full p-3 bg-gray-50 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 text-sm font-medium" />
          </div>

          <button
            onclick={handleAllocate}
            disabled={isAllocating || !selectedPayment || !selectedCharge || allocationAmount <= 0}
            class="w-full bg-green-600 hover:bg-green-700 text-white font-black py-4 rounded-xl shadow-lg transition-all disabled:opacity-50"
          >
            {isAllocating ? 'Allocating...' : 'Confirm Allocation'}
          </button>
        </div>
      </section>
    </div>

    <!-- Recent Payments Table -->
    <section class="space-y-3 pt-4">
      <h3 class="text-xl font-black text-gray-900 px-1">Recent Payments</h3>
      <div class="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 text-[10px] font-black text-gray-400 uppercase tracking-widest">
              <th class="px-6 py-4">Payer</th>
              <th class="px-6 py-4">Amount</th>
              <th class="px-6 py-4">Status</th>
              <th class="px-6 py-4">Date</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            {#each payments as payment}
              {@const unallocated = getUnallocatedAmount(payment)}
              <tr class="group hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4">
                  <p class="text-sm font-bold text-gray-900">{payment.externalName}</p>
                  {#if payment.reference}
                    <p class="text-[10px] text-gray-400 italic">"{payment.reference}"</p>
                  {/if}
                </td>
                <td class="px-6 py-4 text-sm font-black text-gray-900">${(payment.amountCents / 100).toFixed(2)}</td>
                <td class="px-6 py-4">
                  {#if unallocated === 0}
                    <span class="px-2 py-1 bg-green-50 text-green-600 rounded-md text-[9px] font-black uppercase tracking-widest">Fully Allocated</span>
                  {:else}
                    <span class="px-2 py-1 bg-orange-50 text-orange-600 rounded-md text-[9px] font-black uppercase tracking-widest">${unallocated.toFixed(2)} Remaining</span>
                  {/if}
                </td>
                <td class="px-6 py-4 text-sm font-medium text-gray-400">{new Date(payment.receivedAt).toLocaleDateString()}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </section>
  {/if}
</div>
