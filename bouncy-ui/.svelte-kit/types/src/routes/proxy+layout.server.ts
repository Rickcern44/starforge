// @ts-nocheck
import type {LayoutServerLoad} from "../../.svelte-kit/types/src/routes/$types";
import axios from "axios";
import {getAuthToken} from "$lib/services/cookie_utils";
import {redirect} from "@sveltejs/kit";


export const load = async ({cookies, url}: Parameters<LayoutServerLoad>[0]) => {
    if (!getAuthToken(cookies)) {

        if(url.pathname.startsWith("/auth")){
            return {}
        }

        return redirect(301, "/auth/login");
    }

    let userResponse = await axios.get("http://localhost:3000/api/v1/users/me", {
        headers: {
            Authorization: `Bearer ${getAuthToken(cookies)}`
        }
    })

    return {user: userResponse.data};
}