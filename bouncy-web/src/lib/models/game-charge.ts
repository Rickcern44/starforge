import type { Game } from './game';
import type { PaymentAllocation } from './payment-allocation';

export interface GameCharge {
  id: string;
  gameId: string;
  userId: string;
  externalName?: string;
  amountCents: number;
  createdAt: Date;
  allocations: PaymentAllocation[];
  game?: Game;
}
