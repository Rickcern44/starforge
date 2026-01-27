import axios from "axios";
import { authState } from "$lib/auth";

export const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL
});

/**
 * Request interceptor to attach Keycloak JWT token
 */
api.interceptors.request.use(async (config) => {
    try {
        const token = await authState.getAccessToken();
        
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    } catch (error) {
        console.warn('[API] Failed to get access token:', error);
    }
    
    return config;
});

/**
 * Response interceptor to handle 401 errors
 */
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // If 401 and we haven't retried yet
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                // Try to refresh the token
                const token = await authState.getAccessToken();
                
                if (token) {
                    originalRequest.headers.Authorization = `Bearer ${token}`;
                    return api(originalRequest);
                }
            } catch (refreshError) {
                console.error('[API] Token refresh failed:', refreshError);
            }

            // If refresh failed, redirect to login
            if (typeof window !== 'undefined') {
                window.location.href = '/account/signin';
            }
        }

        return Promise.reject(error);
    }
);