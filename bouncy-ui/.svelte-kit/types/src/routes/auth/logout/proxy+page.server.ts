// @ts-nocheck
import type {PageServerLoad} from './$types';
import {clearAuthCookies} from "$lib/services/cookie_utils";
import {redirect} from "@sveltejs/kit";


export const load = async ({cookies}: Parameters<PageServerLoad>[0]) => {
    clearAuthCookies(cookies)

    redirect(301, "/auth/login")
}