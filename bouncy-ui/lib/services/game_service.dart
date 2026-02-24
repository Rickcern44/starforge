import 'package:logging/logging.dart';
import 'api_service.dart';
import 'dart:convert';
import '../models/game.dart';

class GameService {
  static final _log = Logger('GameService');

  static Future<bool> updateAttendance({
    required String gameId,
    required int status,
    required String comment,
  }) async {
    try {
      _log.info('Updating attendance for game $gameId');
      final response = await ApiService.post(
        '/game/$gameId/attendance',
        data: {
          'status': status,
          'comment': comment,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 204) {
        _log.info('Successfully updated attendance for game $gameId');
        return true;
      }

      _log.severe(
          'Failed to update attendance. Status: ${response.statusCode}, Body: ${response.body}');
      return false;
    } catch (e, stack) {
      _log.severe('Error updating attendance for game $gameId', e, stack);
      return false;
    }
  }

  static Future<Game?> getGameById(String gameId) async {
    try {
      _log.info('Fetching game $gameId');
      final response = await ApiService.get('/game/$gameId');
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data is Map<String, dynamic>) {
          return Game.fromMap(data);
        }
      }
      _log.warning(
          'Failed to fetch game $gameId. Status: ${response.statusCode}');
    } catch (e, stack) {
      _log.severe('Error fetching game $gameId', e, stack);
    }
    return null;
  }
}
