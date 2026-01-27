<script lang="ts">
    import GameCard, {type UpcomingGame} from './GameCard.svelte';

    export let games: UpcomingGame[] = [];
</script>

{#if games.length === 0}
    <!-- empty state -->
{:else}
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {#each games as game (game.game.id)}
            <GameCard item={game}/>
        {/each}
    </div>
{/if}