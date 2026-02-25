import type { User } from '$lib/models';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: User | null;
			token: string | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
