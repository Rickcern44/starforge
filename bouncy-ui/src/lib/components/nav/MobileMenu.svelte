<script lang="ts">
    import {Bolt, PlusIcon} from 'lucide-svelte';
    import {CalendarDays} from 'lucide-svelte';
    import {LayoutDashboard} from 'lucide-svelte';

</script>

<div class="md:hidden">
    <div class="dock dock-md dock-bottom bg-base-300 border-t border-base-300">
        <a
                href="/"
                class="dock-item btn btn-ghost btn-sm flex flex-col gap-0.5"
                aria-label="Dashboard"
        >
            <div class="flex items-center justify-center h-8 w-8">
                <LayoutDashboard/>
            </div>
            <span class="text-[11px]">Home</span>
        </a>

        <a
                href="/games"
                class="dock-item btn btn-ghost btn-sm flex flex-col gap-0.5"
                aria-label="Games"
        >
            <div class="flex items-center justify-center h-8 w-8">
                <CalendarDays/>
            </div>
            <span class="text-[11px]">Games</span>
        </a>
        <a
                href="/settings"
                class="dock-item btn btn-ghost btn-sm flex flex-col gap-0.5"
                aria-label="Settings"
        >
            <div class="flex items-center justify-center h-8 w-8">
                <Bolt/>
            </div>
            <span class="text-[11px]">Settings</span>
        </a>
        <a
                href="/settings"
                class="dock-item btn btn-ghost btn-sm flex flex-col gap-0.5"
                aria-label="Settings"
        >
            <div class="flex items-center justify-center h-8 w-8">
                <Bolt/>
            </div>
            <span class="text-[11px]">Settings</span>
        </a>
    </div>
    <!--Mobile Fab Menu-->
    <div class="fab mb-16 md:hidden">
        <!-- a focusable div with tabindex is necessary to work on all browsers. role="button" is necessary for accessibility -->
        <div tabindex="0" role="button" class="btn btn-lg btn-circle btn-success">F</div>

        <!-- buttons that show up when FAB is open -->
        <div>Create Game
            <button class="btn btn-lg btn-circle">
                <PlusIcon/>
            </button>
        </div>
    </div>
</div>