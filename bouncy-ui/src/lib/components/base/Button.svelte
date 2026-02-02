<script lang="ts">

    let {
        variant = 'primary',
        size = 'md',
        outline = false,
        loading = false,
        disabled = false,
        class: className = '',
        type = 'button',
        children,
        ...rest
    } = $props<{
        variant?:
            | 'primary'
            | 'secondary'
            | 'accent'
            | 'ghost'
            | 'link'
            | 'neutral'
            | 'info'
            | 'success'
            | 'warning'
            | 'error';
        size?: 'sm' | 'md' | 'lg';
        outline?: boolean;
        loading?: boolean;
        disabled?: boolean;
        class?: string;
        type?: 'button' | 'submit' | 'reset';
        children?: () => any
    }>();

    const classes = $derived(() => {
        return [
            'btn',
            `btn-${variant}`,
            size !== 'md' && `btn-${size}`,
            outline && 'btn-outline',
            loading && 'btn-loading',
            disabled && 'btn-disabled',
            className
        ]
            .filter(Boolean)
            .join(' ');
    });
</script>

<button type={type} class={classes}>
    {@render children?.()}
</button>