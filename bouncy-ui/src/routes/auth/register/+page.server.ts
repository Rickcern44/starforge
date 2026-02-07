import {fail, redirect} from "@sveltejs/kit";
import axios from "axios";

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

        const response = await axios.post("http://localhost:3000/api/v1/auth/register", {
            email: data.get("email"),
            password: data.get("password"),
        })

        if (response.status === 409) {
            return fail(409, {
                error: "User already exists",
                description: "Error already exists please login or try again"
            })
        }

        console.log(response)
        
        redirect(301, "/")
    }
}