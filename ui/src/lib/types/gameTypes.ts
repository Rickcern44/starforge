export type PaymentStatus = 'pending' | 'completed' | 'failed' | 'refunded';
export type PaymentMethod = 'cash' | 'venmo' | 'zelle' | 'card';

export type Game = {
    id: string;
    leagueId: string;
    startTime: string;
    location: string;
    costInCents: number;
    isCanceled: boolean;
    createdAt: string;
};

export type League = {
    id: string;
    name: string;
    isActive: boolean;
    createdAt: string;
};

export type GameAttendance = {
    gameId: string;
    userId: string;
    checkedIn: boolean;
    createdAt: string;
};

export type GamePayment = {
    id: string;
    gameId: string;
    userId: string;
    amountCents: number;
    method: PaymentMethod;
    status: PaymentStatus;
    paidAt: string | null;
    confirmedBy: string | null;
};

export type UpcomingGame = {
    game: Game;
    league: League;
    attendance?: GameAttendance;
    payment?: GamePayment;
};