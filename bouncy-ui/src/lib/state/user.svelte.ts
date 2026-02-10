import {browser} from '$app/environment';
import {LeagueRole} from "$lib/types/LeagueRole";


type UserState = {
    value: User | null;
    leagueRole: LeagueRole | null;
};

const initial: UserState = {
    value: null,
    leagueRole: LeagueRole.USER,
};

if (browser) {
    const stored = localStorage.getItem('user');
    if (stored) {
        initial.value = JSON.parse(stored);
    }
}

export const user = $state<UserState>(initial);

export function setUser(value: User | null) {
    user.value = value;

    if (!browser) return;

    if (value) {
        localStorage.setItem('user', JSON.stringify(value));
    } else {
        localStorage.removeItem('user');
    }
}

export function clearUser() {
    localStorage.removeItem('user');
}

export function setUserLeagueRole(role: string) {
    user.leagueRole = role as LeagueRole;
}