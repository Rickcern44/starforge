export interface Invitation {
  token: string;
  email: string;
  leagueId: string;
  invitedBy: string;
  expiresAt: string | Date;
  usedAt?: string | Date | null;
}
