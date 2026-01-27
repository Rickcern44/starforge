<script lang="ts">
    import Card from '$lib/components/Card.svelte';

    export let item: UpcomingGame;

    const formatTime = (iso: string) =>
        new Intl.DateTimeFormat('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: '2-digit'
        }).format(new Date(iso));

    const formatRelative = (iso: string) => {
        const now = new Date();
        const date = new Date(iso);
        const diffMs = date.getTime() - now.getTime();
        const diffDays = Math.round(diffMs / (1000 * 60 * 60 * 24));

        if (diffDays === 0) return 'Today';
        if (diffDays === 1) return 'Tomorrow';
        if (diffDays === -1) return 'Yesterday';
        return new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric'
        }).format(date);
    };

    const formatCurrency = (cents: number) =>
        cents === 0
            ? 'Free'
            : new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: 'USD'
            }).format(cents / 100);

    const attendanceLabel = (g: UpcomingGame) => {
        if (g.game.isCanceled) return 'Canceled';
        if (!g.attendance) return 'Not joined';
        if (g.attendance.checkedIn) return 'Checked in';
        return 'Going';
    };

    const paymentLabel = (g: UpcomingGame) => {
        if (g.game.costInCents === 0) return 'No payment needed';
        if (!g.payment) return 'Payment needed';
        if (g.payment.status === 'completed') return 'Paid';
        if (g.payment.status === 'pending') return 'Payment pending';
        if (g.payment.status === 'failed') return 'Payment failed';
        return 'Payment updated';
    };

    const primaryActionLabel = (g: UpcomingGame) => {
        if (g.game.isCanceled) return 'View details';
        if (!g.attendance) return 'Join game';
        if (g.game.costInCents > 0 && (!g.payment || g.payment.status !== 'completed')) {
            return 'Pay now';
        }
        if (g.attendance.checkedIn) return 'View receipt';
        return 'View details';
    };

    const attendancePillClass = (g: UpcomingGame) => {
        const label = attendanceLabel(g);
        if (label === 'Canceled') return 'bg-rose-500/10 text-rose-300 ring-1 ring-rose-500/40';
        if (label === 'Checked in') return 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40';
        if (label === 'Going') return 'bg-sky-500/10 text-sky-300 ring-1 ring-sky-500/40';
        return 'bg-slate-800 text-slate-300 ring-1 ring-slate-700/60';
    };

    const paymentPillClass = (g: UpcomingGame) => {
        const label = paymentLabel(g);
        if (label === 'Paid') return 'bg-emerald-500/10 text-emerald-300';
        if (label === 'Payment pending') return 'bg-amber-500/10 text-amber-300';
        if (label === 'Payment failed') return 'bg-rose-500/10 text-rose-300';
        return 'bg-slate-900 text-slate-400';
    };
</script>

<Card>
    <div slot="header" class="w-full flex items-start justify-between gap-3">
        <div class="space-y-1">
            <div class="text-xs uppercase tracking-wide text-slate-400">
                {item.league.name}
                {#if !item.league.isActive}
          <span class="ml-1 rounded-full bg-amber-500/10 px-1.5 py-0.5 text-[10px] text-amber-300">
            Inactive
          </span>
                {/if}
            </div>
            <div class="text-sm font-medium text-slate-100">
                {formatRelative(item.game.startTime)}
            </div>
            <div class="text-xs text-slate-400">
                {formatTime(item.game.startTime)}
            </div>
        </div>

        <span class={`inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-medium ${attendancePillClass(item)}`}>
      {attendanceLabel(item)}
    </span>
    </div>

    <div slot="body">
        <div class="space-y-1">
            <div class="text-sm font-medium text-slate-100">
                {item.game.location}
            </div>
            <div class="flex items-center justify-between text-xs text-slate-400">
        <span>
          {#if item.game.isCanceled}
            <span>This game has been canceled.</span>
          {/if}
            </div>
        </div>

        {#if item.game.costInCents > 0}
            <div class="mt-3">
                <span class={`inline-flex items-center rounded-full px-2 py-1 text-[11px] ${paymentPillClass(item)}`}>
                  {#if item.payment}
                    {paymentLabel(item)} • {formatCurrency(item.payment.amountCents)}
                  {:else}
                    Payment needed • {formatCurrency(item.game.costInCents)}
                  {/if}
                </span>
            </div>
        {/if}
    </div>

    <div slot="footer">
    <span class="text-[11px] text-slate-500">
      Created {new Date(item.game.createdAt).toLocaleDateString()}
    </span>

        <button
                class="inline-flex items-center rounded-xl bg-emerald-500/90 px-3 py-1.5 text-xs font-semibold text-slate-950 hover:bg-emerald-500 disabled:bg-slate-700 disabled:text-slate-400 transition-all"
                disabled={item.game.isCanceled}
                type="button"
        >
            {primaryActionLabel(item)}
        </button>
    </div>
</Card>
