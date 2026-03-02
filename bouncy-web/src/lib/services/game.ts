import * as api from './api';
import type { Game } from '$lib/models';

export async function updateAttendance(gameId: string, status: number, comment: string): Promise<boolean> {
  try {
    await api.post(`game/${gameId}/attendance`, { status, comment });
    return true;
  } catch (error) {
    console.error(`Error updating attendance for game ${gameId}:`, error);
    return false;
  }
}

export async function removeAttendance(gameId: string, userId: string): Promise<boolean> {
  try {
    await api.request('DELETE', `game/${gameId}/attendance/${userId}`);
    return true;
  } catch (error) {
    console.error(`Error removing attendance for user ${userId} from game ${gameId}:`, error);
    return false;
  }
}

export async function getGameById(gameId: string, fetch?: api.Fetch, token?: string): Promise<Game | null> {
  try {
    const game = await api.get(`game/${gameId}`, fetch, token);
    return game as Game;
  } catch (error) {
    console.error(`Error fetching game ${gameId}:`, error);
    return null;
  }
}

export async function createGame(gameData: {
  leagueId: string;
  location: string;
  startTime: string;
  costInCents: number;
  isRecurring?: boolean;
  recurrenceInterval?: string;
}): Promise<Game | Game[] | null> {
  try {
    const game = await api.post(`league/${gameData.leagueId}/games`, gameData);
    return game as Game | Game[];
  } catch (error) {
    console.error('Error creating game:', error);
    return null;
  }
}
