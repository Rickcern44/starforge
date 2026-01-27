<article
        class="flex flex-col rounded-xl border border-slate-800 bg-slate-900/60 p-4 shadow-sm shadow-slate-950/60 hover:border-emerald-500/40 hover:bg-slate-900/80 transition-all"
>
    {#if $$slots.header}
        <header class="mb-3 flex items-start justify-between gap-3">
            <slot name="header"/>
        </header>
    {/if}

    {#if $$slots.body}
        <div class="flex-1 space-y-3 text-sm text-slate-200">
            <slot name="body"/>
        </div>
    {/if}

    {#if $$slots.footer}
        <footer class="mt-4 flex items-center justify-between gap-3 text-xs text-slate-400">
            <slot name="footer"/>
        </footer>
    {/if}
</article>