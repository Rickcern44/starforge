<script lang="ts">
    // This is to appease tailwind
    const cachedClasses = [
        'primary',
        'secondary',
        'accent',
        'ghost',
        'link',
        'neutral',
        'info',
        'success',
        'warning',
        'error',
        'btn-sm',
        'btn-md',
        'btn-lg',
        'btn-outline'
    ]

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


    const classDefinition = $derived(`btn btn-${variant} btn-${size} ${outline ? 'btn-outlined' : ''}`)
</script>

<button type={type} class="{classDefinition}">
    {@render children?.()}
</button>