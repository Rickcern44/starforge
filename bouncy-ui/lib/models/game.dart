import 'game_attendance.dart';
import 'game_charge.dart';

class Game {
  final String id;
  final String leagueId;
  final DateTime startTime;
  final String location;
  final int costInCents;
  final bool isCanceled;
  final List<GameAttendance> attendance;
  final List<GameCharge> charges;

  Game({
    required this.id,
    required this.leagueId,
    required this.startTime,
    required this.location,
    required this.costInCents,
    required this.isCanceled,
    required this.attendance,
    required this.charges,
  });

  factory Game.fromMap(Map<String, dynamic> map) {
    return Game(
      id: map['ID'] ?? '',
      leagueId: map['LeagueID'] ?? '',
      startTime: DateTime.tryParse(map['StartTime'] ?? '') ?? DateTime.now(),
      location: map['Location'] ?? '',
      costInCents: map['CostInCents'] ?? 0,
      isCanceled: map['IsCanceled'] ?? false,
      attendance: List<GameAttendance>.from(
        (map['Attendance'] ?? []).map((x) => GameAttendance.fromMap(x)),
      ),
      charges: List<GameCharge>.from(
        (map['Charges'] ?? []).map((x) => GameCharge.fromMap(x)),
      ),
    );
  }
}
