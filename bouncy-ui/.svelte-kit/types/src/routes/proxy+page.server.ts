// @ts-nocheck
import type {PageServerLoad} from './$types';
import axios from "axios";
import type {League} from "$lib/types/League";
import {getAuthToken} from "$lib/services/cookie_utils";

export const load = async ({cookies}: Parameters<PageServerLoad>[0]) => {
    let response = await axios.get<League[]>("http://localhost:3000/api/v1/me/leagues", {
        headers: {
            Authorization: `Bearer ${getAuthToken(cookies)}`
        }
    })

    return {leagues: response.data};
}
