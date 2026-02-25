import type { Game } from './game';
import type { LeagueMember } from './league-member';

export interface League {
  id: string;
  name: string;
  isActive: boolean;
  members: LeagueMember[];
  games: Game[];
}
