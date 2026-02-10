<script lang="ts">
    let {form} = $props();
</script>

<div class="justify-center">
    <form method="POST">
        <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
            <legend class="fieldset-legend">Register</legend>

            <label for="user_name" class="label">Name</label>
            <input name="user_name" type="text" class="input" placeholder="Name"/>

            <label for="email" class="label">Email</label>
            <input name="email" type="email" class="input" placeholder="Email"/>

            <label for="password" class="label">Password</label>
            <input name="password" type="password" class="input" placeholder="Password"/>

            <label for="confirm_password" class="label">Confirm Password</label>
            <input name="confirm_password" type="password" class="input" placeholder="Password"/>
            {#if form?.error}
                <div class="toast toast-top toast-center">
                    <div class="alert alert-error ">
                        <span>{form?.error}</span>
                    </div>
                </div>
            {/if}
            <button class="btn btn-neutral mt-4">Login</button>
            <p>Already have an account? <a class="link text-blue-400" href="/auth/login">Sign in</a></p>
        </fieldset>
    </form>
</div>