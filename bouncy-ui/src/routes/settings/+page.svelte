<script>

    import {goto} from "$app/navigation";

    const navigateToLogout = () => {
        goto("/auth/logout");
    }
</script>
<div class="flex">
    <button on:click={navigateToLogout} type="button" class="btn btn-error btn-lg btn-block flex flex-1">Logout
    </button>
</div>
