class GameAttendance {
  final String userId;
  final bool checkedIn;
  final String checkInComment;
  final DateTime createdAt;
  final DateTime updatedAt;

  GameAttendance({
    required this.userId,
    required this.checkedIn,
    required this.checkInComment,
    required this.createdAt,
    required this.updatedAt,
  });

  factory GameAttendance.fromMap(Map<String, dynamic> map) {
    return GameAttendance(
      userId: map['UserID'] ?? '',
      checkedIn: map['CheckedIn'] ?? false,
      checkInComment: map['CheckInComment'] ?? '',
      createdAt: DateTime.tryParse(map['CreatedAt'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(map['UpdatedAt'] ?? '') ?? DateTime.now(),
    );
  }
}
