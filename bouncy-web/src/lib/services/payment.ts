import * as api from './api';
import type { Payment, PaymentAllocation, GameCharge, LeagueFinancialSummary } from '$lib/models';

export async function getPaymentsForLeague(leagueId: string, fetch?: api.Fetch, token?: string): Promise<Payment[]> {
  try {
    const payments = await api.get(`league/${leagueId}/payments`, fetch, token);
    return payments as Payment[];
  } catch (error) {
    console.error(`Error fetching payments for league ${leagueId}:`, error);
    return [];
  }
}

export async function getFinancialSummaryForLeague(leagueId: string, fetch?: api.Fetch, token?: string): Promise<LeagueFinancialSummary | null> {
  try {
    const summary = await api.get(`league/${leagueId}/financials`, fetch, token);
    return summary as LeagueFinancialSummary;
  } catch (error) {
    console.error(`Error fetching financial summary for league ${leagueId}:`, error);
    return null;
  }
}

export async function addPayment(leagueId: string, paymentData: Partial<Payment> & { recordedBy?: string }): Promise<boolean> {
  try {
    await api.post(`league/${leagueId}/payments`, paymentData);
    return true;
  } catch (error) {
    console.error(`Error adding payment to league ${leagueId}:`, error);
    return false;
  }
}

export async function addAllocation(paymentId: string, allocationData: Partial<PaymentAllocation>): Promise<boolean> {
  try {
    await api.post(`payments/${paymentId}/allocations`, allocationData);
    return true;
  } catch (error) {
    console.error(`Error adding allocation to payment ${paymentId}:`, error);
    return false;
  }
}

export async function getUserPayments(fetch?: api.Fetch, token?: string): Promise<Payment[]> {
  try {
    const payments = await api.get('me/payments', fetch, token);
    return payments as Payment[];
  } catch (error) {
    console.error('Error fetching user payments:', error);
    return [];
  }
}

export async function getUserCharges(fetch?: api.Fetch, token?: string): Promise<Payment[]> {
  try {
    const charges = await api.get('me/charges', fetch, token);
    return charges as Payment[];
  } catch (error) {
    console.error('Error fetching user charges:', error);
    return [];
  }
}
