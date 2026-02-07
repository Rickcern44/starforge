<script lang="ts">

    import type {PageProps} from './$types';
    import {findNextGame} from "$lib/utils/leagueUtils";
    import {selectLeague, setLeagues, leagueState} from "$lib/state/league.svelte";
    import PageContent from "$lib/components/layout/PageContent.svelte";
    import {Check, CheckCheck, CheckCircleIcon, ChevronRight} from "lucide-svelte";
    import type {Game} from "$lib/types/Game";

    let {data}: PageProps = $props()

    let selectedLeagueId = $state("0")

    $effect(() => {
        if (data.leagues.length > 0) {
            setLeagues(data.leagues);
        }
    })

    type UpcomingGameVM = Game & {
        leagueId: string;
        leagueName: string;
    };

    let combinedGames = $derived(() => {
        if (!leagueState.leagues.length) return [];

        return leagueState.leagues.flatMap(league =>
            league.Games
                // .filter(g => g.startTime && g.startTime > new Date())
                .map(g => ({
                    ...g,
                    leagueId: league.ID,
                    leagueName: league.Name
                }))
        );
    });
    $inspect(combinedGames());

</script>


<PageContent>
    <div class="flex flex-col items-center justify-center px-6 space-y-6 sm:px-6 lg:px-8">
        <div class="px-4lg:px-8">
            <p class="text-3xl mt-6">Dashboard</p>
        </div>
        <section>
            <ul class="list bg-base-100 rounded-box shadow-md">
                {#each combinedGames() as game}
                    <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">Upcoming Game</li>
                    <li class="list-row">

                        <div>
                            <div>{game?.location}</div>
                            <div class="text-xs uppercase font-semibold opacity-60">{game?.leagueName}</div>
                        </div>
                        <button class="btn btn-square btn-ghost text-green-600">
                            <CheckCheck/>
                        </button>
                        <button class="btn btn-square btn-ghost">
                            <ChevronRight/>
                        </button>
                    </li>
                {/each}
            </ul>
        </section>
        <div class="divider"></div>
        <section class="w-full">
            <p class="text-lg mb-2">Recent Games</p>
            <div class="stats shadow w-full">
                <div class="stat">
                    <div class="stat-title">This month</div>
                    <div class="stat-value">2</div>
                    <div class="stat-desc">Games played</div>
                </div>

                <div class="stat">
                    <div class="stat-title">Total</div>
                    <div class="stat-value">46</div>
                    <div class="stat-desc">Games played</div>
                </div>
            </div>
        </section>
        <div class="divider"></div>
        <section class="w-full">
            <p class="text-lg mb-2">Finances</p>
            <div class="stats shadow w-full">
                <div class="stat">
                    <div class="stat-title">Amount Owed</div>
                    <div class="stat-value">$7.00</div>
                </div>

                <div class="stat">
                    <div class="stat-title">Total Paid</div>
                    <div class="stat-value">$780.00</div>
                </div>
            </div>
        </section>
    </div>
</PageContent>
