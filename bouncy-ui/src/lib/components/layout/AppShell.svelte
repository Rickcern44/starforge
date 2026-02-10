<script lang="ts">
    import AppBar from "$lib/components/nav/AppBar.svelte";
    import MobileMenu from "$lib/components/nav/MobileMenu.svelte";

    let {children, user} = $props()
</script>

<AppBar user={user}/>
<main class="flex-1 bg-base-100 rounded-t-2xl p-2">
    {@render children?.()}
</main>
<MobileMenu/>
