import { authService } from '$lib/services/auth.svelte';

export const ssr = false;
export const prerender = false;

export const load = async () => {
  await authService.initialized;
  return {};
};
