import 'package:ui_test/models/payment_allocation.dart';

class GameCharge {
  final String id;
  final String gameId;
  final String userId;
  final int amountCents;
  final DateTime createdAt;
  final List<PaymentAllocation> allocations;

  GameCharge({
    required this.id,
    required this.gameId,
    required this.userId,
    required this.amountCents,
    required this.createdAt,
    required this.allocations,
  });

  factory GameCharge.fromMap(Map<String, dynamic> map) {
    return GameCharge(
      id: map['ID'] ?? '',
      gameId: map['GameID'] ?? '',
      userId: map['UserID'] ?? '',
      amountCents: map['AmountCents'] ?? 0,
      createdAt: DateTime.tryParse(map['CreatedAt'] ?? '') ?? DateTime.now(),
      allocations: List<PaymentAllocation>.from(
        (map['Allocations'] ?? []).map((x) => PaymentAllocation.fromMap(x)),
      ),
    );
  }
}
