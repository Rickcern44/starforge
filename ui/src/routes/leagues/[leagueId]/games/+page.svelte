<script lang="ts">
    import type {PageProps} from "./$types"
    import {games, loadGames, loading} from "$lib/stores/games"
    import {onMount} from "svelte";

    let {data}: PageProps = $props()

    onMount(() => {
        loadGames(data.pageData.leagueId)
    })
</script>

{#if $loading}
    <p>Loading Games...</p>
{:else if $games.length === 0}
    <p>There are no games!</p>
{:else}
    <p>got a game or two</p>
{/if}