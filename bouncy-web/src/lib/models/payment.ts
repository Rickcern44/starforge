import type { PaymentAllocation } from './payment-allocation';

export type PaymentMethod = 'venmo' | 'cash';

export interface Payment {
  id: string;
  userId?: string | null;
  externalName: string;
  leagueId: string;
  amountInCents: number;
  method: PaymentMethod;
  receivedAt: string | Date;
  recordedBy: string;
  reference?: string | null;
  allocations: PaymentAllocation[];
}
