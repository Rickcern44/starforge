/**
 * Design Tokens
 * 
 * Centralized design system tokens for consistent styling across the app.
 * Import and use these in components to maintain design consistency.
 */

export const tokens = {
    // Cards
    card: 'bg-slate-900/60 border border-slate-800 rounded-xl shadow-sm shadow-slate-950/60',
    cardHover: 'hover:border-emerald-500/40 hover:bg-slate-900/80 transition-all',
    cardPadding: 'p-4',
    cardPaddingLarge: 'p-6',

    // Buttons
    buttonPrimary: 'bg-emerald-500/90 hover:bg-emerald-500 text-slate-950 font-semibold rounded-xl px-4 py-2.5 shadow-lg shadow-emerald-500/25 transition-all disabled:bg-emerald-500/50 disabled:cursor-not-allowed',
    buttonSecondary: 'bg-slate-800/50 hover:bg-slate-800 border border-slate-700 text-slate-300 font-medium rounded-xl px-4 py-2.5 transition-all',
    buttonGhost: 'text-slate-400 hover:text-slate-200 hover:bg-slate-800 rounded-xl px-3 py-2 transition-all',
    buttonDanger: 'text-rose-400 hover:text-rose-300 hover:bg-slate-800 transition-all',

    // Inputs
    input: 'bg-slate-800/50 border border-slate-700 rounded-xl px-4 py-2.5 text-slate-200 placeholder-slate-500 focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500/50 transition-all',
    
    // Focus states
    focusRing: 'focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500/50',

    // Typography
    headingPage: 'text-2xl lg:text-3xl font-bold',
    headingSection: 'text-lg font-semibold',
    textBody: 'text-sm text-slate-300',
    textMuted: 'text-xs text-slate-500',

    // Navigation
    navItem: 'px-3 py-2 text-sm font-medium rounded-xl transition-all',
    navItemActive: 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40',
    navItemInactive: 'text-slate-400 hover:bg-slate-800 hover:text-slate-200',

    // Pills/Badges
    pillBase: 'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium',
    pillSuccess: 'bg-emerald-500/10 text-emerald-300 ring-1 ring-emerald-500/40',
    pillWarning: 'bg-amber-500/10 text-amber-300 ring-1 ring-amber-500/40',
    pillDanger: 'bg-rose-500/10 text-rose-300 ring-1 ring-rose-500/40',
    pillNeutral: 'bg-slate-800 text-slate-300 ring-1 ring-slate-700/60',

    // Dropdown/Menu
    dropdown: 'bg-slate-900/95 backdrop-blur-xl border border-slate-800 rounded-xl shadow-sm shadow-slate-950/60',
    dropdownItem: 'flex items-center gap-3 px-4 py-2 text-sm text-slate-300 hover:bg-slate-800 hover:text-slate-100 transition-all',

    // Loading
    spinner: 'animate-spin rounded-full border-b-2 border-emerald-500',
} as const;

/**
 * Combine multiple token classes
 */
export function cn(...classes: (string | undefined | null | false)[]): string {
    return classes.filter(Boolean).join(' ');
}
