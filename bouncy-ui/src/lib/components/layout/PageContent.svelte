<script lang="ts">
    import AppBar from "$lib/components/nav/AppBar.svelte";

    let {children} = $props()
</script>

<div class="flex flex-col items-center justify-center px-6 rounded-t-lg sm:px-6 lg:px-8">
    {@render children?.()}
</div>