<script>
    import {enhance} from '$app/forms';
    import EmailInput from "$lib/components/auth/EmailInput.svelte";
    import PasswordInput from "$lib/components/auth/PasswordInput.svelte";

</script>

<div class="flex justify-center">
    <div class="text-center">
        <p class="text-center">Login</p>
        <p class="text-center">Sign into you account</p>
        <form class="flex flex-col w-full space-y-4 mt-4" method="POST" use:enhance>
            <EmailInput name="email"/>
            <PasswordInput name="password"/>

            <button class="btn btn-neutral mt-4">Login</button>
        </form>


        <!--        <form method="POST" use:enhance>-->
        <!--            <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">-->
        <!--                <legend class="fieldset-legend">Login</legend>-->

        <!--                <label for="email" class="label">Email</label>-->
        <!--                <input name="email" type="email" class="input" placeholder="Email" required/>-->

        <!--                <label for="password" class="label">Password</label>-->
        <!--                <input name="password" type="password" class="input" placeholder="Password" required/>-->

        <!--                <button class="btn btn-neutral mt-4">Login</button>-->
        <!--                <p>Dont have an account? <a class="link text-blue-400" href="/auth/register">Register here</a></p>-->
        <!--            </fieldset>-->
        <!--        </form>-->
    </div>
</div>