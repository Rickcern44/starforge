<script lang="ts">
    import LeagueSelect from "$lib/components/nav/LeagueSelect.svelte";

    let {user} = $props();

    export function formatUserInitials(): string {
        if (!user?.name) {
            return ''
        }
        console.log(user)

        let name = user?.name;

        const parts = name
            .trim()
            .split(/\s+/)
            .filter(Boolean);

        if (parts.length === 0) return '';

        if (parts.length === 1) {
            return parts[0][0].toUpperCase();
        }

        return (
            parts[0][0] + parts[1][0]
        ).toUpperCase();
    }
</script>

<header>
    <nav class="navbar bg-secondary  border-b border-base-300">
        <div class="navbar bg-base-30">
            <div class="flex-1">
                <a href="/" class="btn btn-ghost text-xl">Bouncy League Manager</a>
            </div>
            <div class="flex gap-2">
                <div class="dropdown dropdown-end">
                    <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
                        <div class="w-10 rounded-full">
                            <img
                                    alt="Tailwind CSS Navbar component"
                                    src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"/>
                        </div>
                    </div>
                    <ul
                            tabindex="-1"
                            class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
                        <li><a href="/settings">Settings</a></li>
                        <li><a href="/auth/logout">Logout</a></li>
                    </ul>
                </div>
            </div>
        </div>
    </nav>
</header>