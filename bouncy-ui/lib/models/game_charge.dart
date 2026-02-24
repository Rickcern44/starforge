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
      id: map['ID'] ?? map['id'] ?? '',
      gameId: map['GameID'] ?? map['gameId'] ?? map['game_id'] ?? '',
      userId: map['UserID'] ?? map['userId'] ?? map['user_id'] ?? '',
      amountCents: map['AmountCents'] ?? map['amountCents'] ?? 0,
      createdAt:
          DateTime.tryParse(map['CreatedAt'] ?? map['createdAt'] ?? '') ??
              DateTime.now(),
      allocations: List<PaymentAllocation>.from(
        ((map['Allocations'] ?? map['allocations']) ?? [])
            .map((x) => PaymentAllocation.fromMap(x)),
      ),
    );
  }
}
