import type {League} from "$lib/types/League";
import type {Game} from "$lib/types/Game";

export function findNextGame(
    leagues: League[],
    now: Date = new Date()
): Game | null {
    let nextGame: Game | null = null;
    let smallestDiff = Infinity;

    console.log(leagues);

    for (const league of leagues) {
        for (const game of league.Games) {
            const gameTime =
                game.startTime
                    ? game.startTime
                    : new Date(game.startTime);

            const diff = gameTime.getTime() - now.getTime();

            // skip past games
            if (diff <= 0) continue;

            if (diff < smallestDiff) {
                smallestDiff = diff;
                nextGame = game;
            }
        }
    }

    return nextGame;
}
