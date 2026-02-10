class PaymentAllocation {
  final String paymentId;
  final String gameChargeId;
  final int amountInCents;

  PaymentAllocation({
    required this.paymentId,
    required this.gameChargeId,
    required this.amountInCents,
  });

  factory PaymentAllocation.fromMap(Map<String, dynamic> map) {
    return PaymentAllocation(
      paymentId: map['PaymentID'] ?? '',
      gameChargeId: map['GameChargeID'] ?? '',
      amountInCents: map['AmountInCents'] ?? 0,
    );
  }
}
