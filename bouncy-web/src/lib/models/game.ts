import type { GameAttendance } from './game-attendance';
import type { GameCharge } from './game-charge';

export interface Game {
  id: string;
  leagueId: string;
  startTime: Date;
  location: string;
  costInCents: number;
  isCanceled: boolean;
  attendance: GameAttendance[];
  charges: GameCharge[];
}
