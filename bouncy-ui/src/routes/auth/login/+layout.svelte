<script>
    import "../../../app.css"

    let {children} = $props();
</script>

<div class={`min-h-screen flex flex-col items-center justify-center px-4 sm:px-6 lg:px-8`}>
    {@render children()}
</div>
