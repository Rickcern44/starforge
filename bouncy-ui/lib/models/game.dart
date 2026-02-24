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
      id: map['ID'] ?? map['id'] ?? '',
      leagueId: map['LeagueID'] ?? map['leagueId'] ?? map['league_id'] ?? '',
      startTime:
          DateTime.tryParse(map['StartTime'] ?? map['startTime'] ?? '') ??
              DateTime.now(),
      location: map['Location'] ?? map['location'] ?? '',
      costInCents: map['CostInCents'] ?? map['costInCents'] ?? 0,
      isCanceled: map['IsCanceled'] ?? map['isCanceled'] ?? false,
      attendance: List<GameAttendance>.from(
        ((map['Attendance'] ?? map['attendance']) ?? [])
            .map((x) => GameAttendance.fromMap(x)),
      ),
      charges: List<GameCharge>.from(
        ((map['Charges'] ?? map['charges']) ?? [])
            .map((x) => GameCharge.fromMap(x)),
      ),
    );
  }
}
