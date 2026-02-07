import {type League} from "$lib/types/League"


type LeagueState = {
    leagues: League[];
    selectedLeagueId: string | null;
};

export const leagueState = $state<LeagueState>({
    leagues: [],
    selectedLeagueId: null
});

export function setLeagues(leagues: League[]) {
    leagueState.leagues = leagues;

    // auto-select if only one league
    if (leagues.length === 1) {
        leagueState.selectedLeagueId = leagues[0].ID;
    }
}

export function selectLeague(id: string) {
    leagueState.selectedLeagueId = id;
}