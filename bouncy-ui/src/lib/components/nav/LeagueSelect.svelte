<script lang="ts">
    import {leagueState, selectLeague} from '$lib/state/league.svelte';


    function onChange(event: Event) {
        const target = event.target as HTMLSelectElement;
        selectLeague(target.value);
    }
</script>

{#if leagueState.leagues.length === 1}
    {@const league = leagueState.leagues[0]}
    <span class="font-semibold">
    {league.Name}
  </span>
{:else if leagueState.leagues.length > 1}
    <select
            class="select appearance-none"
            onchange={onChange}
    >
        {#each leagueState.leagues as league}
            <option value={league.ID}>
                {league.Name}
            </option>
        {/each}
    </select>
{:else}
    <span>No leagues</span>
{/if}
