<script lang="ts">
    import favicon from '$lib/assets/favicon.svg';
    import "../app.css"
    import {setUser, user} from "$lib/state/user.svelte";
    import {page} from "$app/state"
    import type {LayoutProps} from "./$types";
    import AppShell from "$lib/components/layout/AppShell.svelte";

    let {children, data}: LayoutProps = $props();

    $effect(() => {
        if (data?.user) {
            setUser(data.user);
        }
    })



    // This and the server short circuit hide the layout components when on an auth page
    let hideNavComponents = $state(!page.url.pathname.startsWith("/auth"));

</script>

<svelte:head>
    <link rel="icon" href={favicon}/>
</svelte:head>

<AppShell user={user}>
    {@render children?.()}
</AppShell>






