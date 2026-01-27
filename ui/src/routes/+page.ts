import type {PageLoad} from "./$types";
import {authState} from "$lib";
import {redirect} from "@sveltejs/kit";

// export const load: PageLoad = async () => {
//     if(!authState.isAuthenticated) {
//         redirect(301, "/account/register");
//     }
// }