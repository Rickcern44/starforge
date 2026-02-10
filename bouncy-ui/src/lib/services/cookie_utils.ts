import type {Cookies} from "@sveltejs/kit";
import {dev} from "$app/environment";

// Configuration for authentication cookies
export const COOKIE_OPTIONS = {
    path: "/",
    httpOnly: true,
    secure: !dev, // Only use secure in production
    sameSite: "strict" as const,
    maxAge: 60 * 60, // 1 hour in seconds
};

export interface JwtPayload {
    userId: string;
    email: string;
    role?: string;
    iat?: number;
    exp?: number;
}

interface TokenResponse {
    success: boolean;
    message?: string;
    user?: Omit<JwtPayload, "iat" | "exp">;
    error?: any;
}

/**
 * Set authentication cookie in the response
 */
export const setAuthCookie = (
    cookies: Cookies,
    accessToken: string
): void => {
    cookies.set("access_token", accessToken, COOKIE_OPTIONS);
};

/**
 * Clear authentication cookies
 */
export const clearAuthCookies = (cookies: Cookies): void => {
    cookies.delete("access_token", {path: "/"});
    cookies.delete("user", {path: "/"}); // Clear the existing user cookie as well
};

/**
 * Get authentication token from cookies
 */
export const getAuthToken = (
    cookies: Cookies
): string | undefined => {
    return cookies.get("access_token");
};
