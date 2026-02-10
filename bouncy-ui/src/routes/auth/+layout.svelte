<script>
    import "../../app.css"
    import PageContent from "$lib/components/layout/PageContent.svelte";

    let {children} = $props();
</script>

<PageContent>
    {@render children()}
</PageContent>

