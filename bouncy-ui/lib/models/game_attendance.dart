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
      userId: map['UserID'] ?? map['userId'] ?? map['user_id'] ?? '',
      checkedIn: map['CheckedIn'] ?? map['checkedIn'] ?? false,
      checkInComment: map['CheckInComment'] ??
          map['checkInComment'] ??
          map['check_in_comment'] ??
          '',
      createdAt:
          DateTime.tryParse(map['CreatedAt'] ?? map['createdAt'] ?? '') ??
              DateTime.now(),
      updatedAt:
          DateTime.tryParse(map['UpdatedAt'] ?? map['updatedAt'] ?? '') ??
              DateTime.now(),
    );
  }
}
