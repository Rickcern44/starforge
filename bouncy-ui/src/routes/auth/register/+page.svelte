<script lang="ts">
    import Register from "$lib/components/auth/Register.svelte";
    let {data, form} = $props();
</script>

<Register form={form} />