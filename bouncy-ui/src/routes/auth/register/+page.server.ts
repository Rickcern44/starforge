import {fail, redirect} from "@sveltejs/kit";

export const actions = {
    default: async ({cookies, request}) => {
        const data = await request.formData()
        if (data.get("password") !== data.get("confirm_password")) {
            console.log(data)
            console.log("passwords don't match")
            return fail(400, {
                description: "Passwords do not match",
                error: "Please make sure the passwords are matching"
            })
        }

        console.log("Succuss!")
        redirect(301, "/")
    }
}